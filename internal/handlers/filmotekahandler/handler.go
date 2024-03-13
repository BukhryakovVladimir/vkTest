package filmotekahandler

import (
	"github.com/BukhryakovVladimir/vkTest/internal/routes"
	"net/http"
)

func SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/signup", routes.SignupPerson)
	mux.HandleFunc("POST /api/login", routes.LoginPerson)

	mux.HandleFunc("POST /api/add-actor", routes.AddActor)
	mux.HandleFunc("PUT /api/update-actor", routes.UpdateActor)
	mux.HandleFunc("DELETE /api/delete-actor", routes.DeleteActor)
	mux.HandleFunc("POST /api/get-actors-with-id", routes.GetActorsWithID)
}
