package hashcash

import (
	"context"
	"main/internal/shared/db"
	"strconv"
)

type IRepository interface {
	AddIndicator(ctx context.Context, indicator int64) error
	GetIndicator(ctx context.Context, indicator int64) (int64, error)
	RemoveIndicator(ctx context.Context, indicator int64) error
}

type Repository struct {
	db db.IDB
}

func NewRepository(db db.IDB) IRepository {
	return &Repository{
		db: db,
	}
}

func (r Repository) AddIndicator(ctx context.Context, indicator int64) error {
	err := r.db.Set(ctx, strconv.FormatInt(indicator, 10), "true")
	if err != nil {
		return ErrCouldNotAddIndicator
	}
	return nil
}

func (r Repository) GetIndicator(ctx context.Context, indicator int64) (int64, error) {
	value, err := r.db.Get(ctx, strconv.FormatInt(indicator, 10))
	if err != nil {
		return 0, ErrCouldNotGetIndicator
	}
	if value == nil {
		return 0, ErrIndicatorNotFound
	}
	return indicator, nil
}

func (r Repository) RemoveIndicator(ctx context.Context, indicator int64) error {
	r.db.Remove(ctx, strconv.FormatInt(indicator, 10))
	return nil
}
