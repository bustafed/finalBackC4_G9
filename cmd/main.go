package main

import (
	"net/http"
	"os"

	"github.com/bustafed/finalBackC4_G9/cmd/config"
	"github.com/bustafed/finalBackC4_G9/cmd/external/database"
	"github.com/bustafed/finalBackC4_G9/cmd/middlewares"
	"github.com/bustafed/finalBackC4_G9/cmd/server/handler"
	"github.com/bustafed/finalBackC4_G9/docs"
	"github.com/bustafed/finalBackC4_G9/internal/dentists"
	"github.com/bustafed/finalBackC4_G9/internal/patients"
	"github.com/bustafed/finalBackC4_G9/internal/appointments"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Proyecto Final Especialización backend 3
// @version 1.0
// @description Este es el proyecto final de la materia espcialización backend 3 para crear, editar, consultar y borrar las entidades de paciente, dentista y turnos.

// @contact.name API Support (Natalia Garcia, Federico Bustamante, Damian, Camilo Zuleta)

// @license.name Apache 2.0
// @license.url http:www.apache.org/licenses/LICENSE-2.0.html
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

	appointmentService := appointments.NewService(myDatabase)

	appointmentsHandler := handler.NewAppointmentsHandler(appointmentService, appointmentService, appointmentService)

	appointmentsGroup := router.Group("/appointments")

	appointmentsGroup.GET("/:id", appointmentsHandler.GetAppointmentByID)
	
	appointmentsGroup.GET("/", appointmentsHandler.GetAppointmentByDni)

	appointmentsGroup.POST("/", authMidd.AuthHeader, appointmentsHandler.CreateAppointment)

	 appointmentsGroup.PUT("/:id", authMidd.AuthHeader, appointmentsHandler.FullUpdateAppointmentByID)

	appointmentsGroup.PATCH("/:id", authMidd.AuthHeader, appointmentsHandler.UpdateAppointmentByID)

	appointmentsGroup.DELETE("/:id", authMidd.AuthHeader, appointmentsHandler.DeleteAppointmentByID)

	err = router.Run()

	if err != nil {
		panic(err)
	}
}
