package infrastructure

import (
	"os"

	"github.com/aws/aws-xray-sdk-go/xray"
	handlers "github.com/horsewin/echo-playground-v2/handler"
	"github.com/horsewin/echo-playground-v2/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const (
	// TODO: 適切な名前に変更をする
	projectName = "echo-playground-v2"
)

// Router ...
func Router() *echo.Echo {
	e := echo.New()
	apiConfig := utils.NewAPIConfig()

	// X-Ray設定
	if apiConfig.EnableTracing {
		if err := xray.Configure(xray.Config{
			DaemonAddr:     "127.0.0.1:2000", // X-Rayデーモンのアドレス
			ServiceVersion: "2.14.0",
		}); err != nil {
			e.Logger.Errorf("Failed to configure X-Ray: %v", err)
			// X-Ray設定失敗時はデフォルトの設定を使用
			if configErr := xray.Configure(xray.Config{}); configErr != nil {
				e.Logger.Fatalf("Failed to configure default X-Ray settings: %v", configErr)
			}
		}
		os.Setenv("AWS_XRAY_CONTEXT_MISSING", "LOG_ERROR")
	}

	// Middleware
	logger := middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"id":"${id}","time":"${time_rfc3339}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}"}` + "\n",
		Output: os.Stdout,
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/healthcheck"
		},
	})
	e.Use(logger)
	e.Use(middleware.Recover())
	e.Logger.SetLevel(log.INFO)

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
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
				c.Logger().Errorf("Failed to create X-Ray segment")
				return next(c)
			}

			var err error
			defer func() {
				seg.Close(err)
			}()

			// リクエストのコンテキストを更新
			c.SetRequest(req.WithContext(ctx))

			// 次のハンドラーを呼び出し
			err = next(c)

			// レスポンスステータスをセグメントに記録
			if addErr := seg.AddMetadata("response_status", res.Status); addErr != nil {
				c.Logger().Errorf("Failed to add response_status metadata: %v", addErr)
			}

			// エラーが発生した場合、セグメントにエラー情報を記録
			if err != nil {
				if addErr := seg.AddError(err); addErr != nil {
					c.Logger().Errorf("Failed to add error to segment: %v", addErr)
				}
			}

			return err
		}
	})

	e.Logger.SetLevel(log.INFO)
	e.HideBanner = true
	e.HidePort = false

	healthCheckHandler := handlers.NewHealthCheckHandler()
	helloWorldHandler := handlers.NewHelloWorldHandler()

	// APIルートの定義
	e.GET("/", healthCheckHandler.HealthCheck())
	e.GET("/healthcheck", healthCheckHandler.HealthCheck())
	e.GET("/v1/helloworld", helloWorldHandler.SayHelloWorld())
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
