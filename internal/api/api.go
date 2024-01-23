// internal/api/api.go

package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"test_effective_mobile/test-effective-mobile/internal/database"
	"test_effective_mobile/test-effective-mobile/internal/models"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Response struct {
	Age         int               `json:"age"`
	Gender      string            `json:"gender"`
	Nationality []CountryResponse `json:"country"`
}
type CountryResponse struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

func handleError(w http.ResponseWriter, errMsg string, statusCode int) {
	http.Error(w, errMsg, statusCode)
}

func EnrichPersonInfo(person *models.Person) error {
	name := person.Name
	var response Response

	agifyURL := "https://api.agify.io/?name=" + name
	genderURL := "https://api.genderize.io/?name=" + name
	nationalityURL := "https://api.nationalize.io/?name=" + name

	var apiURLs = map[string]string{
		"age":         agifyURL,
		"gender":      genderURL,
		"nationality": nationalityURL,
	}

	for v, url := range apiURLs {
		resp, err := http.Get(url)
		if err != nil {
			return errors.Wrap(err, "error in http.Get")
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			resp.Body.Close()
			return errors.Wrap(err, "error in io.ReadAll")
		}
		switch {
		case v == "age":
			if err = json.Unmarshal(body, &response); err != nil {
				return fmt.Errorf("error in json.Unmarshal: %v", err)
			}
			person.Age = response.Age

		case v == "gender":
			if err = json.Unmarshal(body, &response); err != nil {
				return fmt.Errorf("error in json.Unmarshal: %v", err)
			}
			person.Gender = response.Gender
		case v == "nationality":
			if err = json.Unmarshal(body, &response); err != nil {
				return fmt.Errorf("error in json.Unmarshal: %v", err)
			}
			if len(response.Nationality) > 0 {
				for _, countryInfo := range response.Nationality {
					country := models.PersonCountry{
						CountryID:   countryInfo.CountryID,
						Probability: countryInfo.Probability,
					}

					person.Country = append(person.Country, country)
				}
			}
		}
	}
	return nil
}

func (s *Service) HandleAddPerson(w http.ResponseWriter, r *http.Request) {
	var newPerson = new(models.Person)

	err := json.NewDecoder(r.Body).Decode(newPerson)
	if err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = EnrichPersonInfo(newPerson); err != nil {
		handleError(w, "Failed to add extra data to the person", http.StatusInternalServerError)
		return
	}
	err = database.AddPerson(newPerson, s.DB)
	if err != nil {
		handleError(w, "Failed to add person to the database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPerson)
}

func (s *Service) HandleDeletePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	personID := vars["id"]

	id, err := strconv.Atoi(personID)
	if err != nil {
		handleError(w, "Invalid person ID", http.StatusBadRequest)
		return
	}
	if err = database.DeletePersonByID(id, s.DB); err != nil {
		handleError(w, "Person not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func (s *Service) HandleChangePerson(w http.ResponseWriter, r *http.Request) error {
	userID := mux.Vars(r)["tweet_id"]
	if userID == "" {
		return errors.New("Missing person ID")
	}

	person, err := s.GetPersonByID(userID)
	if err != nil {
		return err
	}

	var updatePerson = new(models.ChangePerson)
	err = json.NewDecoder(r.Body).Decode(&updatePerson)
	if err != nil {
		return err
	}

	s.UpdateFromChangePerson(person, updatePerson)

	if err = database.UpdatePerson(person, s.DB); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
func (s *Service) GetPersonByID(id string) (*models.Person, error) {
	query := "SELECT id, name, age FROM persons WHERE id = ?"
	person := &models.Person{}
	err := s.DB.QueryRow(query, id).Scan(&person.ID, &person.Name, &person.Age)
	if err != nil {
		return nil, err
	}
	return person, nil
	return nil, errors.New("GetPersonByID method not implemented")
}
func (s *Service) UpdateFromChangePerson(p *models.Person, updatePerson *models.ChangePerson) {
	if updatePerson.Name != "" {
		p.Name = updatePerson.Name
	}
	if updatePerson.Surname != "" {
		p.Surname = updatePerson.Surname
	}
	if updatePerson.Gender != "" {
		p.Gender = updatePerson.Gender
	}
	if updatePerson.Age != 0 {
		p.Age = updatePerson.Age
	}

}

func SetupRoutes(api *Service) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/people/{id}", api.HandleDeletePerson).Methods(http.MethodDelete)
	router.HandleFunc("/people", api.HandleAddPerson).Methods(http.MethodPost)
	return router
}
