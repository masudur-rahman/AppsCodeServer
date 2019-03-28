package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
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

// List of workers and authenticated users
var Workers = make(map[string]Worker)
var authUser = make(map[string]string)

var srvr http.Server
var byPass bool = true
var stopTime int16

// Handler Functions....
func Welcome(w http.ResponseWriter, r *http.Request) {

	if err := json.NewEncoder(w).Encode("Congratulations...! Your API Server is up and running... :) "); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func WelcomeToAppsCode(w http.ResponseWriter, r *http.Request) {

	if err := json.NewEncoder(w).Encode("Welcome to AppsCode Ltd.. Available Links are : `/appscode/workers`, `/appscode/workers/{username}`"); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}


func ShowAllWorkers(w http.ResponseWriter, r *http.Request) {

	if info, valid := basicAuth(r); !valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(info))
		return
	}
	if err := json.NewEncoder(w).Encode(Workers); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}

func ShowSingleWorker(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("okay")
	params := mux.Vars(r)
	//fmt.Println("Username from parameter:", params["username"])

	if info, valid := basicAuth(r); !valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(info))
		return
	}

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

	if info, valid := basicAuth(r); !valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(info))
		return
	}

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
	//json.NewEncoder(w).Encode(Workers)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("201 - Created successfully"))
}

func UpdateWorkerProfile(w http.ResponseWriter, r *http.Request) {

	if info, valid := basicAuth(r); !valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(info))
		return
	}

	params := mux.Vars(r)

	var worker Worker

	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil {
		panic(err)
	}


	if _, exist := Workers[params["username"]]; !exist {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Username Doesn't Exist"))
		return
	}

	if params["username"] != worker.Username {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Username can't be changed"))
		return
	}

	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()
	Workers[worker.Username] = worker
	// json.NewEncoder(w).Encode(Workers)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("201 - Updated successfully"))
}

func DeleteWorker(w http.ResponseWriter, r *http.Request) {

	if info, valid := basicAuth(r); !valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(info))
		return
	}

	params := mux.Vars(r)

	if _, exist := Workers[params["username"]]; !exist {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Content Not Found"))
		return
	}

	delete(Workers, params["username"])
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Deleted Successfully"))
}

// Creating initial worker profiles
func CreateInitialWorkerProfile() {

	/*Workers = map[string]Worker{
		"masud": {Person{
			Username: "masud", FirstName:"Masudur", LastName:"Rahman", Address{City: "Madaripur", Division: "Dhaka"} },Position: "Software Engineer", Salary: 55}
	}*/

	worker := Worker{
		Person: Person{Username: "masud",
			FirstName: "Masudur",
			LastName:  "Rahman",
			Address:   Address{City: "Madaripur", Division: "Dhaka"}},
		Position: "Software Engineer",
		Salary:   55,
	}
	Workers["masud"] = worker

	worker = Worker{
		Person: Person{Username: "fahim",
			FirstName: "Fahim",
			LastName:  "Abrar",
			Address:   Address{City: "Chittagong", Division: "Chittagong"}},
		Position: "Software Engineer",
		Salary:   55,
	}
	Workers["fahim"] = worker

	worker = Worker{
		Person: Person{Username: "tahsin",
			FirstName: "Tahsin",
			LastName:  "Rahman",
			Address:   Address{City: "Chittagong", Division: "Chittagong"}},
		Position: "Software Engineer",
		Salary:   55,
	}
	Workers["tahsin"] = worker

	worker = Worker{
		Person: Person{Username: "jenny",
			FirstName: "Jannatul",
			LastName:  "Ferdows",
			Address:   Address{City: "Chittagong", Division: "Chittagong"}},
		Position: "Software Engineer",
		Salary:   55,
	}
	Workers["jenny"] = worker

	authUser["masud"] = "pass"
	authUser["admin"] = "admin"

}

func basicAuth(r *http.Request) (string, bool) {
	if byPass{
		return "", true
	}
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "Error: Authorization Needed...!", false
	}

	authInfo := strings.SplitN(authHeader, " ", 2)

	userInfo, err := base64.StdEncoding.DecodeString(authInfo[1])

	if err != nil {
		return "Error: Error while decoding...!", false
	}
	userPass := strings.SplitN(string(userInfo), ":", 2)

	if len(userPass) != 2 {
		return "Error: Authorization failed...!", false
	}

	if pass, exist := authUser[userPass[0]]; exist {
		if pass != userPass[1] {
			return "Error: Unauthorized User", false
		} else {
			return "Success: Authorization Successful...!!", true
		}
	} else {
		return "Error: Unauthorized User...!", false
	}
}

func AssignValues(port string, bypass bool, stop int16) {
	srvr.Addr = ":" + port
	byPass = bypass
	stopTime = stop
}

func StartTheApp() {

	log.Println("Starting the server...!")

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second * 15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	CreateInitialWorkerProfile()

	router := mux.NewRouter()

	router.HandleFunc("/", Welcome).Methods("GET")
	router.HandleFunc("/appscode", WelcomeToAppsCode).Methods("GET")
	router.HandleFunc("/appscode/workers", ShowAllWorkers).Methods("GET")
	router.HandleFunc("/appscode/workers/{username}", ShowSingleWorker).Methods("GET")
	router.HandleFunc("/appscode/workers", AddNewWorker).Methods("POST")
	router.HandleFunc("/appscode/workers/{username}", UpdateWorkerProfile).Methods("PUT")
	router.HandleFunc("/appscode/workers/{username}", DeleteWorker).Methods("DELETE")

	srvr.WriteTimeout = time.Second * 15
	srvr.ReadTimeout = time.Second * 15
	srvr.IdleTimeout = time.Second * 60

	srvr.Handler = router
	log.Println("Just before starting the server")

	go func() {
		if err := srvr.ListenAndServe(); err != nil {
			log.Println("Nothing")
		}
	}()


	// Channel to interrupt the server from keyboard
	channel := make(chan os.Signal, 1)

	signal.Notify(channel, os.Interrupt)
	<- channel

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	//  Shutting down the server
	log.Println("Shutting down the server...!")

	time.Sleep(time.Second * time.Duration(stopTime))

	srvr.Shutdown(ctx)

	log.Println("The server has been shut down...!")

	os.Exit(0)
}

/*
func main2() {
	router := mux.NewRouter()

	CreateInitialWorkerProfile()

	router.HandleFunc("/appscode/workers", ShowAllWorkers).Methods("GET")
	router.HandleFunc("/appscode/workers/{username}", ShowSinigleWorker).Methods("GET")
	router.HandleFunc("/appscode/workers", AddNewWorker).Methods("POST")
	router.HandleFunc("/appscode/workers/{username}", UpdateWorkerProfile).Methods("PUT")
	router.HandleFunc("/appscode/workers/{username}", DeleteWorker).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}*/
