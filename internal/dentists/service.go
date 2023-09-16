package dentists

type Repository interface {
	GetDentistByID(id int) (Dentist, error)
}

type Service struct {
	repository Repository
}

func (s *Service) ModifyByID(id int, dentist Dentist) (Dentist, error) {
	//TODO implement me
	panic("implement me")
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetDentistByID(id int) (Dentist, error) {
	return s.repository.GetDentistByID(id)
}
