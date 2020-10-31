package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"demorestapi/internal/database"
	"demorestapi/internal/object"
	"demorestapi/internal/util"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

//Declaration of token Active Time to 10 mins
const (
	dur = "10m"
)

//Declaration of cache maps
var (
	mutx            *sync.Mutex
	accesstokenmap  map[string]string
	userloggedinmap map[string]string
	durationmap     map[string]string
	adminaccess     map[string]bool
)

//InitializeMaps Initializes Cache Map
func InitializeMaps() {
	mutx = &sync.Mutex{}
	if accesstokenmap == nil {
		accesstokenmap = make(map[string]string)
	}
	if userloggedinmap == nil {
		userloggedinmap = make(map[string]string)
	}
	if durationmap == nil {
		durationmap = make(map[string]string)
	}
	if adminaccess == nil {
		adminaccess = make(map[string]bool)
	}
}

//Signup SignUp Handler User, Admin Get Registered Through This Handler
func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusOK)
		var reqobj object.SignupRequest
		var resp *object.Response
		if err := json.NewDecoder(r.Body).Decode(&reqobj); err != nil {
			resp = &object.Response{
				Message: err.Error(),
			}
		}
		resp = database.Signup(&reqobj)
		byteresp, _ := json.Marshal(resp)
		w.Write(byteresp)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := object.Response{
			Message: "Method Not Supported",
		}
		byteresp, _ := json.Marshal(response)
		w.Write(byteresp)
	}
}

//Login Handler: User Gets Authenticated Through This Handler,
//On Successfull Auth start time is stored in cache map
//Will have to Use JWT, but for demo just created, a token generator
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		var reqobj object.LoginRequest
		var resp *object.LoginResponse
		if err := json.NewDecoder(r.Body).Decode(&reqobj); err != nil {
			resp = &object.LoginResponse{
				Message: err.Error(),
			}
		}
		if accesstokenmap == nil {
			accesstokenmap = make(map[string]string)
		}
		if userloggedinmap == nil {
			userloggedinmap = make(map[string]string)
		}
		if durationmap == nil {
			durationmap = make(map[string]string)
		}
		if adminaccess == nil {
			adminaccess = make(map[string]bool)
		}
		if mutx == nil {
			mutx = &sync.Mutex{}
		}
		admin := false
		if accesstoken, ok := userloggedinmap[reqobj.Username]; ok {
			if settime, k := durationmap[accesstoken]; k {
				startTime, _ := time.Parse("2006-01-02 15:04:05.00", settime)
				currentTimeVal := time.Now().Format("2006-01-02 15:04:05.00")
				currentTime, _ := time.Parse("2006-01-02 15:04:05.00", currentTimeVal)
				duration := currentTime.Sub(startTime)
				definedDuration, _ := time.ParseDuration(dur)
				if duration <= definedDuration {
					resp = &object.LoginResponse{
						Message:     "User Logged In",
						AccessToken: fmt.Sprintf("%s", accesstoken),
					}
				} else {
					resp, admin = database.Login(&reqobj)
				}
			} else {
				resp, admin = database.Login(&reqobj)
			}
		} else {
			resp, admin = database.Login(&reqobj)
		}
		byteresp, _ := json.Marshal(resp)
		currenttime := time.Now().Format("2006-01-02 15:04:05.00")
		if admin {
			mutx.Lock()
			adminaccess[resp.AccessToken] = true
			mutx.Unlock()
		}
		mutx.Lock()
		userloggedinmap[reqobj.Username] = resp.AccessToken
		durationmap[resp.AccessToken] = currenttime
		accesstokenmap[resp.AccessToken] = reqobj.Username
		mutx.Unlock()
		if resp.AccessToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
		}
		w.Write(byteresp)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := object.LoginResponse{
			Message: "Method Not Supported",
		}
		byteresp, _ := json.Marshal(response)
		w.Write(byteresp)
	}
}

//Search Handler, Searches Movies,Current All calls go to database,as users grow have to use some caching
func Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		var resp *object.GetAllMovieResponse
		var reqobj *object.SearchMovie
		w.WriteHeader(http.StatusOK)
		if err := json.NewDecoder(r.Body).Decode(&reqobj); err != nil {
			resp = &object.GetAllMovieResponse{
				Message: err.Error(),
			}
		}
		resp = database.SearchMovies(reqobj)
		byteresp, _ := json.Marshal(resp)
		w.Write(byteresp)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := object.Response{
			Message: "Method Not Supported",
		}
		byteresp, _ := json.Marshal(response)
		w.Write(byteresp)
	}
}

//Movies Handler does CRUD operation on movies, Only allows Admin to Perform CUD
func Movies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		var resp *object.GetAllMovieResponse
		var id int
		var err error
		uri := strings.Split(util.After(r.RequestURI, "/api/v1/movies/"), "/")
		if uri[0] != "" {
			id, err = CheckURIAndRetrieveID(uri)
			if err != nil {
				resp = &object.GetAllMovieResponse{
					Message: err.Error(),
				}
				w.WriteHeader(http.StatusNotFound)
				byteresp, _ := json.Marshal(resp)
				w.Write(byteresp)
				return
			}
		}
		if id != 0 {
			resp = database.FetchMoviesByID(id)
		} else {
			resp = database.FetchAllMovies()
		}
		w.WriteHeader(http.StatusOK)
		byteresp, _ := json.Marshal(resp)
		w.Write([]byte(byteresp))
	case "POST":
		var reqobj object.AddMovieRequest
		var resp *object.Response
		if err := json.NewDecoder(r.Body).Decode(&reqobj); err != nil {
			resp = &object.Response{
				Message: err.Error(),
			}
		}
		authorizationHeader := r.Header.Get("authorization")
		bearerToken := strings.Split(authorizationHeader, " ")
		accessToken := bearerToken[1]
		val := adminaccess[accessToken]
		if !val {
			w.WriteHeader(http.StatusForbidden)
			resp = &object.Response{
				Message: "POST Permission Denied",
			}
		}
		if _, ok := adminaccess[accessToken]; ok {
			resp = database.AddMovies(&reqobj, accesstokenmap[accessToken])
		} else {
			w.WriteHeader(http.StatusForbidden)
			resp = &object.Response{
				Message: "Access Denied",
			}
		}
		byteresp, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(byteresp))
	case "PUT":
		var reqobj object.UpdateMovie
		var resp *object.Response
		if err := json.NewDecoder(r.Body).Decode(&reqobj); err != nil {
			resp = &object.Response{
				Message: err.Error(),
			}
		}
		uri := strings.Split(util.After(r.RequestURI, "/api/v1/movies/"), "/")
		id, err := CheckURIAndRetrieveID(uri)
		if err != nil {
			resp = &object.Response{
				Message: err.Error(),
			}
			w.WriteHeader(http.StatusNotFound)
		} else {
			authorizationHeader := r.Header.Get("authorization")
			bearerToken := strings.Split(authorizationHeader, " ")
			accessToken := bearerToken[1]
			val := adminaccess[accessToken]
			if !val {
				w.WriteHeader(http.StatusForbidden)
				resp = &object.Response{
					Message: "POST Permission Denied",
				}
			}
			if _, ok := adminaccess[accessToken]; ok {
				resp = database.UpdateMovieByID(&reqobj, id, accesstokenmap[accessToken])
			} else {
				w.WriteHeader(http.StatusForbidden)
				resp = &object.Response{
					Message: "Access Denied",
				}
			}
		}
		byteresp, _ := json.Marshal(resp)
		w.Write(byteresp)
	case "DELETE":
		var resp *object.Response
		uri := strings.Split(util.After(r.RequestURI, "/api/v1/movies/"), "/")
		id, err := CheckURIAndRetrieveID(uri)
		if err != nil {
			resp = &object.Response{
				Message: err.Error(),
			}
			w.WriteHeader(http.StatusNotFound)
		} else {
			authorizationHeader := r.Header.Get("authorization")
			bearerToken := strings.Split(authorizationHeader, " ")
			accessToken := bearerToken[1]
			val := adminaccess[accessToken]
			if !val {
				w.WriteHeader(http.StatusForbidden)
				resp = &object.Response{
					Message: "POST Permission Denied",
				}
				byteresp, _ := json.Marshal(resp)
				w.Write(byteresp)
				return
			}
			if _, ok := adminaccess[accessToken]; ok {
				resp = database.DeleteMovieByID(id)
			} else {
				w.WriteHeader(http.StatusForbidden)
				resp = &object.Response{
					Message: "Access Denied",
				}
			}
		}
		byteresp, _ := json.Marshal(resp)
		w.Write(byteresp)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := object.Response{
			Message: "Method Not Supported",
		}
		byteresp, _ := json.Marshal(response)
		w.Write(byteresp)
	}
}

//CheckUserLoggedIn Method to check whether Token is coorect or not
func CheckUserLoggedIn(accessToken string) error {
	if accessToken == "" {
		return errors.New("Empty Access Token")
	}
	if settime, k := durationmap[accessToken]; k {
		startTime, _ := time.Parse("2006-01-02 15:04:05.00", settime)
		currentTimeVal := time.Now().Format("2006-01-02 15:04:05.00")
		currentTime, _ := time.Parse("2006-01-02 15:04:05.00", currentTimeVal)
		duration := currentTime.Sub(startTime)
		definedDuration, _ := time.ParseDuration(dur)
		if duration <= definedDuration {
			return nil
		}
		mutx.Lock()
		delete(durationmap, accessToken)
		delete(adminaccess, accessToken)
		delete(accesstokenmap, accessToken)
		mutx.Unlock()
		return errors.New("User Logged Out")
	}
	return errors.New("Invalid Access Token")
}

// CheckURIAndRetrieveID Method checks uri and Fetches Id when Operation By Id is requested
//It allows baseurl/{id} , returns error (method not found) if pattern is not there
//Else returns id
func CheckURIAndRetrieveID(uri []string) (int, error) {
	var id int
	var err error
	if len(uri) > 1 || len(uri) < 1 {
		return 0, errors.New("Method Not Found")
	}
	id, err = strconv.Atoi(uri[0])
	if err != nil {
		return 0, errors.New("Method Not Found")
	}
	return id, nil
}
