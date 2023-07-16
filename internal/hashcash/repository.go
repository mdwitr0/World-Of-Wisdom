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

func NewRepository(database db.IDB) IRepository {
	return &Repository{
		db: database,
	}
}

func (repo *Repository) AddIndicator(indicator int) error {
	return repo.db.Set(strconv.Itoa(indicator), "true")
}

func (repo *Repository) GetIndicator(indicator int) (int, error) {
	value, err := repo.db.Get(strconv.Itoa(indicator))

	if err != nil || value == nil {
		return 0, ErrIndicatorNotFound
	}

	return indicator, nil
}

func (repo *Repository) RemoveIndicator(indicator int) error {
	return repo.db.Remove(strconv.Itoa(indicator))
}
