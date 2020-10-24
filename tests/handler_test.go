package tests

import (
	"bytes"
	"encoding/json"
	"fyndtest/internal/handlers"
	"fyndtest/internal/object"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

//Test function for SignUp Handler
func TestSignUpHandler(t *testing.T) {
	var jsonStr = []byte(`{
		"username":"rishabhadmin4",
		"password":"12345",
		"admin_secret":"admin",
		"email":"rishab25@gmail.com"
	}`)
	req, err := http.NewRequest("POST", "/api/v1/signup", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.Signup)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"message":"User Created Successfully"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

//Test function for Login Handler
func TestLoginHandler(t *testing.T) {
	var jsonStr = []byte(`{
		"username":"rishabhadmin4",
		"password":"12345"
	}`)
	req, err := http.NewRequest("GET", "/api/v1/login", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.Login)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

////Test function for Retrieve All Movies Handler
func TestFetchAllMovies(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/movies", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.Movies)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

////Test function for Add Movie Handler
func TestAddMovies(t *testing.T) {
	var jsonStr = []byte(`{
			"username":"rishabhadmin",
			"password":"12345"
		}`)
	req, err := http.NewRequest("GET", "/api/v1/login", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.Login)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	body, readErr := ioutil.ReadAll(rr.Body)
	if readErr != nil {
		t.Fatal(readErr)
	}
	var loginresp object.LoginResponse
	err = json.Unmarshal(body, &loginresp)
	if err != nil {
		t.Fatal(err)
	}
	jsonStr = []byte(`{
		"99popularity": 79.0,
		"director": "William Cottrell",
		"genre": [
		  "Animation",
		  " Family",
		  " Fantasy",
		  " Musical",
		  " Romance"
		],
		"imdb_score": 7.9,
		"name": "Snow White and the Seven Dwarfs"
	  }`)
	req, err = http.NewRequest("POST", "/api/v1/movies", bytes.NewBuffer(jsonStr))
	var bearer = "Bearer " + loginresp.AccessToken
	req.Header.Add("Authorization", bearer)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(handlers.Movies)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

//Test function to Test Handler to Retrieve All Movies Based on Query
func TestSearchMovies(t *testing.T) {
	var jsonstr = []byte(`{
		"term":{
			"imdb_score":9.1
		},
		"range":[
			{
				"key":"imdb_score",
				"operator":"eq",
				"value":9.1
			},
			{
				"key":"99popularity",
				"operator":"gt",
				"value":50
			}
		]
	}`)
	req, err := http.NewRequest("GET", "/api/v1/movies/search", bytes.NewBuffer(jsonstr))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.Movies)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

//Test function to Test Handler to Update Movies By Id
func TestUpdateMovies(t *testing.T) {
	var jsonStr = []byte(`{
			"username":"rishabhadmin",
			"password":"12345"
		}`)
	req, err := http.NewRequest("GET", "/api/v1/login", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.Login)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	body, readErr := ioutil.ReadAll(rr.Body)
	if readErr != nil {
		t.Fatal(readErr)
	}
	var loginresp object.LoginResponse
	err = json.Unmarshal(body, &loginresp)
	if err != nil {
		t.Fatal(err)
	}
	jsonStr = []byte(`{
		"key":"imdb_score",
		"value":"9.1"
	}`)
	req, err = http.NewRequest("PUT", "/api/v1/movies/11", bytes.NewBuffer(jsonStr))
	req.RequestURI = "/api/v1/movies/12"
	var bearer = "Bearer " + loginresp.AccessToken
	req.Header.Add("Authorization", bearer)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(handlers.Movies)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

//Test function to Test Handler to Delete Movies By Id
func TestDeleteMovies(t *testing.T) {
	var jsonStr = []byte(`{
			"username":"rishabhadmin",
			"password":"12345"
		}`)
	req, err := http.NewRequest("GET", "/api/v1/login", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.Login)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	body, readErr := ioutil.ReadAll(rr.Body)
	if readErr != nil {
		t.Fatal(readErr)
	}
	var loginresp object.LoginResponse
	err = json.Unmarshal(body, &loginresp)
	if err != nil {
		t.Fatal(err)
	}
	req, err = http.NewRequest("DELETE", "/api/v1/movies/12", nil)
	req.RequestURI = "/api/v1/movies/12"
	var bearer = "Bearer " + loginresp.AccessToken
	req.Header.Add("Authorization", bearer)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(handlers.Movies)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
