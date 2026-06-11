package server

import (
	"context"

	"github.com/google/uuid"

	"github.com/LemuriiL/GophKeeperPassword/internal/model"
)

// ItemService работает с секретами
type ItemService struct {
	store *SQLite
}

// NewItemService создает сервис секретов
func NewItemService(store *SQLite) *ItemService {
	return &ItemService{store: store}
}

// Create создает секрет
func (s *ItemService) Create(ctx context.Context, userID int64, item model.Item) (model.Item, error) {
	if item.ID == "" {
		item.ID = uuid.NewString()
	}

	item.UserID = userID
	item.UpdatedAt = Now()

	err := s.store.CreateItem(ctx, item)
	return item, err
}

// Update обновляет секрет
func (s *ItemService) Update(ctx context.Context, userID int64, item model.Item) (model.Item, error) {
	item.UserID = userID
	item.UpdatedAt = Now()

	err := s.store.UpdateItem(ctx, item)
	return item, err
}

// Get получает секрет
func (s *ItemService) Get(ctx context.Context, userID int64, id string) (model.Item, error) {
	return s.store.GetItem(ctx, userID, id)
}

// List получает список секретов
func (s *ItemService) List(ctx context.Context, userID int64) ([]model.Item, error) {
	return s.store.ListItems(ctx, userID)
}

// Delete удаляет секрет
func (s *ItemService) Delete(ctx context.Context, userID int64, id string) error {
	return s.store.DeleteItem(ctx, userID, id)
}
