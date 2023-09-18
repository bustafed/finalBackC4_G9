package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bustafed/finalBackC4_G9/internal/dentists"
	"github.com/gin-gonic/gin"
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

// GetDentistByID godoc
// @Summary      Gets a Dentist by id
// @Description  Gets a Dentist by id using the repository principal
// @Tags         Dentist
// @Produce      application/json
// @Param        PUBLIC-KEY header string true "publicKey"
// @Param        SECRET_KEY header string true "secretKey"
// @Param        id path string true "ID"
// @Success      200 {object} dentists.Dentist
// @Responses:
//
//	200: {object} dentists.Dentist (updated)
//	400: The id passed is in the wrong format
//	401: Either The PUBLIC-KEY or the SECRET-KEY or both are not correct
//	404: The dentist with the given id was not found
//	500: Internal error occured
//
// @Router       /dentists/{id} [get]
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

// UpdateDentistByID godoc
// @Summary      Updates a Dentist by id
// @Description  Updates a Dentist by id, you must send all of the dentist fields to process your request
// @Tags         Dentist
// @Produce      application/json
// @Param        PUBLIC-KEY header string true "publicKey"
// @Param        SECRET_KEY header string true "secretKey"
// @Param        id path string true "ID"
// @Success      200 {object} dentists.Dentist
// @Responses:
//
//	200: {object} dentists.Dentist (updated)
//	400: Either the id passed is in the wrong format or there are missing fields
//	401: Either The PUBLIC-KEY or the SECRET-KEY or both are not correct
//	404: The dentist with the given id was not found
//	500: Internal error occurred
//
// @Router       /dentists/{id} [put]
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

// UpdateDentistByID godoc
// @Summary      Updates a Dentist by id
// @Description  Updates a Dentist by id, you must send all of the dentist required fields to process your request
// @Tags         Dentist
// @Produce      application/json
// @Param        PUBLIC-KEY header string true "publicKey"
// @Param        SECRET_KEY header string true "secretKey"
// @Param        id path string true "ID"
// @Success      200 {object} dentists.Dentist
// @Responses:
//
//	200: {object} dentists.Dentist (updated)
//	400: Either the id passed is in the wrong format or there are missing fields
//	401: Either The PUBLIC-KEY or the SECRET-KEY or both are not correct
//	404: The dentist with the given id was not found
//	500: Internal error occurred
//
// @Router       /dentists/{id} [patch]
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

// CreateDentist godoc
// @Summary      Creates a Dentist
// @Description  Creates a Dentist, you must send the fields required to process your request they are name, surname, address, dni, and registration date.
// @Tags         Dentist
// @Produce      application/json
// @Param        PUBLIC-KEY header string true "publicKey"
// @Param        SECRET_KEY header string true "secretKey"
// @Param        Dentist body dentists.Dentist true "Create Dentist"
// @Success      200 {object} dentists.Dentist
// @Responses:
//
//	200: {object} dentists.Dentist
//	401: Either The PUBLIC-KEY or the SECRET-KEY or both are not correct
//	400: Either the request wasn't valid or all of the required fields weren't sent
//	500: Internal error occured
//
// @Router       /dentists [post]
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

// DeleteDentistByID godoc
// @Summary      Deletes a Dentist by id
// @Description  Deletes a Dentist by ID, be careful with this option!
// @Tags         Dentist
// @Produce      application/json
// @Param        PUBLIC-KEY header string true "publicKey"
// @Param        SECRET_KEY header string true "secretKey"
// @Param        id path string true "ID"
// @Success      200 {object} dentists.Dentist (updated)
// @Responses:
//
//	200: {object} dentists.Dentist (updated)
//	400: The id passed is in the wrong format
//	401: Either The PUBLIC-KEY or the SECRET-KEY or both are not correct
//	404: The Dentist with the given id was not found
//	500: Internal error occured
//
// @Router       /dentists/{id} [delete]
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
