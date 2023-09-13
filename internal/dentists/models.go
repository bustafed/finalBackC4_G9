package dentists

type Dentist struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	License string `json:"license"`
}
