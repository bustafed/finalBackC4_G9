package dentists

type Repository interface {
	GetDentistByID(id int) (Dentist, error)
	CreateDentist(d Dentist) (Dentist, error)
	UpdateDentistByID(id int, d Dentist) (Dentist, error)
	DeleteDentistByID(id int) error
}

type Service struct {
	repository Repository
}

func (s *Service) UpdateDentistByID(id int, dentist Dentist) (Dentist, error) {
	return s.repository.UpdateDentistByID(id, dentist)
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetDentistByID(id int) (Dentist, error) {
	return s.repository.GetDentistByID(id)
}

func (s *Service) CreateDentist(d Dentist) (Dentist, error) {
	return s.repository.CreateDentist(d)
}

func (s *Service) DeleteDentistByID(id int) error {
	return s.repository.DeleteDentistByID(id)
}
