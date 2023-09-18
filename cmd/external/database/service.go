package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bustafed/finalBackC4_G9/internal/dentists"
	"github.com/bustafed/finalBackC4_G9/internal/patients"
	"github.com/bustafed/finalBackC4_G9/internal/appointments"
)

type SqlStore struct {
	*sql.DB
}

func NewDatabase(db *sql.DB) *SqlStore {
	return &SqlStore{db}
}

func (s *SqlStore) GetPatientByID(id int) (patients.Patient, error) {
	var patientReturn patients.Patient

	query := fmt.Sprintf("SELECT * FROM patients WHERE id = %d;", id)
	row := s.DB.QueryRow(query)
	err := row.Scan(&patientReturn.ID, &patientReturn.Name, &patientReturn.Surname, &patientReturn.Address,
		&patientReturn.Dni, &patientReturn.RegistrationDate)
	if err != nil {
		return patients.Patient{}, err
	}
	return patientReturn, nil
}

func (s *SqlStore) GetPatientByDni(dni string) (patients.Patient, error) {
	var patientReturn patients.Patient

	query := fmt.Sprintf("SELECT * FROM patients WHERE dni = '%s';", dni)
	row := s.DB.QueryRow(query)
	err := row.Scan(&patientReturn.ID, &patientReturn.Name, &patientReturn.Surname, &patientReturn.Address,
		&patientReturn.Dni, &patientReturn.RegistrationDate)
	if err != nil {
		return patients.Patient{}, err
	}
	return patientReturn, nil
}

func (s *SqlStore) UpdatePatientByID(id int, patient patients.Patient) (patients.Patient, error) {
	query := fmt.Sprintf("UPDATE patients SET name = '%s', surname = '%s', address = '%s',"+
		" dni = '%s', registration_date = '%s' WHERE id = %v;", patient.Name, patient.Surname,
		patient.Address, patient.Dni, patient.RegistrationDate, id)
	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return patients.Patient{}, err
	}

	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return patients.Patient{}, err
	}

	return patient, nil
}

func (s *SqlStore) CreatePatient(patient patients.Patient) (patients.Patient, error) {
	query := fmt.Sprintf("INSERT INTO patients (name, surname, address, dni, registration_date)"+
		" VALUES ('%s', '%s', '%s', '%s', '%s') RETURNING id;",
		patient.Name, patient.Surname, patient.Address, patient.Dni, patient.RegistrationDate)
	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return patients.Patient{}, err
	}
	defer stmt.Close()

	var insertedId int

	err = stmt.QueryRow().Scan(&insertedId)
	if err != nil {
		return patients.Patient{}, err
	}

	patient.ID = insertedId
	return patient, nil
}

func (s *SqlStore) DeletePatientByID(id int) error {

	query := fmt.Sprintf("DELETE FROM patients WHERE id = %v;", id)
	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("error deleting patient: %v", err)
	}

	return nil
}

func (s *SqlStore) GetDentistByID(id int) (dentists.Dentist, error) {
	var dentistReturn dentists.Dentist

	query := fmt.Sprintf("SELECT * FROM dentists WHERE id = %d;", id)
	row := s.DB.QueryRow(query)
	err := row.Scan(&dentistReturn.ID, &dentistReturn.Name, &dentistReturn.Surname, &dentistReturn.License)
	if err != nil {
		return dentists.Dentist{}, err
	}
	return dentistReturn, nil
}

func (s *SqlStore) GetDentistByLicense(license string) (dentists.Dentist, error) {
	var dentistReturn dentists.Dentist

	query := fmt.Sprintf("SELECT * FROM dentists WHERE license = '%s';", license)
	row := s.DB.QueryRow(query)
	err := row.Scan(&dentistReturn.ID, &dentistReturn.Name, &dentistReturn.Surname, &dentistReturn.License)
	if err != nil {
		return dentists.Dentist{}, err
	}
	return dentistReturn, nil
}

func (s *SqlStore) CreateDentist(d dentists.Dentist) (dentists.Dentist, error) {
	stmt, err := s.DB.Prepare("INSERT INTO dentists (name, surname, license) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var insertedId int
	err = stmt.QueryRow(d.Name, d.Surname, d.License).Scan(&insertedId)

	if err != nil {
		return dentists.Dentist{}, err
	}

	d.ID = insertedId
	return d, nil
}

func (s *SqlStore) UpdateDentistByID(id int, d dentists.Dentist) (dentists.Dentist, error) {
	stmt, err := s.DB.Prepare("UPDATE dentists SET name = $1, surname = $2, license = $3 WHERE id = $4")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Query(d.Name, d.Surname, d.License, id)
	if err != nil {
		return dentists.Dentist{}, err
	}
	d.ID = id
	return d, nil
}

func (s *SqlStore) DeleteDentistByID(id int) error {
	stmt, err := s.DB.Prepare("DELETE FROM dentists WHERE id = $1")
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
	
	query := fmt.Sprintf("SELECT a.id, a.date, a.description, p.id, p.name, " + 
			"p.surname, p.address, p.dni, p.registration_date, " +
			"d.id, d.name, d.surname, d.license " + 
			"FROM appointments a " +
			"JOIN dentists d ON a.dentist_id = d.id " +
			"JOIN patients p ON a.patient_id = p.id " +
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
