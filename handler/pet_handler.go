package handlers

import (
	"github.com/horsewin/echo-playground-v2/domain/model"
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/horsewin/echo-playground-v2/domain/repository"
	"github.com/horsewin/echo-playground-v2/interface/database"
	"github.com/horsewin/echo-playground-v2/usecase"
	"github.com/horsewin/echo-playground-v2/utils"
)

// PetHandler ...
type PetHandler struct {
	Interactor usecase.PetInteractor
}

// NewPetHandler ...
func NewPetHandler(sqlHandler database.SQLHandler) *PetHandler {
	return &PetHandler{
		Interactor: usecase.PetInteractor{
			PetRepository: &repository.PetRepository{
				SQLHandler: sqlHandler,
			},
		},
	}
}

// GetPets ...
func (handler *PetHandler) GetPets() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		gender := c.QueryParam("gender")
		res, err := handler.Interactor.GetPets(gender)
		if err != nil {
			return utils.GetErrorMassage(c, "en", err)
		}

		// resの中身をJSONにして返却
		resJSON := model.APIResponse{
			Data: res,
		}

		return c.JSON(http.StatusOK, resJSON)
	}
}

// CreateItem ...
func (handler *PetHandler) CreateItem() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		input := new(model.Pet)
		// Bindでリクエストの中身をiに詰める
		if err = c.Bind(input); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if input.Name == "" {
			return c.JSON(http.StatusBadRequest, model.Response{
				Message: "No name param found",
			})
		}

		resJSON, err := handler.Interactor.CreateItem(input)
		if err != nil {
			return utils.GetErrorMassage(c, "en", err)
		}

		return c.JSON(http.StatusOK, resJSON)
	}
}

// UpdateFavoriteAttr ...
func (handler *PetHandler) UpdateFavoriteAttr() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		input := new(model.Pet)
		if err = c.Bind(input); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = handler.Interactor.UpdateFavoriteAttr(*input)
		if err != nil {
			return utils.GetErrorMassage(c, "en", err)
		}

		return c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusOK,
			Message: "OK",
		})
	}
}
