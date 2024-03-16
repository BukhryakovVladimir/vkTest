package routes

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/BukhryakovVladimir/vkTest/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"time"
)

func GetMoviesOrdered(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	cookie, err := r.Cookie(jwtName)

	if err != nil {
		http.Error(w, "Unauthenticated", http.StatusUnauthorized)
		return
	}

	token, err := jwtCheck(cookie)

	if err != nil {
		http.Error(w, "Unauthenticated", http.StatusUnauthorized)
		return
	}

	claims := token.Claims.(*jwt.RegisteredClaims)

	userExists, err := checkUserExists(claims.Issuer)
	if err != nil {
		http.Error(w, "Error while checking user authorization", http.StatusInternalServerError)
		return
	}

	if !userExists {
		log.Println("User with id ", claims.Issuer, "does not exist: ", err)
		http.Error(w, "You are not logged in", http.StatusUnauthorized)
		return
	}

	order := r.URL.Query().Get("order")

	if order != "asc" && order != "desc" {
		order = "desc"
	}

	by := r.URL.Query().Get("by")

	if by != "name" && by != "rating" && by != "date" {
		by = "rating"
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(queryTimeLimit)*time.Second)
	defer cancel()

	var getMoviesQuery string

	if order == "desc" {
		switch by {
		case "name":
			getMoviesQuery = `
		SELECT m.name, m.description, m.date, m.rating, a.firstName, a.lastName, a.sex, a.birthDate
		FROM movie m
		JOIN actormovie ma ON m.id = ma.movie_id
		JOIN actor a ON a.id = ma.actor_id
		ORDER BY m.name DESC
		`
		case "rating":
			getMoviesQuery = `
		SELECT m.name, m.description, m.date, m.rating, a.firstName, a.lastName, a.sex, a.birthDate
		FROM movie m
		JOIN actormovie ma ON m.id = ma.movie_id
		JOIN actor a ON a.id = ma.actor_id
		ORDER BY m.rating DESC
		`
		case "date":
			getMoviesQuery = `
		SELECT m.name, m.description, m.date, m.rating, a.firstName, a.lastName, a.sex, a.birthDate
		FROM movie m
		JOIN actormovie ma ON m.id = ma.movie_id
		JOIN actor a ON a.id = ma.actor_id
		ORDER BY m.date DESC
		`
		default:
		}
	} else {
		switch by {
		case "name":
			getMoviesQuery = `
		SELECT m.name, m.description, m.date, m.rating, a.firstName, a.lastName, a.sex, a.birthDate
		FROM movie m
		JOIN actormovie ma ON m.id = ma.movie_id
		JOIN actor a ON a.id = ma.actor_id
		ORDER BY m.name
		`
		case "rating":
			getMoviesQuery = `
		SELECT m.name, m.description, m.date, m.rating, a.firstName, a.lastName, a.sex, a.birthDate
		FROM movie m
		JOIN actormovie ma ON m.id = ma.movie_id
		JOIN actor a ON a.id = ma.actor_id
		ORDER BY m.rating
		`
		case "date":
			getMoviesQuery = `
		SELECT m.name, m.description, m.date, m.rating, a.firstName, a.lastName, a.sex, a.birthDate
		FROM movie m
		JOIN actormovie ma ON m.id = ma.movie_id
		JOIN actor a ON a.id = ma.actor_id
		ORDER BY m.date
		`
		default:
		}
	}

	rows, err := db.QueryContext(ctx, getMoviesQuery)
	defer rows.Close()
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Println("GetMoviesOrdered QueryContext deadline exceeded: ", err)
			http.Error(w, "Database query time limit exceeded", http.StatusGatewayTimeout)
			return
		} else {
			log.Println("Database error: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	var movies []model.Movie
	var currentMovie *model.Movie

	for rows.Next() {
		var movie model.Movie
		var actor model.Actor

		if err := rows.Scan(&movie.Name, &movie.Description, &movie.Date, &movie.Rating,
			&actor.FirstName, &actor.LastName, &actor.Sex, &actor.BirthDate); err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Check if we're still processing the same movie
		if currentMovie != nil && movie.Name == currentMovie.Name && movie.Date == currentMovie.Date {
			// Add actor to the current movie's actor list
			currentMovie.Actors = append(currentMovie.Actors, actor)
		} else {
			// We've encountered a new movie, so add the previous one to the movies slice
			if currentMovie != nil {
				movies = append(movies, *currentMovie)
			}
			// Start aggregating actors for the new movie
			movie.Actors = []model.Actor{actor}
			currentMovie = &movie
		}
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if currentMovie != nil {
		movies = append(movies, *currentMovie)
	}

	resp, err := json.Marshal(movies)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Printf("Write failed: %v\n", err)
	}
}
