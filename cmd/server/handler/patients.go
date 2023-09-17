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
	ModifyPatientByID(id int, patient patients.Patient) (patients.Patient, error)
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

// GetProductByID godoc
// @Summary      Gets a product by id
// @Description  Gets a product by id from the repository
// @Tags         products
// @Produce      json
// @Param        id path string true "ID"
// @Success      200 {object} patients.Patient
// @Router       /products/{id} [get]

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

	updatedPatient, err := ph.patientsCreator.ModifyPatientByID(id, patient)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "internal error"})
		return
	}
	ctx.JSON(http.StatusOK, updatedPatient)

}

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

	patient, err := ph.patientsCreator.ModifyPatientByID(id, patientRequest)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "internal error"})
		return
	}
	ctx.JSON(http.StatusOK, patient)
}

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
