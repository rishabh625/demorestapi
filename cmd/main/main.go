package main

import (
	"flag"
	"fmt"
	. "fyndtest/internal/database"
	"fyndtest/internal/handlers"
	. "fyndtest/internal/middleware"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	port    int
	logfile string
	ver     bool
)

//Init To Initialize Cache Maps and flags require on startup
func init() {
	flag.IntVar(&port, "port", 9090, "The port to listen on.")
	flag.StringVar(&logfile, "logfile", "", "Location of the logfile.")
	flag.BoolVar(&ver, "version", false, "Print server version.")
	handlers.InitializeMaps()
}

const (
	// base HTTP paths.
	apiVersion  = "v1"
	apiBasePath = "/api/" + apiVersion + "/"

	//http path .
	signupPath    = apiBasePath + "signup"
	loginPath     = apiBasePath + "login"
	moviesPath    = apiBasePath + "movies/"
	addmoviesPath = apiBasePath + "movies"
	searchPath    = apiBasePath + "movies/search"
	// server version.
	version = "1.0.0"
)

//Main Function: Starts Server,Exposes Endpoint and Initialized DataBase Connection
func main() {
	flag.Parse()
	if ver {
		fmt.Printf("HTTP Server v%s", version)
		os.Exit(0)
	}
	var logger *log.Logger
	if logfile == "" {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	} else {
		f, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}
		logger = log.New(f, "", log.LstdFlags)
	}
	err := InitConnection()
	if err != nil {
		panic(err)
	}
	http.Handle(signupPath, ServiceLoader(http.HandlerFunc(handlers.Signup), RequestMetrics(logger)))
	http.Handle(loginPath, ServiceLoader(http.HandlerFunc(handlers.Login), RequestMetrics(logger)))
	http.Handle(moviesPath, ServiceLoader(http.HandlerFunc(handlers.Movies), RequestMetrics(logger), Auth(logger)))
	http.Handle(searchPath, ServiceLoader(http.HandlerFunc(handlers.Search), RequestMetrics(logger), Auth(logger)))
	http.Handle(addmoviesPath, ServiceLoader(http.HandlerFunc(handlers.Movies), RequestMetrics(logger), Auth(logger)))
	logger.Printf("starting server on :%d", port)

	strPort := ":" + strconv.Itoa(port)
	logger.Fatal(http.ListenAndServe(strPort, nil))
	logger.Printf("started server on :%d", port)
}
