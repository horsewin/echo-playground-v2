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

// configureXRay X-Rayの設定を行う
func configureXRay(apiConfig *utils.APIConfig, logger zerolog.Logger) {
	if !apiConfig.EnableTracing {
		return
	}

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

// setupRequestLogger リクエストロガーミドルウェアを設定
func setupRequestLogger(logger zerolog.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
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
				Str("latency", v.Latency.String()).
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
	})
}

// setupXRayMiddleware X-Rayミドルウェアを設定
func setupXRayMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
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
	}
}

// registerRoutes ルートを登録する
func registerRoutes(e *echo.Echo) {
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
		e.POST("/v1/notifications/read", notificationHandler.PostNotificationsRead())

	}
}

// Router ...
func Router() *echo.Echo {
	// Setup
	logger := zerolog.New(os.Stdout)
	e := echo.New()
	apiConfig := utils.NewAPIConfig()

	// Configure X-Ray
	configureXRay(apiConfig, logger)

	// Configure Echo settings
	e.HideBanner = true
	e.HidePort = false

	// Setup middlewares
	setupMiddlewares(e, logger)

	// Register routes
	registerRoutes(e)

	return e
}

// setupMiddlewares ミドルウェアを設定
func setupMiddlewares(e *echo.Echo, logger zerolog.Logger) {
	// リクエストIDの生成
	e.Use(middleware.RequestID())

	// ログ出力設定
	e.Use(setupRequestLogger(logger))

	// X-Rayミドルウェア
	e.Use(setupXRayMiddleware())

	// recoveryミドルウェアの設定
	e.Use(middleware.Recover())
}
