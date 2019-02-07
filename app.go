package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"

	"github.com/mleyb/go-api-mongo/config"
	"github.com/mleyb/go-api-mongo/model"
	"github.com/mleyb/go-api-mongo/dao"
)

var cfg = config.Config{}
var db = dao.MoviesDAO{}

// GET - get all movies
func GetMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := db.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, movies)
}

// GET - get an existing movie by id
func GetMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movie, err := db.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Movie Id")
		return
	}
	respondWithJson(w, http.StatusOK, movie)
}

// POST - create a new movie
func CreateMovie(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie model.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	movie.ID = bson.NewObjectId()
	if err := db.Insert(movie); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, movie)
}

// PUT - update an existing movie
func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie model.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := db.Update(movie); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// DELETE - delete an existing movie
func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie model.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := db.Delete(movie); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/movies", GetMovies).Methods("GET")
	router.HandleFunc("/movies/{id}", GetMovie).Methods("GET")
	router.HandleFunc("/movies", CreateMovie).Methods("POST")
	router.HandleFunc("/movies", UpdateMovie).Methods("PUT")
	router.HandleFunc("/movies", DeleteMovie).Methods("DELETE")

	if err := http.ListenAndServe(":3000", router); err != nil {
		log.Fatal(err)
	}
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}