package database

import (
	"fmt"
	"github.com/bustafed/finalBackC4_G9/internal/appointments"
	"log"
)

func (s *SqlStore) GetAppointmentByID(id int) (appointments.Appointment, error) {
	var appointmentReturn appointments.Appointment
	var dentistId int
	var patientId int

	query := fmt.Sprintf("SELECT * FROM appointments WHERE id = %d;", id)
	row := s.DB.QueryRow(query)
	err := row.Scan(&appointmentReturn.ID, &dentistId, &patientId,
		&appointmentReturn.Date, &appointmentReturn.Description)

	if err != nil {
		return appointments.Appointment{}, err
	}

	dentistFound, err := s.GetDentistByID(dentistId)
	if err != nil {
		return appointments.Appointment{}, err
	}
	appointmentReturn.Dentist = dentistFound

	patientFound, err := s.GetPatientByID(patientId)
	if err != nil {
		return appointments.Appointment{}, err
	}
	appointmentReturn.Patient = patientFound

	return appointmentReturn, nil
}

func (s *SqlStore) GetAppointmentByDni(dni string) ([]appointments.Appointment, error) {

	var appointmentsReturn []appointments.Appointment

	query := fmt.Sprintf("SELECT a.id, a.date, a.description, p.id, p.name, "+
		"p.surname, p.address, p.dni, p.registration_date, "+
		"d.id, d.name, d.surname, d.license "+
		"FROM appointments a "+
		"JOIN dentists d ON a.dentist_id = d.id "+
		"JOIN patients p ON a.patient_id = p.id "+
		"WHERE p.dni = '%s';", dni)

	rows, err := s.DB.Query(query)
	if err != nil {
		return []appointments.Appointment{}, err
	}

	for rows.Next() {
		var appointmentReturn appointments.Appointment
		err := rows.Scan(&appointmentReturn.ID, &appointmentReturn.Date, &appointmentReturn.Description,
			&appointmentReturn.Patient.ID, &appointmentReturn.Patient.Name, &appointmentReturn.Patient.Surname,
			&appointmentReturn.Patient.Address, &appointmentReturn.Patient.Dni, &appointmentReturn.Patient.RegistrationDate,
			&appointmentReturn.Dentist.ID, &appointmentReturn.Dentist.Name, &appointmentReturn.Dentist.Surname,
			&appointmentReturn.Dentist.License)
		if err != nil {
			return []appointments.Appointment{}, err
		}
		appointmentsReturn = append(appointmentsReturn, appointmentReturn)
	}

	return appointmentsReturn, nil
}

func (s *SqlStore) CreateAppointment(a appointments.Appointment) (appointments.Appointment, error) {

	dentistFound, err := s.GetDentistByLicense(a.Dentist.License)
	if err != nil {
		return appointments.Appointment{}, err
	}
	a.Dentist = dentistFound

	patientFound, err := s.GetPatientByDni(a.Patient.Dni)
	if err != nil {
		return appointments.Appointment{}, err
	}
	a.Patient = patientFound

	stmt, err := s.DB.Prepare("INSERT INTO appointments (dentist_id, patient_id, date, description) VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var insertedId int
	err = stmt.QueryRow(a.Dentist.ID, a.Patient.ID, a.Date, a.Description).Scan(&insertedId)

	if err != nil {
		return appointments.Appointment{}, err
	}

	a.ID = insertedId
	return a, nil
}

func (s *SqlStore) UpdateAppointmentByID(id int, a appointments.Appointment) (appointments.Appointment, error) {
	stmt, err := s.DB.Prepare("UPDATE appointments SET dentist_id = $1, patient_id = $2, date = $3, description = $4 WHERE id = $5")
	if err != nil {
		return appointments.Appointment{}, err

	}
	defer stmt.Close()

	_, err = stmt.Query(a.Dentist.ID, a.Patient.ID, a.Date, a.Description, id)
	if err != nil {
		return appointments.Appointment{}, err
	}
	a.ID = id
	return a, nil
}

func (s *SqlStore) DeleteAppointmentByID(id int) error {
	stmt, err := s.DB.Prepare("DELETE FROM appointments WHERE id = $1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Query(id)
	if err != nil {
		return err
	}
	return nil
}
