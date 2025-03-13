package infrastructure

import (
	"os"

	"github.com/aws/aws-xray-sdk-go/xray"
	handlers "github.com/horsewin/echo-playground-v2/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

// Router ...
func Router() *echo.Echo {
	e := echo.New()

	// X-Ray設定
	if err := xray.Configure(xray.Config{
		DaemonAddr:     "127.0.0.1:2000", // X-Rayデーモンのアドレス
		ServiceVersion: "1.0.0",
	}); err != nil {
		e.Logger.Errorf("Failed to configure X-Ray: %v", err)
	}

	// コンテキストがない場合のエラー処理を設定
	os.Setenv("AWS_XRAY_CONTEXT_MISSING", "LOG_ERROR")

	// Middleware
	logger := middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"id":"${id}","time":"${time_rfc3339}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}"}` + "\n",
		Output: os.Stdout,
	})
	e.Use(logger)
	e.Use(middleware.Recover())

	// X-Rayミドルウェアを追加（修正版）
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			println("x-ray middleware")
			req := c.Request()
			res := c.Response()

			// リクエストごとに新しいセグメントを作成
			ctx, seg := xray.BeginSegment(req.Context(), "echo-playground-v2")
			defer seg.Close(nil)

			// リクエストのコンテキストを更新
			c.SetRequest(req.WithContext(ctx))

			// 次のハンドラーを呼び出し
			err := next(c)

			// レスポンスステータスをセグメントに記録
			if addErr := seg.AddMetadata("response_status", res.Status); addErr != nil {
				c.Logger().Errorf("Failed to add response_status metadata: %v", addErr)
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
