package main

import (
	"net/http"
	"os"

	"github.com/bustafed/finalBackC4_G9/cmd/config"
	"github.com/bustafed/finalBackC4_G9/cmd/external/database"
	"github.com/bustafed/finalBackC4_G9/cmd/middlewares"
	"github.com/bustafed/finalBackC4_G9/cmd/server/handler"
	"github.com/bustafed/finalBackC4_G9/internal/dentists"
	"github.com/bustafed/finalBackC4_G9/internal/patients"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
)

func main() {

	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	if env == "local" {
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
	}

	cfg, err := config.NewConfig(env)

	if err != nil {
		panic(err)
	}

	authMidd := middlewares.NewAuth(cfg.PublicConfig.PublicKey, cfg.PrivateConfig.SecretKey)

	router := gin.New()

	customRecovery := gin.CustomRecovery(middlewares.RecoveryWithLog)

	router.Use(customRecovery)

	// docs endpoint
	docs.SwaggerInfo.Host = os.Getenv("HOST")
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"ok": "ok"})
	})

	postgresDatabase, err := database.NewPostgresSQLDatabase(cfg.PublicConfig.PostgresHost,
		cfg.PublicConfig.PostgresPort, cfg.PublicConfig.PostgresUser, cfg.PrivateConfig.PostgresPassword,
		cfg.PublicConfig.PostgresDatabase)

	if err != nil {
		panic(err)
	}

	myDatabase := database.NewDatabase(postgresDatabase)

	patientsService := patients.NewService(myDatabase)

	patientsHandler := handler.NewPatientsHandler(patientsService, patientsService, patientsService)

	patientsGroup := router.Group("/patients")

	patientsGroup.GET("/:id", patientsHandler.GetPatientByID)

	patientsGroup.PUT("/:id", authMidd.AuthHeader, patientsHandler.PutPatient)

	patientsGroup.PATCH("/:id", authMidd.AuthHeader, patientsHandler.ModifyPatientByProperty)

	patientsGroup.POST("/", authMidd.AuthHeader, patientsHandler.CreatePatient)

	patientsGroup.DELETE("/:id", authMidd.AuthHeader, patientsHandler.DeletePatientByID)

	dentistService := dentists.NewService(myDatabase)

	dentistsHandler := handler.NewDentistsHandler(dentistService, dentistService, dentistService)

	dentistsGroup := router.Group("/dentists")

	dentistsGroup.GET("/:id", dentistsHandler.GetDentistByID)

	dentistsGroup.POST("/", authMidd.AuthHeader, dentistsHandler.CreateDentist)

	dentistsGroup.PUT("/:id", authMidd.AuthHeader, dentistsHandler.FullUpdateDentistByID)

	dentistsGroup.PATCH("/:id", authMidd.AuthHeader, dentistsHandler.UpdateDentistByID)

	dentistsGroup.DELETE("/:id", authMidd.AuthHeader, dentistsHandler.DeleteDentistByID)

	err = router.Run()

	if err != nil {
		panic(err)
	}
}
