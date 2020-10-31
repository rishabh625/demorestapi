package middleware

import (
	"encoding/json"
	"demorestapi/internal/handlers"
	"demorestapi/internal/object"
	"log"
	"net/http"
	"strings"
	"time"
)

// our base middleware implementation.
type service func(http.Handler) http.Handler

//ServiceLoader chain load middleware services.
func ServiceLoader(h http.Handler, svcs ...service) http.Handler {
	for _, svc := range svcs {
		h = svc(h)
	}
	return h
}

//RequestMetrics middleware for request length metrics.
func RequestMetrics(l *log.Logger) service {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			h.ServeHTTP(w, r)
			l.Printf("%s request to %s took %vns.", r.Method, r.URL.Path, time.Now().Sub(start).Nanoseconds())
		})
	}
}

//Auth middleware to AUthenticate Tokens,If auth successfull serves application else return unauthorized
func Auth(l *log.Logger) service {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get("authorization")
			if authorizationHeader != "" {
				bearerToken := strings.Split(authorizationHeader, " ")
				if len(bearerToken) == 2 {
					err := handlers.CheckUserLoggedIn(bearerToken[1])
					if err != nil {
						w.WriteHeader(http.StatusUnauthorized)
						resp := &object.Response{
							Message: err.Error(),
						}
						byteresp, _ := json.Marshal(resp)
						w.Write(byteresp)
						l.Printf("Auth Failed")
					} else {
						h.ServeHTTP(w, r)
					}
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					resp := &object.Response{
						Message: "Authentication Token Not Found",
					}
					byteresp, _ := json.Marshal(resp)
					w.Write(byteresp)
					l.Printf("Auth Failed")
				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				resp := &object.Response{
					Message: "An authorization header is required",
				}
				byteresp, _ := json.Marshal(resp)
				w.Write(byteresp)
				l.Printf("Auth Failed")
			}
		})
	}
}
