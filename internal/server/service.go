package server

import (
	"context"
	"github.com/google/uuid"

	"github.com/LemuriiL/GophKeeperPassword/internal/model"
)

// ItemService отвечает за секреты.
type ItemService struct {
	store *SQLite
}

// NewItemService создаёт сервис секретов.
func NewItemService(store *SQLite) *ItemService {
	return &ItemService{store: store}
}

// Save сохраняет секрет пользователя.
func (s *ItemService) Save(ctx context.Context, userID int64, item model.Item) (model.Item, error) {
	if item.ID == "" {
		item.ID = uuid.NewString()
	}
	item.UserID = userID
	item.UpdatedAt = Now()

	err := s.store.UpsertItem(ctx, item)
	return item, err
}

// List возвращает секреты пользователя.
func (s *ItemService) List(ctx context.Context, userID int64) ([]model.Item, error) {
	return s.store.ListItems(ctx, userID)
}

// Delete удаляет секрет пользователя.
func (s *ItemService) Delete(ctx context.Context, userID int64, id string) error {
	return s.store.DeleteItem(ctx, userID, id)
}
