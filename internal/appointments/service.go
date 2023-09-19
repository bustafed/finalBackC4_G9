package appointments

type Repository interface {
	GetAppointmentByID(id int) (Appointment, error)
	GetAppointmentByDni(dni string) ([]Appointment, error)
	CreateAppointment(d Appointment) (Appointment, error)
	UpdateAppointmentByID(id int, d Appointment) (Appointment, error)
	DeleteAppointmentByID(id int) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAppointmentByID(id int) (Appointment, error) {
	return s.repository.GetAppointmentByID(id)
}

func (s *Service) GetAppointmentByDni(dni string) ([]Appointment, error) {
	return s.repository.GetAppointmentByDni(dni)
}

func (s *Service) CreateAppointment(d Appointment) (Appointment, error) {
	return s.repository.CreateAppointment(d)
}

func (s *Service) UpdateAppointmentByID(id int, appointment Appointment) (Appointment, error) {
	return s.repository.UpdateAppointmentByID(id, appointment)
}

func (s *Service) DeleteAppointmentByID(id int) error {
	return s.repository.DeleteAppointmentByID(id)
}
