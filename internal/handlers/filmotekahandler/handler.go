package filmotekahandler

import (
	"github.com/BukhryakovVladimir/vkTest/internal/routes"
	"net/http"
)

func SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/signup", routes.SignupPerson)
	mux.HandleFunc("POST /api/login", routes.LoginPerson)
}
