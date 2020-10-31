package database

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"demorestapi/internal/object"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//Signup Adds User,Admin to database, gives success message or error as message
func Signup(request *object.SignupRequest) *object.Response {
	db = GetConnection()
	if request.Username == "" || request.Password == "" || request.Email == "" {
		return &object.Response{
			Message: "Username, Password,Email is mandatory",
		}
	}
	var admin_secret bool
	if request.AdminSecret == "admin" {
		admin_secret = true
	}
	password := fmt.Sprintf("%x", md5.Sum([]byte(request.Password)))
	insForm, err := db.Prepare("INSERT INTO imdb_users (admin,deleted,username,password,email) values ($1, $2, $3, $4, $5)")
	if err != nil {
		return &object.Response{
			Message: err.Error(),
		}
	}
	_, err = insForm.Exec(admin_secret, false, request.Username, password, request.Email)
	if err != nil {
		return &object.Response{
			Message: err.Error(),
		}
	}
	return &object.Response{
		Message: "User Created Successfully",
	}
}

//Login Authenticates User Credentials from db, gives success message,access token on login or error as message
func Login(request *object.LoginRequest) (*object.LoginResponse, bool) {
	db = GetConnection()
	if request.Username == "" || request.Password == "" {
		return &object.LoginResponse{
			Message: "Username, Password is mandatory",
		}, false
	}
	query := "Select admin,username from imdb_users where username = $1 and password = $2 and deleted = $3"
	password := fmt.Sprintf("%x", md5.Sum([]byte(request.Password)))
	var admin bool
	var username string
	rows, err := db.Query(query, request.Username, password, false)
	if err != nil {
		return &object.LoginResponse{
			Message: err.Error(),
		}, false
	}
	if rows.Next() {
		rows.Scan(&admin, &username)
		key := fmt.Sprintf("%s%s", strconv.FormatInt(time.Now().UnixNano(), 10), username)
		accesstoken := md5.Sum([]byte(key))
		return &object.LoginResponse{
			Message:     "User Logged In",
			AccessToken: fmt.Sprintf("%x", accesstoken),
		}, admin
	}
	return &object.LoginResponse{
		Message: "Failed to Login",
	}, false
}

//AddMovies Adds Movie into Database, Returns Success Message or Error as Message
func AddMovies(request *object.AddMovieRequest, user string) *object.Response {
	db = GetConnection()
	if reflect.DeepEqual(request, &object.AddMovieRequest{}) || reflect.DeepEqual(request.Movies, []object.Movie{}) {
		return &object.Response{
			Message: "No Movies to Add",
		}
	}
	for _, movie := range request.Movies {
		jsondata, _ := json.Marshal(movie)
		insForm, err := db.Prepare("INSERT INTO movies (data,added_by) values ($1,$2)")
		if err != nil {
			return &object.Response{
				Message: err.Error(),
			}
		}
		_, err = insForm.Exec(jsondata, user)
		if err != nil {
			return &object.Response{
				Message: err.Error(),
			}
		}
	}
	return &object.Response{
		Message: "Movies Added",
	}
}

//FetchAllMovies Fetches All Movies from Db and returns list of movies and success message or error as message
func FetchAllMovies() *object.GetAllMovieResponse {
	db = GetConnection()
	var moviebytes []byte
	var movie object.Movie
	var movies []object.Movie
	var id int
	query := "Select id,data from movies"
	rows, err := db.Query(query)
	if err != nil {
		return &object.GetAllMovieResponse{
			Message: err.Error(),
		}
	}
	for rows.Next() {
		err := rows.Scan(&id, &moviebytes)
		if err != nil {
			return &object.GetAllMovieResponse{
				Message: err.Error(),
			}
		}
		err = json.Unmarshal(moviebytes, &movie)
		movie.Id = id
		if err != nil {
			return &object.GetAllMovieResponse{
				Message: err.Error(),
			}
		}
		movies = append(movies, movie)
	}
	return &object.GetAllMovieResponse{
		Message: "SuccessFully Retrieved All Movies",
		Movies:  movies,
	}
}

//FetchMoviesByID Fetches All Movies from Db by Id and returns list of 1 movie and success message or error as message
func FetchMoviesByID(id int) *object.GetAllMovieResponse {
	db = GetConnection()
	var moviebytes []byte
	var movie object.Movie
	var movies []object.Movie
	query := "Select data from movies where id = $1"
	row := db.QueryRow(query, id)
	err := row.Scan(&moviebytes)
	if err != nil {
		return &object.GetAllMovieResponse{
			Message: err.Error(),
		}
	}
	err = json.Unmarshal(moviebytes, &movie)
	if err != nil {
		return &object.GetAllMovieResponse{
			Message: err.Error(),
		}
	}
	movies = append(movies, movie)
	return &object.GetAllMovieResponse{
		Message: "SuccessFully Retrieved Movie",
		Movies:  movies,
	}
}

//DeleteMovieByID from Db for specified id and returns success message with count or error as message
func DeleteMovieByID(id int) *object.Response {
	db = GetConnection()
	delquery, err := db.Prepare("delete from movies where id =$1 ")
	if err != nil {
		return &object.Response{
			Message: err.Error(),
		}
	}
	res, err := delquery.Exec(id)
	if err != nil {
		return &object.Response{
			Message: err.Error(),
		}
	}
	count, _ := res.RowsAffected()
	return &object.Response{
		Message: fmt.Sprintf("%d %s", count, "Documents Deleted Successfully"),
	}
}

//UpdateMovieByID from Db for id specified and returns success message with count or error as message
func UpdateMovieByID(req *object.UpdateMovie, id int, user string) *object.Response {
	db = GetConnection()
	switch req.Key {
	case "genre":
		var val string
		arr := strings.Split(req.Value, ",")
		for _, i := range arr {
			tempval := fmt.Sprintf("\"%s\",", i)
			val += tempval
		}
		req.Value = "[" + strings.TrimSuffix(val, ",") + "]"
	case "director":
		fallthrough
	case "name":
		req.Value = "\"" + req.Value + "\""
	case "imdb_score":
		fallthrough
	case "99popularity":
	default:
		return &object.Response{
			Message: "Invalid Key",
		}
	}

	updateQuery, err := db.Prepare("update movies set data = jsonb_set(data,'{" + req.Key + "}','" + req.Value + "',true),updated_by = $1 where id = $2")
	if err != nil {
		return &object.Response{
			Message: err.Error(),
		}
	}
	res, err := updateQuery.Exec(user, id)
	if err != nil {
		return &object.Response{
			Message: err.Error(),
		}
	}
	count, _ := res.RowsAffected()
	return &object.Response{
		Message: fmt.Sprintf("%d %s", count, " Document Updated Successfully"),
	}
}

//SearchMovies from Db for given term or range query or both  and returns movies list with success message or error as message
func SearchMovies(request *object.SearchMovie) *object.GetAllMovieResponse {
	db = GetConnection()
	var moviebytes []byte
	var movie object.Movie
	var movies []object.Movie
	var id int
	if !reflect.DeepEqual(request.Term, object.Movie{}) && reflect.DeepEqual(request.Range, []object.Range{}) {
		t, _ := json.Marshal(request.Term)
		term := string(t)
		query := "select id,data from movies where (data)::jsonb @> $1::jsonb"
		rows, err := db.Query(query, term)
		if err != nil {
			return &object.GetAllMovieResponse{
				Message: err.Error(),
			}
		}
		for rows.Next() {
			err := rows.Scan(&id, &moviebytes)
			if err != nil {
				return &object.GetAllMovieResponse{
					Message: err.Error(),
				}
			}
			err = json.Unmarshal(moviebytes, &movie)
			movie.Id = id
			if err != nil {
				return &object.GetAllMovieResponse{
					Message: err.Error(),
				}
			}
			movies = append(movies, movie)
		}
	} else if !reflect.DeepEqual(request.Range, []object.Range{}) && reflect.DeepEqual(request.Term, object.Movie{}) {
		query := "select id,data from movies where"
		for _, rang := range request.Range {
			var op string
			switch rang.Operator {
			case "lt":
				op = "<"
			case "lte":
				op = "<="
			case "gt":
				op = "gt"
			case "gte":
				op = ">="
			default:
				op = "="
			}
			num := fmt.Sprintf("%f", rang.Value)
			query += " ((data->'" + rang.Key + "')::jsonb::NUMERIC " + op + num + ") and "
		}
		query = strings.TrimSuffix(query, " and ")
		rows, err := db.Query(query)
		if err != nil {
			return &object.GetAllMovieResponse{
				Message: err.Error(),
			}
		}
		for rows.Next() {
			err := rows.Scan(&id, &moviebytes)
			if err != nil {
				return &object.GetAllMovieResponse{
					Message: err.Error(),
				}
			}
			err = json.Unmarshal(moviebytes, &movie)
			movie.Id = id
			if err != nil {
				return &object.GetAllMovieResponse{
					Message: err.Error(),
				}
			}
			movies = append(movies, movie)
		}
	} else if !reflect.DeepEqual(request.Range, []object.Range{}) && !reflect.DeepEqual(request.Term, object.Movie{}) {
		t, _ := json.Marshal(request.Term)
		term := string(t)
		query := "select id,data from movies where ((data)::jsonb @> $1::jsonb) and "
		for _, rang := range request.Range {
			var op string
			switch rang.Operator {
			case "lt":
				op = "<"
			case "lte":
				op = "<="
			case "gt":
				op = ">"
			case "gte":
				op = ">="
			default:
				op = "="
			}
			num := fmt.Sprintf("%f", rang.Value)
			query += " ((data->'" + rang.Key + "')::jsonb::NUMERIC " + op + num + ") and "
		}
		query = strings.TrimSuffix(query, " and ")
		rows, err := db.Query(query, term)
		if err != nil {
			return &object.GetAllMovieResponse{
				Message: err.Error(),
			}
		}
		for rows.Next() {
			err := rows.Scan(&id, &moviebytes)
			if err != nil {
				return &object.GetAllMovieResponse{
					Message: err.Error(),
				}
			}
			err = json.Unmarshal(moviebytes, &movie)
			movie.Id = id
			if err != nil {
				return &object.GetAllMovieResponse{
					Message: err.Error(),
				}
			}
			movies = append(movies, movie)
		}
	} else {
		return &object.GetAllMovieResponse{
			Message: "Invalid Search Query",
			Movies:  movies,
		}
	}
	return &object.GetAllMovieResponse{
		Message: "SuccessFully Retrieved Movie",
		Movies:  movies,
	}
}
