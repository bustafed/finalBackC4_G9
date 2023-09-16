package handler

import (
	"github.com/bustafed/finalBackC4_G9/internal/dentists"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type DentistsGetter interface {
	GetDentistByID(id int) (dentists.Dentist, error)
}

type DentistCreator interface {
	ModifyByID(id int, dentist dentists.Dentist) (dentists.Dentist, error)
}

type DentistsHandler struct {
	dentistsGetter  DentistsGetter
	dentistsCreator DentistCreator
}

func NewDentistsHandler(getter DentistsGetter, creator DentistCreator) *DentistsHandler {
	return &DentistsHandler{
		dentistsGetter:  getter,
		dentistsCreator: creator,
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

func (dh *DentistsHandler) GetDentistByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	patient, err := dh.dentistsGetter.GetDentistByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, patient)
}

func (dh *DentistsHandler) PutDentist(context *gin.Context) {

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
