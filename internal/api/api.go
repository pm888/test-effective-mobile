// internal/api/api.go

package api

import (
	"net/http"
	"strconv"

	"test_effective_mobile/test-effective-mobile/internal/database"

	"github.com/gorilla/mux"
)

func (s *Service) HandleDeletePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	personID := vars["id"]

	id, err := strconv.Atoi(personID)
	if err != nil {
		http.Error(w, "Invalid person ID", http.StatusBadRequest)
		return
	}
	if err = database.DeletePersonByID(id, s.DB); err != nil {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func SetupRoutes(api *Service) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/people/{id}", api.HandleDeletePerson).Methods(http.MethodDelete)
	return router
}
