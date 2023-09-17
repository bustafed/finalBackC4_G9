package patients

type Repository interface {
	GetPatientByID(id int) (Patient, error)
	UpdatePatientByID(id int, patient Patient) (Patient, error)
	CreatePatient(patient Patient) (Patient, error)
	DeletePatientByID(id int) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetPatientByID(id int) (Patient, error) {
	return s.repository.GetPatientByID(id)
}

func (s *Service) UpdatePatientByID(id int, patient Patient) (Patient, error) {
	return s.repository.UpdatePatientByID(id, patient)

}

func (s *Service) CreatePatient(patient Patient) (Patient, error) {
	return s.repository.CreatePatient(patient)
}

func (s *Service) DeletePatientByID(id int) error {
	return s.repository.DeletePatientByID(id)
}
