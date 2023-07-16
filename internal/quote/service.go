package quote

import "fmt"

type IService interface {
	GetRandomQuote() string
}

type Service struct {
	repository IRepository
}

func NewService(repository IRepository) IService {
	return &Service{
		repository: repository,
	}
}

func (r *Service) GetRandomQuote() string {
	q := r.repository.GetRandomQuote()

	return fmt.Sprintf("%s %s", q.Text, q.Author)
}
