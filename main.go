package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// Structures of users
type Person struct {
	Username  string  `json:"username"`
	FirstName string  `json:"firstname"`
	LastName  string  `json:"lastname"`
	Address   Address `json:"address"`
}

type Address struct {
	City     string `json:"city"`
	Division string `json:"division"`
}


type Worker struct {
	Person
	Position string `json:"position"`
	Salary   int    `json:"salary"`
}


var Workers = make(map[string]Worker)


// Handler Functions....
func ShowAllWorkers(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(Workers); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}


func ShowSinigleWorker(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if info, exist := Workers[params["username"]]; exist {
		json.NewEncoder(w).Encode(info)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Content Not Found"))
		return
	}

	w.WriteHeader(http.StatusOK)
}
func AddNewWorker(w http.ResponseWriter, r *http.Request) {
	var worker Worker
	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil {
		panic(err)
	}

	if _, exist := Workers[worker.Username]; exist {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("409 - username already exists"))
		return
	}

	Workers[worker.Username] = worker
	json.NewEncoder(w).Encode(Workers)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("201 - Created successfully"))
}

func UpdateWorkerProfile(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)
	var worker Worker

	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil {
		panic(err)
	}

	if _, exist := Workers[worker.Username]; !exist {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Content Not Found"))
		return
	}

	Workers[worker.Username] = worker
	json.NewEncoder(w).Encode(Workers)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("201 - Updated successfully"))
}

func DeleteWorker(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if _, exist := Workers[params["username"]]; !exist {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Content Not Found"))
		return
	}

	delete(Workers, params["username"])
	w.WriteHeader(http.StatusOK)
}

// Creating initial worker profiles
func CreateInitialWorkerProfile() {
	/*Workers = map[string]Worker{
		"masud": {Person{
			Username: "masud", FirstName:"Masudur", LastName:"Rahman", Address{City: "Madaripur", Division: "Dhaka"} },Position: "Software Engineer", Salary: 55}
	}*/

	worker := Worker{
		Person:Person{Username:"masud",
			FirstName: "Masudur",
			LastName: "Rahman",
			Address: Address{City: "Madaripur", Division: "Dhaka"}},
			Position:"Software Engineer",
			Salary: 55,
	}
	Workers["masud"] = worker

	worker = Worker{
		Person:Person{Username:"fahim",
			FirstName: "Fahim",
			LastName: "Abrar",
			Address: Address{City: "Chittagong", Division: "Chittagong"}},
			Position:"Software Engineer",
			Salary: 55,
	}
	Workers["fahim"] = worker

	worker = Worker{
		Person:Person{Username:"tahsin",
			FirstName: "Tahsin",
			LastName: "Rahman",
			Address: Address{City: "Chittagong", Division: "Chittagong"}},
			Position:"Software Engineer",
			Salary: 55,
	}
	Workers["tahsin"] = worker

	worker = Worker{
		Person:Person{Username:"jenny",
			FirstName: "Jannatul",
			LastName: "Ferdows",
			Address: Address{City: "Chittagong", Division: "Chittagong"}},
			Position:"Software Engineer",
			Salary: 55,
	}
	Workers["jenny"] = worker

}



func main() {
	router := mux.NewRouter()

	CreateInitialWorkerProfile()

	router.HandleFunc("/appscode/workers", ShowAllWorkers).Methods("GET")
	router.HandleFunc("/appscode/workers/{username}", ShowSinigleWorker).Methods("GET")
	router.HandleFunc("/appscode/workers", AddNewWorker).Methods("POST")
	router.HandleFunc("/appscode/workers/{username}", UpdateWorkerProfile).Methods("PUT")
	router.HandleFunc("/appscode/workers/{username}", DeleteWorker).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
