package hashcash

import (
	"fmt"
	"strconv"
	"time"
)

const (
	futuristicDays = 2 * 24 * time.Hour
	expiredDays    = 28 * 24 * time.Hour
	zeroByte       = '0'
)

type IService interface {
	AddIndicator(indicator int) error
	GetIndicator(indicator int) (int, error)
	RemoveIndicator(indicator int) error
	CheckStamp(stamp Stamp) bool
}

type Service struct {
	repository IRepository
}

func NewService(repository IRepository) IService {
	return &Service{
		repository: repository,
	}
}

func (r Service) AddIndicator(indicator int) error {
	return r.repository.AddIndicator(indicator)
}

func (r Service) GetIndicator(indicator int) (int, error) {
	return r.repository.GetIndicator(indicator)
}

func (r Service) RemoveIndicator(indicator int) error {
	return r.repository.RemoveIndicator(indicator)
}

func (r Service) CheckStamp(stamp Stamp) bool {
	if stamp.Date > time.Now().Add(futuristicDays).Unix() {
		return false
	}
	if stamp.Date < time.Now().Add(-expiredDays).Unix() {
		return false
	}

	if !stamp.IsSolved() {
		return false
	}

	v, err := strconv.Atoi(stamp.Rand)
	if err != nil {
		return false
	}

	ii, err := r.repository.GetIndicator(v)

	fmt.Println(ii)
	return err == nil
}
