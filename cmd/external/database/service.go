package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bustafed/finalBackC4_G9/internal/dentists"
	"github.com/bustafed/finalBackC4_G9/internal/patients"
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
		return fmt.Errorf("Error deleting patient: %v", err)
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
