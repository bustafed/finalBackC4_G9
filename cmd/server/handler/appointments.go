package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bustafed/finalBackC4_G9/internal/appointments"
	"github.com/gin-gonic/gin"
)

type AppointmentsGetter interface {
	GetAppointmentByID(id int) (appointments.Appointment, error)
	GetAppointmentByDni(dni string) ([]appointments.Appointment, error)
}

type AppointmentCreator interface {
	UpdateAppointmentByID(id int, appointment appointments.Appointment) (appointments.Appointment, error)
	CreateAppointment(appointment appointments.Appointment) (appointments.Appointment, error)
}
type AppointmentDeleter interface {
	DeleteAppointmentByID(id int) error
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

// GetAppointmentByID godoc
// @Summary      Gets an Appointment by id
// @Description  Gets an Appointment by id using the repository principal
// @Tags         Appointment
// @Produce      application/json
// @Param        PUBLIC-KEY header string true "publicKey"
// @Param        SECRET_KEY header string true "secretKey"
// @Param        id path string true "ID"
// @Success      200 {object} appointments.Appointment
// @Responses:
//
//	200: {object} appointments.Appointment (updated)
//	400: The id passed is in the wrong format
//	401: Either The PUBLIC-KEY or the SECRET-KEY or both are not correct
//	404: The patient with the given id was not found
//	500: Internal error occured
//
// @Router       /appointments/{id} [get]
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

// GetAppointmentByDNI godoc
// @Summary      Gets all appointments by dni
// @Description  Gets all appointments if any by patient dni
// @Tags         Appointment
// @Produce      application/json
// @Param        PUBLIC-KEY header string true "publicKey"
// @Param        SECRET_KEY header string true "secretKey"
// @Param        DNI path string true "DNI"
// @Success      200 {object} appointments.Appointment
// @Responses:
//
//	200: {object} appointments.Appointment (updated)
//	400: The id passed is in the wrong format
//	401: Either The PUBLIC-KEY or the SECRET-KEY or both are not correct
//	404: The patient with the given dni was not found
//	500: Internal error occured
//
// @Router       /appointments/{id} [get]
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

// CreateAppointment godoc
// @Summary      Creates an Appointment
// @Description  Creates an Appointment, you must send the fields required to process your request Patient, Dentist, Date, Description
// @Tags         Appointment
// @Produce      application/json
// @Param        PUBLIC-KEY header string true "publicKey"
// @Param        SECRET_KEY header string true "secretKey"
// @Param        Appointment body appointments.Appointment true "Create Appointment"
// @Success      200 {object} appointments.Appointment
// @Responses:
//
//	200: {object} appointments.Appointment
//	401: Either The PUBLIC-KEY or the SECRET-KEY or both are not correct
//	400: Either the request wasn't valid or all of the required fields weren't sent
//	404: Either the dentist or the patient were not found
//	500: Internal error occured
//
// @Router       /appointments [post]
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, appointment)
}

// UpdateAppointmentByID godoc
// @Summary      Updates an Appointment by id
// @Description  Updates an Appointment by id, you must send all of the appointment fields to process your request
// @Tags         Appointment
// @Produce      application/json
// @Param        PUBLIC-KEY header string true "publicKey"
// @Param        SECRET_KEY header string true "secretKey"
// @Param        id path string true "ID"
// @Param        Appointment body appointments.Appointment true "Update Appointment"
// @Success      200 {object} appointments.Appointment
// @Responses:
//
//	200: {object} appointments.Appointment (updated)
//	400: Either the id passed is in the wrong format or there are missing fields
//	401: Either The PUBLIC-KEY or the SECRET-KEY or both are not correct
//	404: The appointment with the given id was not found
//	500: Internal error occurred
//
// @Router       /appointments/{id} [put]
func (ah *AppointmentsHandler) FullUpdateAppointmentByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	_, err = ah.appointmentsGetter.GetAppointmentByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "appointment doesn't exist"})
		return
	}

	appointmentRequest := appointments.Appointment{}
	err = ctx.BindJSON(&appointmentRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if appointmentRequest.Dentist.ID == 0 || appointmentRequest.Patient.ID == 0 || appointmentRequest.Date == "" || appointmentRequest.Description == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "appointment field missing, check sent JSON"})
		return
	}

	appointment, err := ah.appointmentsCreator.UpdateAppointmentByID(id, appointmentRequest)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, appointment)
}

// UpdateAppointmentByID godoc
// @Summary      Updates an Appointment by id
// @Description  Updates an Appointment by id, you can send only the appointment fields you need to change
// @Tags         Appointment
// @Produce      application/json
// @Param        PUBLIC-KEY header string true "publicKey"
// @Param        SECRET_KEY header string true "secretKey"
// @Param        id path string true "ID"
// @Param        Appointment body appointments.Appointment true "Update Appointment"
// @Success      200 {object} appointments.Appointment
// @Responses:
//
//	200: {object} appointments.Appointment (updated)
//	400: Either the id passed is in the wrong format or there are missing fields
//	401: Either The PUBLIC-KEY or the SECRET-KEY or both are not correct
//	404: The appointment with the given id was not found
//	500: Internal error occurred
//
// @Router       /appointments/{id} [patch]
func (ah *AppointmentsHandler) UpdateAppointmentByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	appointment, err := ah.appointmentsGetter.GetAppointmentByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "appointment doesn't exist"})
		return
	}

	appointmentRequest := appointments.Appointment{}
	err = ctx.BindJSON(&appointmentRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if appointmentRequest.Dentist.ID == 0 && appointmentRequest.Patient.ID == 0 && appointmentRequest.Date == "" && appointmentRequest.Description == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "appointment field missing, check sent JSON"})
		return
	}

	if appointmentRequest.Dentist.ID == 0 {
		appointmentRequest.Dentist.ID = appointment.Dentist.ID
	}
	if appointmentRequest.Patient.ID == 0 {
		appointmentRequest.Patient.ID = appointment.Patient.ID
	}
	if appointmentRequest.Date == "" {
		appointmentRequest.Date = appointment.Date
	}
	if appointmentRequest.Description == "" {
		appointmentRequest.Description = appointment.Description
	}

	appointment, err = ah.appointmentsCreator.UpdateAppointmentByID(id, appointmentRequest)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, appointment)
}

// DeleteAppointmentByID godoc
// @Summary      Deletes an Appointment by id
// @Description  Deletes an Appointment by id, be careful with this action.
// @Tags         Appointment
// @Produce      application/json
// @Param        PUBLIC-KEY header string true "publicKey"
// @Param        SECRET_KEY header string true "secretKey"
// @Param        id path string true "ID"
// @Param        Appointment body appointments.Appointment true "Delete Appointment"
// @Success      200 {object} appointments.Appointment
// @Responses:
//
//	200: {object} appointments.Appointment (updated)
//	400: Either the id passed is in the wrong format or there are missing fields
//	401: Either The PUBLIC-KEY or the SECRET-KEY or both are not correct
//	404: The appointment with the given id was not found or the dentist associated or the patient associated
//	500: Internal error occurred
//
// @Router       /appointments/{id} [delete]
func (ah *AppointmentsHandler) DeleteAppointmentByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	_, err = ah.appointmentsGetter.GetAppointmentByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "appointment doesn't exist"})
		return
	}

	err = ah.appointmentDeleter.DeleteAppointmentByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, fmt.Sprintf("appointment with ID: %v deleted", id))
}
