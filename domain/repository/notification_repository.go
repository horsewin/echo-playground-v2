package repository

import (
	"context"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/horsewin/echo-playground-v2/domain/model"
	"github.com/horsewin/echo-playground-v2/interface/database"
	"github.com/horsewin/echo-playground-v2/utils"
)

// NotificationRepositoryInterface ...
type NotificationRepositoryInterface interface {
	Find(ctx context.Context, id string) (account model.Notifications, err error)
	FindAll(ctx context.Context) (notifications model.Notifications, err error)
	Count(ctx context.Context, query string, args map[string]interface{}) (data model.NotificationCount, err error)
	Update(ctx context.Context, in map[string]interface{}, query string, args map[string]interface{}) (err error)
}

// NotificationRepository ....
type NotificationRepository struct {
	database.SQLHandler
}

const NotificationTable = "notifications"

// Find ...
func (repo *NotificationRepository) Find(ctx context.Context, id string) (notifications model.Notifications, err error) {
	// サブセグメントを作成
	_, seg := xray.BeginSubsegment(ctx, "NotificationRepository.Find")
	defer seg.Close(err)

	// メタデータを追加
	if err := seg.AddMetadata("id", id); err != nil {
		utils.LogError("Failed to add id metadata: %v", err)
	}

	whereClause := "id = :id"
	whereArgs := map[string]interface{}{"id": id}
	err = repo.SQLHandler.Where(ctx, &notifications.Data, NotificationTable, whereClause, whereArgs)
	return
}

// FindAll ...
func (repo *NotificationRepository) FindAll(ctx context.Context) (notifications model.Notifications, err error) {
	// サブセグメントを作成
	_, seg := xray.BeginSubsegment(ctx, "NotificationRepository.FindAll")
	defer seg.Close(err)

	err = repo.SQLHandler.Scan(ctx, &notifications.Data, NotificationTable, "id desc")
	return notifications, err
}

// Count ...
func (repo *NotificationRepository) Count(ctx context.Context, query string, args map[string]interface{}) (data model.NotificationCount, err error) {
	// サブセグメントを作成
	_, seg := xray.BeginSubsegment(ctx, "NotificationRepository.Count")
	defer seg.Close(err)

	// メタデータを追加
	if err := seg.AddMetadata("query", query); err != nil {
		utils.LogError("Failed to add query metadata: %v", err)
	}
	if err := seg.AddMetadata("args", args); err != nil {
		utils.LogError("Failed to add args metadata: %v", err)
	}

	var count int
	err = repo.SQLHandler.Count(ctx, &count, NotificationTable, query, args)
	return model.NotificationCount{Data: count}, err
}

// Update ...
func (repo *NotificationRepository) Update(ctx context.Context, in map[string]interface{}, query string, args map[string]interface{}) (err error) {
	// サブセグメントを作成
	_, seg := xray.BeginSubsegment(ctx, "NotificationRepository.Update")
	defer seg.Close(err)

	// メタデータを追加
	if err := seg.AddMetadata("query", query); err != nil {
		utils.LogError("Failed to add query metadata: %v", err)
	}
	if err := seg.AddMetadata("args", args); err != nil {
		utils.LogError("Failed to add args metadata: %v", err)
	}
	if err := seg.AddMetadata("input", in); err != nil {
		utils.LogError("Failed to add input metadata: %v", err)
	}

	//err = repo.SQLHandler.Update(ctx, in, NotificationTable, query)
	return
}
