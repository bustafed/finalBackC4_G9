package patients

type Patient struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Surname          string `json:"surname"`
	Address          string `json:"address"`
	Dni              string `json:"dni"`
	RegistrationDate string `json:"registration_date"`
}
