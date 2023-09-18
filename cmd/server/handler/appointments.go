package handler

import (
	"github.com/bustafed/finalBackC4_G9/internal/appointments"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type AppointmentsGetter interface {
	GetAppointmentByID(id int) (appointments.Appointment, error)
	GetAppointmentByDni(dni string) ([]appointments.Appointment, error)
}

type AppointmentCreator interface {
	// UpdateAppointmentByID(id int, appointment appointments.Appointment) (appointments.Appointment, error)
	// UpdateAppointmentByDni(dni string, appointment appointments.Appointment) (appointments.Appointment, error)
	CreateAppointment(appointment appointments.Appointment) (appointments.Appointment, error)
}
type AppointmentDeleter interface {
	// DeleteAppointmentByID(id int) error
}

type AppointmentsHandler struct {
	appointmentsGetter  AppointmentsGetter
	appointmentsCreator AppointmentCreator
	appointmentDeleter  AppointmentDeleter
}

func NewAppointmentsHandler(getter AppointmentsGetter, creator AppointmentCreator, deleter AppointmentDeleter) *AppointmentsHandler {
	return &AppointmentsHandler{
		appointmentsGetter:  getter,
		appointmentsCreator: creator,
		appointmentDeleter:  deleter,
	}
}

func (ah *AppointmentsHandler) GetAppointmentByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	appointment, err := ah.appointmentsGetter.GetAppointmentByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, appointment)
}

func (ah *AppointmentsHandler) GetAppointmentByDni(ctx *gin.Context) {
	dni := ctx.Query("dni")

	
	if dni == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid dni"})
		return
	}

	appointments, err := ah.appointmentsGetter.GetAppointmentByDni(dni)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, appointments)
}

func (ah *AppointmentsHandler) CreateAppointment(ctx *gin.Context) {
	appointmentRequest := appointments.Appointment{}
	err := ctx.BindJSON(&appointmentRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if appointmentRequest.Date == "" ||
		appointmentRequest.Description == "" ||
		appointmentRequest.Patient.Dni == "" ||
		appointmentRequest.Dentist.License == "" {
		
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "appointment field missing, check sent JSON"})
		return
	}

	appointment, err := ah.appointmentsCreator.CreateAppointment(appointmentRequest)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, appointment)
}