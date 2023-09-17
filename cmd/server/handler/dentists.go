package handler

import (
	"fmt"
	"github.com/bustafed/finalBackC4_G9/internal/dentists"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type DentistsGetter interface {
	GetDentistByID(id int) (dentists.Dentist, error)
}

type DentistCreator interface {
	CreateDentist(d dentists.Dentist) (dentists.Dentist, error)
	UpdateDentistByID(id int, dentist dentists.Dentist) (dentists.Dentist, error)
}

type DentistDeleter interface {
	DeleteDentistByID(id int) error
}

type DentistsHandler struct {
	dentistsGetter  DentistsGetter
	dentistsCreator DentistCreator
	dentistDeleter  DentistDeleter
}

func NewDentistsHandler(getter DentistsGetter, creator DentistCreator, deleter DentistDeleter) *DentistsHandler {
	return &DentistsHandler{
		dentistsGetter:  getter,
		dentistsCreator: creator,
		dentistDeleter:  deleter,
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
	dentist, err := dh.dentistsGetter.GetDentistByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, dentist)
}

func (dh *DentistsHandler) FullUpdateDentistByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	_, err = dh.dentistsGetter.GetDentistByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "dentist doesn't exist"})
		return
	}

	dentistRequest := dentists.Dentist{}
	err = ctx.BindJSON(&dentistRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if dentistRequest.Name == "" || dentistRequest.Surname == "" || dentistRequest.License == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "dentist field missing, check sent JSON"})
		return
	}

	dentist, err := dh.dentistsCreator.UpdateDentistByID(id, dentistRequest)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, dentist)
}

func (dh *DentistsHandler) UpdateDentistByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	dentist, err := dh.dentistsGetter.GetDentistByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "dentist doesn't exist"})
		return
	}

	dentistRequest := dentists.Dentist{}
	err = ctx.BindJSON(&dentistRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if dentistRequest.Name == "" && dentistRequest.Surname == "" && dentistRequest.License == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "dentist field missing, check sent JSON"})
		return
	}

	if dentistRequest.Name == "" {
		dentistRequest.Name = dentist.Name
	}
	if dentistRequest.Surname == "" {
		dentistRequest.Surname = dentist.Surname
	}
	if dentistRequest.License == "" {
		dentistRequest.License = dentist.License
	}

	dentist, err = dh.dentistsCreator.UpdateDentistByID(id, dentistRequest)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, dentist)
}

func (dh *DentistsHandler) CreateDentist(ctx *gin.Context) {
	dentistRequest := dentists.Dentist{}
	err := ctx.BindJSON(&dentistRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if dentistRequest.Name == "" || dentistRequest.Surname == "" || dentistRequest.License == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "dentist field missing, check sent JSON"})
		return
	}

	dentist, err := dh.dentistsCreator.CreateDentist(dentistRequest)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, dentist)
}

func (dh *DentistsHandler) DeleteDentistByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	_, err = dh.dentistsGetter.GetDentistByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "dentist doesn't exist"})
		return
	}

	err = dh.dentistDeleter.DeleteDentistByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, fmt.Sprintf("dentist with ID: %v deleted", id))
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
