package rest

import (
	"fmt"

	userHdl "banking/app/api/v1/handler/user"

	userRepo "banking/app/repo/mysql/user"
	userSrv "banking/app/service/user"
	_ "banking/docs"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func InitRouter(router *gin.Engine, db *gorm.DB) *gin.Engine {
	// middleware
	// * if need cors then uncomment this line
	// router.Use(middleware.Cors())

	// swagger
	// docs.SwaggerInfo.BasePath = fmt.Sprintf("/api/%s", viper.GetString("server.apiVersion"))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// user handler
	userHandler := userHdl.NewUserHandler(userSrv.NewUserService(
		userRepo.NewUserQueryRepo(db),
		userRepo.NewUserCommandRepo(db),
	))

	// v1 group
	v1 := router.Group(fmt.Sprintf("/api/%s", viper.GetString("server.apiVersion")))

	// user router
	user := v1.Group("/user")
	user.POST("", userHandler.CreateUser())
	user.GET("", userHandler.GetUsers())
	user.GET("/:id", userHandler.GetUser())
	user.POST("/transfer", userHandler.Transfer())
	user.POST("/deposit", userHandler.Deposit())
	user.POST("/withdraw", userHandler.Withdraw())

	return router
}
