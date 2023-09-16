package patients

type Repository interface {
	GetPatientByID(id int) (Patient, error)
}

type Service struct {
	repository Repository
}

func (s *Service) ModifyByID(id int, patient Patient) (Patient, error) {
	//TODO implement me
	panic("implement me")
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetPatientByID(id int) (Patient, error) {
	return s.repository.GetPatientByID(id)
}
