package api

import (
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)
type testData struct{
	method string
	url string
	status int
	handler http.HandlerFunc
	path string
	body io.Reader
	response string
}

func init() {
	CreateInitialWorkerProfile()
}

func TestShowAllWorkers(t *testing.T) {

	test := []testData{
		{
			"GET",
			"/appscode/workers",
			200,
			ShowAllWorkers,
			"/appscode/workers",
			nil,
			`{"fahim":{"username":"fahim","firstname":"Fahim","lastname":"Abrar","address":{"city":"Chittagong","division":"Chittagong"},"position":"Software Engineer","salary":55},"jenny":{"username":"jenny","firstname":"Jannatul","lastname":"Ferdows","address":{"city":"Chittagong","division":"Chittagong"},"position":"Software Engineer","salary":55},"masud":{"username":"masud","firstname":"Masudur","lastname":"Rahman","address":{"city":"Madaripur","division":"Dhaka"},"position":"Software Engineer","salary":55},"tahsin":{"username":"tahsin","firstname":"Tahsin","lastname":"Rahman","address":{"city":"Chittagong","division":"Chittagong"},"position":"Software Engineer","salary":55}}`,
		},
	}


	for _, data := range test{
		runTest(data, t)
	}

}

func TestShowSingleWorker(t *testing.T) {
	test := []testData{
		{
			"GET",
			"/appscode/workers/masud",
			200,
			ShowSingleWorker,
			"/appscode/workers/{username}",
			nil,
			`{"username":"masud","firstname":"Masudur","lastname":"Rahman","address":{"city":"Madaripur","division":"Dhaka"},"position":"Software Engineer","salary":55}`,
		},
		{
			"GET",
			"/appscode/workers/fahim",
			200,
			ShowSingleWorker,
			"/appscode/workers/{username}",
			nil,
			`{"username":"fahim","firstname":"Fahim","lastname":"Abrar","address":{"city":"Chittagong","division":"Chittagong"},"position":"Software Engineer","salary":55}`,
		},
		{
			"GET",
			"/appscode/workers/tahsin",
			200,
			ShowSingleWorker,
			"/appscode/workers/{username}",
			nil,
			`{"username":"tahsin","firstname":"Tahsin","lastname":"Rahman","address":{"city":"Chittagong","division":"Chittagong"},"position":"Software Engineer","salary":55}`,
		},
		{
			"GET",
			"/appscode/workers/jenny",
			200,
			ShowSingleWorker,
			"/appscode/workers/{username}",
			nil,
			`{"username":"jenny","firstname":"Jannatul","lastname":"Ferdows","address":{"city":"Chittagong","division":"Chittagong"},"position":"Software Engineer","salary":55}`,
		},
		{
			"GET",
			"/appscode/workers/abcd",
			404,
			ShowSingleWorker,
			"/appscode/workers/{username}",
			nil,
			`404 - Content Not Found`,
		},


/*
		{method: "GET", url: "/appscode/workers/masud", status: 200, handler: ShowSinigleWorker, path: "/appscode/workers/{username}", body: nil},
		{method: "GET", url: "/appscode/workers/tahsin", status: 200, handler: ShowSinigleWorker, path: "/appscode/workers/{username}", body: nil},
		{method: "GET", url: "/appscode/workers/fahim", status: 200, handler: ShowSinigleWorker, path: "/appscode/workers/{username}", body: nil},*/
	}

	for _, data := range test{
		runTest(data, t)

	}

}

func TestAddNewWorker(t *testing.T) {
	test := []testData{
		{
			"POST",
			"/appscode/workers",
			409,
			AddNewWorker,
			"/appscode/workers",
			strings.NewReader(`{"username":"masud","firstname":"Masudur","lastname":"Rahman","address":{"city":"Madaripur","division":"Dhaka"},"position":"Software Engineer","salary":55}`),
			`409 - username already exists`,
		},
		{
			"POST",
			"/appscode/workers",
			201,
			AddNewWorker,
			"/appscode/workers",
			strings.NewReader(`{"username":"masudur","firstname":"Masudur","lastname":"Rahman","address":{"city":"Madaripur","division":"Dhaka"},"position":"Software Engineer","salary":55}`),
			`201 - Created successfully`,
		},
	}

	for _, data := range test {
		runTest(data, t)
	}

}

func TestUpdateWorkerProfile(t *testing.T) {
	test := []testData{
		{
			"POST",
			"/appscode/workers/masud",
			405,
			UpdateWorkerProfile,
			"/appscode/workers/{username}",
			strings.NewReader(`{"username":"masudd","firstname":"Masudur","lastname":"Rahman","address":{"city":"Madaripur","division":"Dhaka"},"position":"Software Engineer","salary":55}`),
			`405 - Username can't be changed`,
		},
		{
			"POST",
			"/appscode/workers/masudd",
			404,
			UpdateWorkerProfile,
			"/appscode/workers/{username}",
			strings.NewReader(`{"username":"masudd","firstname":"Masudur","lastname":"Rahman","address":{"city":"Madaripur","division":"Dhaka"},"position":"Software Engineer","salary":55}`),
			`404 - Username Doesn't Exist`,
		},
		{
			"POST",
			"/appscode/workers/masud",
			201,
			UpdateWorkerProfile,
			"/appscode/workers/{username}",
			strings.NewReader(`{"username":"masud","firstname":"Masud","lastname":"Rahman","address":{"city":"M","division":"D"},"position":"Software Engineer","salary":55}`),
			`201 - Updated successfully`,
		},
	}

	for _, data := range test{
		runTest(data, t)
	}
}

func TestDeleteWorker(t *testing.T) {
	test := []testData{
		{
			"DELETE",
			"/appscode/workers/masud",
			200,
			DeleteWorker,
			"/appscode/workers/{username}",
			nil,
			`200 - Deleted Successfully`,
		},
		{
			"DELETE",
			"/appscode/workers/fahim",
			200,
			DeleteWorker,
			"/appscode/workers/{username}",
			nil,
			`200 - Deleted Successfully`,
		},
		{
			"DELETE",
			"/appscode/workers/hello",
			404,
			DeleteWorker,
			"/appscode/workers/{username}",
			nil,
			`404 - Content Not Found`,
		},
	}
	for _, data := range test{
		runTest(data, t)
	}
}

func runTest(test testData, t *testing.T){
	req, err := http.NewRequest(test.method, test.url, test.body)
	if err != nil{
		t.Fatal(err)
	}
	responseRecorder := httptest.NewRecorder()
	//handler := http.HandlerFunc(test.handler)
	//handler.ServeHTTP(responseRecorder, req)

	router := mux.NewRouter()
	router.HandleFunc(test.path, test.handler)
	router.ServeHTTP(responseRecorder, req)

	/*if test.body != nil {
		fmt.Println("1: ", test.path, test.url, test.body)
		fmt.Println("2: ", responseRecorder.Code, test.status)
	}*/

	if status := responseRecorder.Code; status != test.status {
		t.Errorf("handler returned wrong status code: got %v expected %v", status, test.status)
	}

}