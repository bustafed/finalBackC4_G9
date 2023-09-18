DROP TABLE IF EXISTS patients;

CREATE TABLE patients
(
    id                SERIAL PRIMARY KEY,
    name              VARCHAR(128) NOT NULL,
    surname           VARCHAR(128) NOT NULL,
    address           VARCHAR(128) NOT NULL,
    dni               VARCHAR(128) NOT NULL,
    registration_date VARCHAR(128) NOT NULL
);

INSERT INTO patients (name, surname, address, dni, registration_date)
VALUES ('Patient 1', 'Surname 1', 'Avenida Siempreviva 752', '11222333', '2023-01-01');

DROP TABLE IF EXISTS dentists;

CREATE TABLE dentists
(
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(128) NOT NULL,
    surname VARCHAR(128) NOT NULL,
    license VARCHAR(128) NOT NULL
);

INSERT INTO dentists (name, surname, license)
VALUES ('Dentist 1', 'Surname 1', 'M1232121');

DROP TABLE IF EXISTS appointments;

CREATE TABLE appointments
(
    id SERIAL PRIMARY KEY,
    dentist_id INT NOT NULL,
    patient_id INT NOT NULL,
    date VARCHAR(128) NOT NULL,
    description VARCHAR(128) NOT NULL,
    CONSTRAINT FK_appointment_dentist FOREIGN KEY(dentist_id)
    REFERENCES dentists(id)
    ON DELETE CASCADE,
    CONSTRAINT FK_appointment_patient FOREIGN KEY(patient_id)
    REFERENCES patients(id)
    ON DELETE CASCADE
);

INSERT INTO appointments (dentist_id, patient_id, date, description)
VALUES (1, 1, '10-10-2023 16:00', 'Chequeo anual');

