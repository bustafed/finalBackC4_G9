package database

import (
	"database/sql"
	"fmt"
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

/*
func (s *SqlStore) Modify(id int, product products.Product) (products.Product, error) {
	query := fmt.Sprintf("UPDATE products SET name = '%s', quantity = %v, code_value = '%s',"+
		" is_published = %v, expiration = '%s', price = %v WHERE id = %v;", product.Name, product.Quantity,
		product.CodeValue, product.IsPublished, product.Expiration, product.Price, product.ID)
	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return products.Product{}, err
	}

	_, err = stmt.Exec()
	if err != nil {
		return products.Product{}, err
	}

	return product, nil
}*/
