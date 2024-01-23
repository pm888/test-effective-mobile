package models

type Person struct {
	ID         int
	Name       string `json:"name,omitempty" `
	Surname    string `json:"surname,omitempty" `
	Patronymic string `json:"patronymic,omitempty"`
	Age        int    `json:"age,omitempty"`
	Gender     string `json:"gender,omitempty"`
	Country    []PersonCountry
}
type PersonCountry struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type ChangePerson struct {
	ID         int    `json:"id,omitempty"`
	Name       string `json:"name,omitempty" `
	Surname    string `json:"surname,omitempty" `
	Patronymic string `json:"patronymic,omitempty"`
	Age        int    `json:"age,omitempty"`
	Gender     string `json:"gender,omitempty"`
}
