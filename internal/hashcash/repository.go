package hashcash

import (
	"main/internal/shared/db"
	"strconv"
)

type IRepository interface {
	AddIndicator(indicator int) error
	GetIndicator(indicator int) (int, error)
	RemoveIndicator(indicator int) error
}

type Repository struct {
	db db.IDB
}

func NewRepository(db db.IDB) IRepository {
	return &Repository{
		db: db,
	}
}

func (r Repository) AddIndicator(indicator int) error {
	err := r.db.Set(strconv.Itoa(indicator), "true")
	if err != nil {
		return ErrCouldNotAddIndicator
	}
	return nil
}

func (r Repository) GetIndicator(indicator int) (int, error) {
	value, err := r.db.Get(strconv.Itoa(indicator))
	if err != nil {
		return 0, ErrCouldNotGetIndicator
	}
	if value == nil {
		return 0, ErrIndicatorNotFound
	}
	return indicator, nil
}

func (r Repository) RemoveIndicator(indicator int) error {
	r.db.Remove(strconv.Itoa(indicator))
	return nil
}
