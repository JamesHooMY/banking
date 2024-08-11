package rest

import (
	"fmt"

	transactionHdl "banking/app/api/restful/v1/handler/transaction"
	userHdl "banking/app/api/restful/v1/handler/user"
	transactionRepo "banking/app/repo/mysql/transaction"
	userRepo "banking/app/repo/mysql/user"
	transactionSrv "banking/app/service/transaction"
	userSrv "banking/app/service/user"
	_ "banking/docs"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.elastic.co/apm/module/apmgin/v2"
	"go.elastic.co/apm/v2"
	"gorm.io/gorm"
)

func InitRouter(router *gin.Engine, masterDB *gorm.DB, slaveDB *gorm.DB, tracer *apm.Tracer) *gin.Engine {
	// Middleware
	router.Use(apmgin.Middleware(router, apmgin.WithTracer(tracer))) // APM gin middleware

	// Swagger
	// docs.SwaggerInfo.BasePath = fmt.Sprintf("/api/%s", viper.GetString("server.apiVersion"))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// User handler with master and slave DBs
	userHandler := userHdl.NewUserHandler(
		userSrv.NewUserService(
			userRepo.NewUserQueryRepo(slaveDB),    // Read operations
			userRepo.NewUserCommandRepo(masterDB), // Write operations
		),
	)

	// Transaction handler with master DB
	transactionHandler := transactionHdl.NewTransactionHandler(
		transactionSrv.NewTransactionService(
			transactionRepo.NewTransactionCommandRepo(masterDB), // Write operations
		),
	)

	// v1 group
	v1 := router.Group(fmt.Sprintf("/api/%s", viper.GetString("server.apiVersion")))

	// user router
	user := v1.Group("/user")
	user.POST("", userHandler.CreateUser())
	user.GET("", userHandler.GetUsers())
	user.GET("/:id", userHandler.GetUser())
	user.POST("/transfer", transactionHandler.Transfer())
	user.POST("/deposit", transactionHandler.Deposit())
	user.POST("/withdraw", transactionHandler.Withdraw())

	return router
}
