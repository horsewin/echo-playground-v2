package usecase

import (
	"context"
	"math/rand"
	"time"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/horsewin/echo-playground-v2/domain/model"
	"github.com/horsewin/echo-playground-v2/domain/repository"
	"github.com/horsewin/echo-playground-v2/utils"
	"github.com/jinzhu/copier"
)

// PetInteractor ...
type PetInteractor struct {
	PetRepository         repository.PetRepositoryInterface
	ReservationRepository repository.ReservationRepositoryInterface
	FavoriteRepository    repository.FavoriteRepositoryInterface
}

// GetPets ...
func (interactor *PetInteractor) GetPets(ctx context.Context, filter *model.PetFilter) (pets []model.Pet, err error) {
	// サブセグメントを作成
	subCtx, seg := xray.BeginSubsegment(ctx, "PetInteractor.GetPets")
	defer seg.Close(err)

	// メタデータを追加
	if err := seg.AddMetadata("filter", filter); err != nil {
		// エラーはログに記録するだけで処理は続行
		utils.LogError("Failed to add filter metadata: %v", err)
	}

	// 性能テスト用：3回に1回だけCPU負荷を発生させる
	if rand.Intn(3) == 0 {
		induceCpuLoad()
	}

	// 性能テスト用：3回に1回だけレイテンシを発生させる
	if rand.Intn(3) == 0 {
		induceLatency()
	}

	// Repository層からデータを取得
	// サブセグメントを作成
	_app, err := interactor.PetRepository.Find(subCtx, filter)
	if err != nil {
		err = utils.ConvertErrorMassage(ctx, "10001E", err)
		return
	}

	// ドメインモデルに変換
	for _, p := range _app.Data {
		// 予約数を取得
		// Note: 意図的にN+1問題を起こしている箇所。X-Rayで確認するため。
		reservationCount, countErr := interactor.ReservationRepository.GetCountByPetID(subCtx, p.ID)
		if countErr != nil {
			utils.LogError("Failed to get reservation count: %v", countErr)
			// エラーが発生しても処理は継続
			reservationCount = 0
		}

		pets = append(pets, model.Pet{
			ID:       p.ID,
			Name:     p.Name,
			Breed:    p.Breed,
			Gender:   p.Gender,
			Price:    p.Price,
			ImageURL: p.ImageURL,
			Likes:    p.Likes,
			Shop: model.Shop{
				Name:     p.ShopName,
				Location: p.ShopLocation,
			},
			BirthDate:        p.BirthDate,
			ReferenceNumber:  p.ReferenceNumber,
			Tags:             p.Tags,
			ReservationCount: reservationCount,
		})
	}

	// 結果のメタデータを追加
	if err := seg.AddMetadata("result_count", len(pets)); err != nil {
		utils.LogError("Failed to add result_count metadata: %v", err)
	}

	return pets, nil
}

// UpdateLikeCount ...
func (interactor *PetInteractor) UpdateLikeCount(ctx context.Context, input *model.InputUpdateLikeRequest) (err error) {
	// サブセグメントを作成
	ctx, seg := xray.BeginSubsegment(ctx, "PetInteractor.UpdateLikeCount")
	defer seg.Close(err)

	// メタデータを追加
	if err := seg.AddMetadata("input", input); err != nil {
		utils.LogError("Failed to add input metadata: %v", err)
	}

	// like状態を取得
	favMap, err := interactor.FavoriteRepository.FindByUserId(ctx, input.UserId)
	if err != nil {
		err = utils.ConvertErrorMassage(ctx, "10001E", err)
		return
	}
	if favMap[input.PetId].Value == input.Value {
		err = utils.ConvertErrorMassage(ctx, "00001I", nil)
		return
	}

	// 現在のペットモデルを取得
	petData, err := interactor.PetRepository.Find(ctx, &model.PetFilter{ID: input.PetId})
	if err != nil {
		err = utils.ConvertErrorMassage(ctx, "10001E", err)
		return
	}
	var pet model.Pet
	for _, p := range petData.Data {
		if p.ID == input.PetId {
			err = copier.Copy(&pet, &p)
			if err != nil {
				err = utils.ConvertErrorMassage(ctx, "10002E", err)
				return
			}

			// マッピングしきれない構造は手動でコピー
			pet.Shop = model.Shop{
				Name:     p.ShopName,
				Location: p.ShopLocation,
			}
		}
	}

	// Like数を更新
	if input.Value {
		pet.Likes = pet.Likes + 1
	} else {
		pet.Likes = pet.Likes - 1
	}

	// TODO: トランザクションを使って更新処理を行う
	err = interactor.PetRepository.Update(ctx, &pet)
	if err != nil {
		err = utils.ConvertErrorMassage(ctx, "10003E", err)
		return
	}

	if input.Value {
		err = interactor.FavoriteRepository.Create(ctx, &model.Favorite{
			PetId:  input.PetId,
			UserId: input.UserId,
			Value:  input.Value,
		})
		if err != nil {
			err = utils.ConvertErrorMassage(ctx, "10004E", err)
			return
		}
	} else {
		err = interactor.FavoriteRepository.Delete(ctx, &model.Favorite{
			PetId:  input.PetId,
			UserId: input.UserId,
		})
		if err != nil {
			err = utils.ConvertErrorMassage(ctx, "10005E", err)
			return
		}
	}

	return
}

// CreateReservation ...
func (interactor *PetInteractor) CreateReservation(ctx context.Context, input *model.Reservation) (err error) {
	// サブセグメントを作成
	ctx, seg := xray.BeginSubsegment(ctx, "PetInteractor.CreateReservation")
	defer seg.Close(err)

	// メタデータを追加
	if err := seg.AddMetadata("input", input); err != nil {
		utils.LogError("Failed to add input metadata: %v", err)
	}

	// サブセグメントを作成
	err = interactor.ReservationRepository.Create(ctx, input)
	if err != nil {
		err = utils.ConvertErrorMassage(ctx, "10003E", err)
		return
	}
	return
}

// induceCpuLoad ... 意図的にテスト用のCPU負荷を発生させる関数
func induceCpuLoad() {
	t := time.NewTimer(3 * time.Second)

	go func() {
		//nolint:staticcheck // 意図的に無限ループを作成してCPU負荷をシミュレートするためのコード
		for {
		}
	}()
	<-t.C
	t.Stop()
}

// induceLatency ... 意図的にテスト用のレイテンシを発生させる関数
func induceLatency() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	milliseconds := r.Intn(500) + 500
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)
}
