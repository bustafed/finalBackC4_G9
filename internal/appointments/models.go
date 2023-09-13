package appointments

import (
	"github.com/bustafed/finalBackC4_G9/internal/dentists"
	"github.com/bustafed/finalBackC4_G9/internal/patients"
)

type Appointment struct {
	ID          int              `json:"id"`
	Patient     patients.Patient `json:"patient"`
	Dentist     dentists.Dentist `json:"dentist"`
	Date        string           `json:"date"`
	Description string           `json:"description"`
}
