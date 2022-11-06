package routes

import (
	"net/http"

	"library-backend/config"
	"library-backend/routes/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler struct {
	config *config.Config
	router *echo.Echo
}

func SkipperFn(skipURLs []string) func(echo.Context) bool {
	return func(context echo.Context) bool {
		for _, url := range skipURLs {
			if url == context.Request().URL.String() {
				return true
			}
		}
		return false
	}
}

func NewApi(cfg *config.Config) *echo.Echo {
	app := &Handler{
		config: cfg,
		router: echo.New(),
	}

	app.router.HideBanner = true

	app.router.GET("/api/v1", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/")
	})

	debug := app.router.Group("/api/v1/debug")
	student := app.router.Group("/api/v1/student")
	admin := app.router.Group("/api/v1/admin")

	admin.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(app.config.Server.Secret),
		Skipper:    SkipperFn([]string{"/api/v1/admin/login"}),
	}))

	debug.GET("/health", handler.Health)
	student.POST("/login", handler.Login(app.config.Database))
	admin.POST("/login", handler.AdminLogin(app.config.Server.Secret, app.config.Auth.Username, app.config.Auth.Password))
	admin.GET("/list", handler.List(app.config.Database))

	return app.router
}
