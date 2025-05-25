package handlers

import (
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/horsewin/echo-playground-v2/domain/model"
	"github.com/labstack/echo/v4"

	"github.com/horsewin/echo-playground-v2/domain/repository"
	"github.com/horsewin/echo-playground-v2/interface/database"
	"github.com/horsewin/echo-playground-v2/usecase"
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
			ReservationRepository: &repository.ReservationRepository{
				SQLHandler: sqlHandler,
			},
			FavoriteRepository: &repository.FavoriteRepository{
				SQLHandler: sqlHandler,
			},
		},
	}
}

// GetPets ...
func (handler *PetHandler) GetPets() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		// サブセグメントを作成
		ctx := c.Request().Context()
		subCtx, seg := xray.BeginSubsegment(ctx, "GetPets")
		if seg == nil {
			// セグメントが作成できない場合はログに記録して処理を続行
			c.Logger().Warn("Failed to begin subsegment: GetPets")

			filter := new(model.PetFilter)
			if err := c.Bind(filter); err != nil {
				return err
			}

			// Pass the context with X-Ray segment to the interactor
			res, err := handler.Interactor.GetPets(subCtx, filter)
			if err != nil {
				return err
			}

			// resの中身をJSONにして返却
			resJSON := model.APIResponse{
				Data: res,
			}

			return c.JSON(http.StatusOK, resJSON)
		}
		defer seg.Close(err)

		filter := new(model.PetFilter)
		if err := c.Bind(filter); err != nil {
			return err
		}

		// Pass the context with X-Ray segment to the interactor
		res, err := handler.Interactor.GetPets(ctx, filter)
		if err != nil {
			return err
		}

		// Add metadata to the segment
		if err := seg.AddMetadata("filter", filter); err != nil {
			c.Logger().Errorf("Failed to add filter metadata: %v", err)
		}
		if err := seg.AddMetadata("result_count", len(res)); err != nil {
			c.Logger().Errorf("Failed to add result_count metadata: %v", err)
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
		// サブセグメントを作成
		ctx := c.Request().Context()
		_, seg := xray.BeginSubsegment(ctx, "UpdateLike")
		if seg == nil {
			// セグメントが作成できない場合はログに記録して処理を続行
			c.Logger().Warn("Failed to begin subsegment: UpdateLike")

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

			// UseCaseの実行
			err = handler.Interactor.UpdateLikeCount(ctx, &model.InputUpdateLikeRequest{
				PetId:  id,
				UserId: input.UserId,
				Value:  input.Value,
			})

			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, model.Response{
				Code:    http.StatusOK,
				Message: "OK",
			})
		}
		defer seg.Close(err)

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

		// Add metadata to the segment
		if err := seg.AddMetadata("pet_id", id); err != nil {
			c.Logger().Errorf("Failed to add pet_id metadata: %v", err)
		}
		if err := seg.AddMetadata("user_id", input.UserId); err != nil {
			c.Logger().Errorf("Failed to add user_id metadata: %v", err)
		}
		if err := seg.AddMetadata("value", input.Value); err != nil {
			c.Logger().Errorf("Failed to add value metadata: %v", err)
		}

		// UseCaseの実行
		err = handler.Interactor.UpdateLikeCount(ctx, &model.InputUpdateLikeRequest{
			PetId:  id,
			UserId: input.UserId,
			Value:  input.Value,
		})

		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusOK,
			Message: "OK",
		})
	}
}

// Reservation ...
func (handler *PetHandler) Reservation() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		// サブセグメントを作成
		ctx := c.Request().Context()
		_, seg := xray.BeginSubsegment(ctx, "Reservation")
		if seg == nil {
			// セグメントが作成できない場合はログに記録して処理を続行
			c.Logger().Warn("Failed to begin subsegment: Reservation")

			// パスパラメータ "id" の値を取得
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

			// UseCaseの実行
			err = handler.Interactor.CreateReservation(ctx, &model.Reservation{
				PetId:           petId,
				UserId:          input.UserId,
				Email:           input.Email,
				FullName:        input.FullName,
				ReservationDate: input.ReservationDate,
			})

			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, model.Response{
				Code:    http.StatusOK,
				Message: "Created",
			})
		}
		defer seg.Close(err)

		// パスパラメータ "id" の値を取得
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

		// Add metadata to the segment
		if err := seg.AddMetadata("pet_id", petId); err != nil {
			c.Logger().Errorf("Failed to add pet_id metadata: %v", err)
		}
		if err := seg.AddMetadata("user_id", input.UserId); err != nil {
			c.Logger().Errorf("Failed to add user_id metadata: %v", err)
		}
		if err := seg.AddMetadata("reservation_date", input.ReservationDate); err != nil {
			c.Logger().Errorf("Failed to add reservation_date metadata: %v", err)
		}

		// UseCaseの実行
		err = handler.Interactor.CreateReservation(ctx, &model.Reservation{
			PetId:           petId,
			UserId:          input.UserId,
			Email:           input.Email,
			FullName:        input.FullName,
			ReservationDate: input.ReservationDate,
		})

		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusOK,
			Message: "Created",
		})
	}
}
