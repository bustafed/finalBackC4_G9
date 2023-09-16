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