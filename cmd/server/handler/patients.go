package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bustafed/finalBackC4_G9/internal/patients"
	"github.com/gin-gonic/gin"
)

type PatientsGetter interface {
	GetPatientByID(id int) (patients.Patient, error)
}

type PatientCreator interface {
	UpdatePatientByID(id int, patient patients.Patient) (patients.Patient, error)
	CreatePatient(patient patients.Patient) (patients.Patient, error)
}
type PatientDeleter interface {
	DeletePatientByID(id int) error
}

type PatientsHandler struct {
	patientsGetter  PatientsGetter
	patientsCreator PatientCreator
	patientsDeleter PatientDeleter
}

func NewPatientsHandler(getter PatientsGetter, creator PatientCreator, deleter PatientDeleter) *PatientsHandler {
	return &PatientsHandler{
		patientsGetter:  getter,
		patientsCreator: creator,
		patientsDeleter: deleter,
	}
}

// GetPatientByID godoc
// @Summary      Gets a Patient by id
// @Description  Gets a Patient by id using the repository principal
// @Tags         Patient
// @Produce      application/json
// @Param        PUBLIC-KEY header string true "publicKey"
// @Param        SECRET_KEY header string true "secretKey"
// @Param        id path string true "ID"
// @Success      200 {object} patients.Patient
// @Responses:
//
//	200: {object} patients.Patient (updated)
//	400: Your the id passed is in the wrong format
//	401: Either The PUBLIC-KEY or the SECRET-KEY or both are not correct
//	404: The patient with the given id was not found
//	500: Internal error occured
//
// @Router       /patients/{id} [get]
func (ph *PatientsHandler) GetPatientByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	patient, err := ph.patientsGetter.GetPatientByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, patient)
}

// UpdatePatientByID godoc
// @Summary      Updates a Patient by id
// @Description  Updates a Patient by ID, you may be noticed is not required to send data in all of the fields
// @Tags         Patient
// @Produce      application/json
// @Param        PUBLIC-KEY header string true "publicKey"
// @Param        SECRET_KEY header string true "secretKey"
// @Param        id path string true "ID"
// @Success      200 {object} patients.Patient (updated)
// @Responses:
//
//	200: {object} patients.Patient (updated)
//	400: The id passed is in the wrong format
//	401: Either The PUBLIC-KEY or the SECRET-KEY or both are not correct
//	404: The patient with the given id was not found
//	500: Internal error occured
//
// @Router       /patients/{id} [patch]
func (ph *PatientsHandler) ModifyPatientByProperty(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	patient, err := ph.patientsGetter.GetPatientByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	patientRequest := patients.Patient{}
	err = ctx.BindJSON(&patientRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if patientRequest.Name != "" {
		patient.Name = patientRequest.Name
	}
	if patientRequest.Surname != "" {
		patient.Surname = patientRequest.Surname
	}
	if patientRequest.Address != "" {
		patient.Address = patientRequest.Address
	}
	if patientRequest.Dni != "" {
		patient.Dni = patientRequest.Dni
	}
	if patientRequest.RegistrationDate != "" {
		patient.RegistrationDate = patientRequest.RegistrationDate
	}

	updatedPatient, err := ph.patientsCreator.UpdatePatientByID(id, patient)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "internal error"})
		return
	}
	ctx.JSON(http.StatusOK, updatedPatient)

}

// UpdatePatientByID godoc
// @Summary      Updates a Patient by id
// @Description  Updates a Patient by ID, you must send all of the patient fields to process your request
// @Tags         Patient
// @Produce      application/json
// @Param        PUBLIC-KEY header string true "publicKey"
// @Param        SECRET_KEY header string true "secretKey"
// @Param        id path string true "ID"
// @Success      200 {object} patients.Patient (updated)
// @Responses:
//
//	200: {object} patients.Patient (updated)
//	400: Either the request wasn't valid or all of the required fields weren't sent
//	401: Either The PUBLIC-KEY or the SECRET-KEY or both are not correct
//	404: The patient with the given id was not found
//	500: Internal error occured
//
// @Router       /patients/{id} [put]
func (ph *PatientsHandler) PutPatient(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	patientRequest := patients.Patient{}
	err = ctx.ShouldBindJSON(&patientRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if patientRequest.Name == "" || patientRequest.Surname == "" || patientRequest.Address == "" || patientRequest.Dni == "" || patientRequest.RegistrationDate == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "required fields are missing"})
		return
	}
	patientRequest.ID = id

	patient, err := ph.patientsCreator.UpdatePatientByID(id, patientRequest)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "internal error"})
		return
	}
	ctx.JSON(http.StatusOK, patient)
}

// CreatePatient godoc
// @Summary      Creates a Patient
// @Description  Creates a Patient, you must send the fields required to process your request they are name, surname, address, dni, and registration date.
// @Tags         Patient
// @Produce      application/json
// @Param        PUBLIC-KEY header string true "publicKey"
// @Param        SECRET_KEY header string true "secretKey"
// @Param        Patient body patients.Patient true "create Patient"
// @Success      200 {object} patients.Patient
// @Responses:
//
//	200: {object} patients.Patient
//	400: Either the request wasn't valid or all of the required fields weren't sent
//	401: Either The PUBLIC-KEY or the SECRET-KEY or both are not correct
//	500: Internal error occured
//
// @Router       /patients [post]
func (ph *PatientsHandler) CreatePatient(ctx *gin.Context) {
	patientRequest := patients.Patient{}
	err := ctx.ShouldBindJSON(&patientRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if patientRequest.Name == "" || patientRequest.Surname == "" || patientRequest.Address == "" || patientRequest.Dni == "" || patientRequest.RegistrationDate == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "required fields are missing"})
		return
	}

	patient, err := ph.patientsCreator.CreatePatient(patientRequest)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "internal error"})
		return
	}
	ctx.JSON(http.StatusOK, patient)
}

// DeletePatientByID godoc
// @Summary      Deletes a Patient by id
// @Description  Deletes a Patient by ID, be careful with this option!
// @Tags         Patient
// @Produce      application/json
// @Param        PUBLIC-KEY header string true "publicKey"
// @Param        SECRET_KEY header string true "secretKey"
// @Param        id path string true "ID"
// @Success      200 {object} patients.Patient (updated)
// @Responses:
//
//	200: {object} patients.Patient (updated)
//	400: The id passed is in the wrong format
//	401: Either The PUBLIC-KEY or the SECRET-KEY or both are not correct
//	404: The patient with the given id was not found
//	500: Internal error occured
//
// @Router       /patients/{id} [delete]
func (ph *PatientsHandler) DeletePatientByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	err2 := ph.patientsDeleter.DeletePatientByID(id)
	if err2 != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err2.Error()})
		return
	}

	ctx.JSON(http.StatusOK, fmt.Sprintf("patient with id: %v deleted", id))
}
