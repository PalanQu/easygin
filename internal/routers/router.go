package routers

import (
	"easygin/internal/controllers"
	"easygin/internal/middlewares"
	"easygin/internal/models"
	"easygin/internal/services"
	"easygin/pkg/ent"
	"easygin/pkg/prom"

	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
)

func SetupRouter(routerPrefix string, db *ent.Client) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	p := prom.New(r)
	r.Use(p.Instrument())
	r.Use(middlewares.ResponseWithErrorContext())
	r.Use(middlewares.LoggerContext())

	registerPingRouters(r)
	registerUserRouters(routerPrefix, r, db, p)
	return r
}

func registerPingRouters(r *gin.Engine) {
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
}

func registerUserRouters(
	prefix string,
	r *gin.Engine,
	db *ent.Client,
	ginprom *ginprom.Prometheus,
) {
	userService := services.NewUserService(db, ginprom)
	userController := controllers.NewUserController(userService)
	userGroup := r.Group(prefix + "/users")
	{
		userGroup.GET("", Endpoint(
			userController.GetUsers,
			EmptyBinder[struct{}](),
		))
		userGroup.POST("", Endpoint(
			userController.CreateUser,
			JSONBinder[*models.CreateUserRequest](),
		))
	}
}
