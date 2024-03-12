package filmotekahandler

import (
	"github.com/BukhryakovVladimir/vkTest/internal/routes"
	"net/http"
)

type apiHandler struct{}

func (apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func SetupRoutes(mux *http.ServeMux) {

	mux.Handle("/api/", apiHandler{})
	mux.HandleFunc("POST /signup", routes.SignupPerson)
}
