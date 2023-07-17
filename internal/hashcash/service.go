package hashcash

const (
	zeroByte = '0'
)

type IService interface {
	AddIndicator(indicator int) error
	GetIndicator(indicator int) (int, error)
	RemoveIndicator(indicator int) error
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
