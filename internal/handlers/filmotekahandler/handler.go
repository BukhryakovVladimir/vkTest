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

	mux.HandleFunc("POST /api/add-movie", routes.AddMovie)
	mux.HandleFunc("PUT /api/update-movie", routes.UpdateMovie)
	mux.HandleFunc("DELETE /api/delete-movie", routes.DeleteMovie)
	mux.HandleFunc("POST /api/get-movies-with-id", routes.GetMoviesWithID)
	mux.HandleFunc("POST /api/add-actor-to-movie", routes.AddActorToMovie)
	mux.HandleFunc("DELETE /api/delete-actor-from-movie", routes.DeleteActorFromMovie)

	mux.HandleFunc("GET /api/movies", routes.GetMoviesOrdered)
	mux.HandleFunc("POST /api/search-movie", routes.SearchMovie)
}
