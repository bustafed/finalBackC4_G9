package database

import (
	"fmt"
	"github.com/bustafed/finalBackC4_G9/internal/patients"
)

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
