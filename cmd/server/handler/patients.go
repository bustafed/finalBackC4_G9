package handler

import (
	"github.com/bustafed/finalBackC4_G9/internal/patients"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PatientsGetter interface {
	GetPatientByID(id int) (patients.Patient, error)
}

type PatientCreator interface {
	ModifyByID(id int, product patients.Patient) (patients.Patient, error)
}

type PatientsHandler struct {
	patientsGetter  PatientsGetter
	patientsCreator PatientCreator
}

func NewPatientsHandler(getter PatientsGetter, creator PatientCreator) *PatientsHandler {
	return &PatientsHandler{
		patientsGetter:  getter,
		patientsCreator: creator,
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

func (ph *PatientsHandler) PutProduct(context *gin.Context) {

}

/*
func (ph *ProductsHandler) PutProduct(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	productRequest := patients.Patient{}
	err = ctx.BindJSON(&productRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, err := ph.productsCreator.ModifyByID(id, productRequest)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "internal error"})
		return
	}
	ctx.JSON(200, product)
}
*/
