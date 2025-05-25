package infrastructure

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-xray-sdk-go/xray"
	handlers "github.com/horsewin/echo-playground-v2/handler"
	"github.com/horsewin/echo-playground-v2/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	zerologlog "github.com/rs/zerolog/log"
)

const (
	// TODO: 適切な名前に変更をする
	projectName = "echo-playground-v2"
)

// // customHTTPErrorHandler handles errors and separates client-facing messages from internal logs
// func customHTTPErrorHandler(err error, c echo.Context) {
// 	fmt.Println("customHTTPErrorHandler")
// 	// Get logger from context
// 	logger := zerolog.Ctx(c.Request().Context())

// 	// Prepare client-facing response
// 	code := http.StatusInternalServerError
// 	msg := "Internal server error" // Generic message for production

// 	// If it's an Echo HTTP error, use its code and possibly its message
// 	if he, ok := err.(*echo.HTTPError); ok {
// 		code = he.Code
// 		// For 4xx errors, we can be more specific with the client
// 		if code >= 400 && code < 500 {
// 			msg = fmt.Sprintf("%v", he.Message)
// 		}
// 	}

// 	// In development, we might want to show more details
// 	if os.Getenv("APP_ENV") == "development" {
// 		msg = err.Error()
// 	}

// 	// Send response to client
// 	// c.JSON() は第二引数にnilを渡すと空のJSONボディになることがあるので注意
// 	// 適切なエラーレスポンス形式 (例: {"message": "エラーメッセージ"}) にする
// 	if !c.Response().Committed { // レスポンスがまだ送信されていなければ送信
// 		if err := c.JSON(code, customErrors.ErrorMessageDef{Code: code, Message: msg}); err != nil {
// 			// JSON送信自体でエラーが起きた場合はさらにログを出すなど
// 			logger.Error().Err(err).Msg("Failed to send error JSON response")
// 		}
// 	}
// }

// Router ...
func Router() *echo.Echo {
	// Setup Zerolog
	logger := zerolog.New(os.Stdout)

	e := echo.New()
	apiConfig := utils.NewAPIConfig()

	// Set custom error handler
	// e.HTTPErrorHandler = customHTTPErrorHandler

	// X-Ray設定
	if apiConfig.EnableTracing {
		if err := xray.Configure(xray.Config{
			DaemonAddr:     "127.0.0.1:2000", // X-Rayデーモンのアドレス
			ServiceVersion: "2.14.0",
		}); err != nil {
			logger.Error().Err(err).Msg("Failed to configure X-Ray")
			// X-Ray設定失敗時はデフォルトの設定を使用
			if configErr := xray.Configure(xray.Config{}); configErr != nil {
				logger.Error().Err(configErr).Msg("Failed to configure default X-Ray settings")
			}
		}
		os.Setenv("AWS_XRAY_CONTEXT_MISSING", "LOG_ERROR")
	}

	// ----------------------------
	// Middlewareの設定
	// ----------------------------
	// リクエストIDの生成
	e.Use(middleware.RequestID())
	// ログ出力設定
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogUserAgent: true,
		LogMethod:    true,
		LogLatency:   true,
		LogRemoteIP:  true,
		LogRequestID: true,
		LogError:     true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			event := logger.Info()
			if v.Status >= 500 {
				event = logger.Error()
			} else if v.Status >= 400 {
				event = logger.Warn()
			}

			rid := c.Response().Header().Get(echo.HeaderXRequestID)

			event.Str("request_id", rid).
				Str("time", time.Now().Format(time.RFC3339Nano)).
				Str("remote_ip", v.RemoteIP).
				Str("method", v.Method).
				Str("uri", v.URI).
				Int("status", v.Status).
				Str("user_agent", v.UserAgent)

			// エラーが発生した場合、エラー情報をログに記録
			if v.Error != nil {
				event.Err(v.Error)
				event.Msg("failed api request")

				return v.Error
			}

			event.Msg("completed api request")

			return nil
		},
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/healthcheck"
		},
	}))
	// ハンドラの設定。この内部でロガーの設定をしてコンテキストにリクエストごとのロガーをセットする
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rid := c.Response().Header().Get(echo.HeaderXRequestID)
			// Create request-specific logger with request details
			reqLogger := zerologlog.With().
				Str("request_id", rid).
				Str("method", c.Request().Method).
				Str("path", c.Request().URL.Path).
				Logger()

			// healthcheckはX-Rayのセグメント作成を行わない
			if c.Path() == "/healthcheck" {
				return next(c)
			}

			req := c.Request()
			res := c.Response()

			var seg *xray.Segment
			ctx := req.Context()

			// セグメント作成を試みる
			ctx, seg = xray.BeginSegment(ctx, projectName)
			if seg == nil {
				// セグメント作成に失敗した場合はエラーをログに記録し、
				// 通常のリクエスト処理を継続
				reqLogger.Err(fmt.Errorf("failed to create X-Ray segment"))
				return next(c)
			}

			var err error
			defer func() {
				seg.Close(err)
			}()

			// リクエストのコンテキストを更新
			c.SetRequest(req.WithContext(reqLogger.WithContext(ctx)))

			// 次のハンドラーを呼び出し
			err = next(c)

			// レスポンスステータスをセグメントに記録
			if addErr := seg.AddMetadata("response_status", res.Status); addErr != nil {
				reqLogger.Err(fmt.Errorf("failed to add response_status metadata: %v", addErr))
			}

			// エラーが発生した場合、セグメントにエラー情報を記録
			if err != nil {
				if addErr := seg.AddError(err); addErr != nil {
					reqLogger.Err(fmt.Errorf("failed to add error metadata: %v", addErr))
				}
			}

			return err
		}
	})

	// recoveryミドルウェアの設定
	e.Use(middleware.Recover())

	// デフォルトで出力するバナーを隠す
	e.HideBanner = true

	// デフォルトで出力するポート番号はそのまま
	e.HidePort = false

	healthCheckHandler := handlers.NewHealthCheckHandler()
	helloWorldHandler := handlers.NewHelloWorldHandler()

	// ---------------------------
	// APIルートの定義
	// ---------------------------
	e.GET("/", healthCheckHandler.HealthCheck())
	e.GET("/healthcheck", healthCheckHandler.HealthCheck())
	e.GET("/v1/helloworld", helloWorldHandler.SayHelloWorld())
	e.GET("/v1/helloworld/error", helloWorldHandler.SayError())
	if os.Getenv("DB_CONN") == "1" {
		sqlHandler := NewSQLHandler()
		petHandler := handlers.NewPetHandler(sqlHandler)
		notificationHandler := handlers.NewNotificationHandler(sqlHandler)

		e.GET("/v1/pets", petHandler.GetPets())
		e.POST("/v1/pets/:id/like", petHandler.UpdateLike())
		e.POST("/v1/pets/:id/reservation", petHandler.Reservation())

		e.GET("/v1/notifications", notificationHandler.GetNotifications())
		e.GET("/v1/notifications/count", notificationHandler.GetUnreadNotificationCount())
		e.POST("/v1/notifications/read", notificationHandler.PostNotificationsRead())

	}

	return e
}
