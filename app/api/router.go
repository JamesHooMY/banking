package rest

import (
	"fmt"
	"time"

	transactionHdl "banking/app/api/restful/v1/handler/transaction"
	userHdl "banking/app/api/restful/v1/handler/user"
	"banking/app/api/restful/v1/middleware"
	apiKeyRepo "banking/app/repo/mysql/apikey"
	transactionRepo "banking/app/repo/mysql/transaction"
	userRepo "banking/app/repo/mysql/user"
	apiKeyRedisRepo "banking/app/repo/redis/apikey"
	jwtRedisRepo "banking/app/repo/redis/jwt"
	apiKeySrv "banking/app/service/apikey"
	authSrv "banking/app/service/auth"
	transactionSrv "banking/app/service/transaction"
	userSrv "banking/app/service/user"
	_ "banking/docs"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.elastic.co/apm/module/apmgin/v2"
	"go.elastic.co/apm/v2"
	"gorm.io/gorm"
)

func InitRouter(router *gin.Engine, masterDB *gorm.DB, slaveDB *gorm.DB, redisClient *redis.Client, tracer *apm.Tracer) *gin.Engine {
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
			userRepo.NewUserCommandRepo(masterDB),            // Write operations
			userRepo.NewUserQueryRepo(slaveDB),               // Read operations
			jwtRedisRepo.NewRedisJWTCommandRepo(redisClient), // Write operations
			jwtRedisRepo.NewRedisJWTQueryRepo(redisClient),   // Read operations
		),
		apiKeySrv.NewAPIKeyService(
			apiKeyRedisRepo.NewRedisAPIKeyCommandRepo(redisClient), // Write operations
			apiKeyRedisRepo.NewRedisAPIKeyQueryRepo(redisClient),   // Read operations
			apiKeyRepo.NewAPIKeyCommandRepo(masterDB),              // Write operations
			apiKeyRepo.NewAPIKeyQueryRepo(slaveDB),                 // Read operations

		),
	)

	// Transaction handler with master DB and slave DB
	transactionHandler := transactionHdl.NewTransactionHandler(
		transactionSrv.NewTransactionService(
			transactionRepo.NewTransactionCommandRepo(masterDB), // Write operations
			transactionRepo.NewTransactionQueryRepo(slaveDB),    // Read operations
		),
	)

	// v1 group
	v1 := router.Group(fmt.Sprintf("/api/%s", viper.GetString("server.apiVersion")))

	// user router
	user := v1.Group("/user")
	user.POST("/register", userHandler.CreateUser())
	user.POST("/login", userHandler.Login())

	userAuthenticated := user.Group("", middleware.JWTAuthMiddleware())
	userAuthenticated.GET("/:userId", userHandler.GetUsers())
	userAuthenticated.POST("/apikey", userHandler.CreateAPIKey())
	userAuthenticated.GET("/apikey", userHandler.GetAPIKeys())

	transaction := v1.Group("/transaction", middleware.RateLimitMiddleware(redisClient, 10, time.Minute), middleware.APIKeyAuthMiddleware(authSrv.NewAuthService(
		apiKeyRedisRepo.NewRedisAPIKeyCommandRepo(redisClient), // Write operations
		apiKeyRedisRepo.NewRedisAPIKeyQueryRepo(redisClient),   // Read operations
		apiKeyRepo.NewAPIKeyQueryRepo(slaveDB),                 // Read operations
	)))
	transaction.POST("/transfer", transactionHandler.Transfer())
	transaction.POST("/deposit", transactionHandler.Deposit())
	transaction.POST("/withdraw", transactionHandler.Withdraw())
	transaction.GET("", transactionHandler.GetTransactions())

	return router
}
