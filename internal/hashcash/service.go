package hashcash

import (
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

func NewService(hashRepo IRepository) IService {
	return &Service{
		repository: hashRepo,
	}
}

func (service *Service) AddIndicator(indicator int) error {
	return service.repository.AddIndicator(indicator)
}

func (service *Service) GetIndicator(indicator int) (int, error) {
	return service.repository.GetIndicator(indicator)
}

func (service *Service) RemoveIndicator(indicator int) error {
	return service.repository.RemoveIndicator(indicator)
}

func (service *Service) CheckStamp(stamp Stamp) bool {

	if stamp.Date > time.Now().Add(futuristicDays).Unix() || stamp.Date < time.Now().Add(-expiredDays).Unix() {
		return false
	}

	if !stamp.IsSolved() {
		return false
	}

	v, err := strconv.Atoi(stamp.Rand)
	if err != nil {
		return false
	}

	_, err = service.repository.GetIndicator(v)

	return err == nil
}
