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

// ユーザごとのLikeを受け付けるインメモリDBを定義
var likeDB = map[string]map[string]bool{}

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

// UpdateLike ...
func (handler *PetHandler) UpdateLike() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		// パスパラメータ "id" の値を取得
		id := c.Param("id")
		if id == "" {
			return c.JSON(http.StatusBadRequest, model.Response{
				Message: "No id param found",
			})
		}

		// Bindでリクエストの中身をiに詰める
		input := new(struct {
			UserId string `json:"user_id"`
			Value  bool   `json:"value"`
		})
		if err = c.Bind(input); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = handler.Interactor.UpdateLikeCount(&model.InputUpdateLikeRequest{
			PetId:  id,
			UserId: input.UserId,
			Value:  input.Value,
		})
		if err != nil {
			return utils.GetErrorMassage(c, "en", err)
		}

		return c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusOK,
			Message: "OK",
		})
	}
}

func (handler *PetHandler) Reservation() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		petId := c.Param("id")
		if petId == "" {
			return c.JSON(http.StatusBadRequest, model.Response{
				Message: "No id param found",
			})
		}

		// Bindでリクエストの中身をinputにつめる
		input := new(struct {
			UserId          string `json:"user_id"`
			Email           string `json:"email"`
			FullName        string `json:"full_name"`
			ReservationDate string `json:"reservation_date"`
		})
		if err = c.Bind(input); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

	}
}
