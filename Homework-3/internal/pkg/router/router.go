package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"homework-3/config"
	"homework-3/internal/pkg/repository"
	"homework-3/internal/pkg/repository/postrgesql"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const queryParamKey = "key"

type Server struct {
	Repo       *postrgesql.PvzRepo
	AuthConfig config.AuthConfigS
}

type addPvzRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact string `json:"contact"`
}

type addPvzResponse struct {
	ID      int64  `json:"ID"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact string `json:"contact"`
}

type UpdatePvzRequest struct {
	ID      int64  `json:"ID"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact string `json:"contact"`
}

type UpdatePvzRespons struct {
	ID      int64  `json:"ID"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact string `json:"contact"`
}

func CreateRouter(implemetation Server) *mux.Router {
	router := mux.NewRouter()
	router.Use(LoggingMiddleware)
	router.Use(BasicAuthMiddleware(implemetation.AuthConfig))
	router.HandleFunc("/pvz", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			implemetation.Create(w, req)
		case http.MethodPut:
			implemetation.Update(w, req)
		default:
			fmt.Println("error")
		}
	})

	router.HandleFunc(fmt.Sprintf("/pvz/{%s:[0-9]+}", queryParamKey), func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			implemetation.GetByID(w, req)
		case http.MethodDelete:
			implemetation.DeleteByID(w, req)
		default:
			fmt.Println("error")
		}
	})
	return router
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
			log.Printf("Received %s request: URL=%s, RemoteAddr=%s", r.Method, r.URL.String(), r.RemoteAddr)
		}
		next.ServeHTTP(w, r)
	})
}

func BasicAuthMiddleware(auth config.AuthConfigS) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if !ok || username != auth.Username || password != auth.Password {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (s *Server) Create(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var unm addPvzRequest
	if err = json.Unmarshal(body, &unm); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pvzRepo := &repository.Pvz{
		Name:    unm.Name,
		Address: unm.Address,
		Contact: unm.Contact,
	}
	id, err := s.Repo.Add(req.Context(), pvzRepo)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp := &addPvzResponse{
		ID:      id,
		Name:    pvzRepo.Name,
		Address: pvzRepo.Address,
		Contact: pvzRepo.Contact,
	}
	pvzJson, _ := json.Marshal(resp)
	w.Write(pvzJson)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetByID(w http.ResponseWriter, req *http.Request) {
	key, ok := mux.Vars(req)[queryParamKey]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pvz, err := s.Repo.GetByID(req.Context(), keyInt)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			fmt.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pvzJson, _ := json.Marshal(pvz)
	w.Write(pvzJson)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) Update(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var pvzRepo repository.Pvz
	if err = json.Unmarshal(body, &pvzRepo); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.Repo.Update(req.Context(), &pvzRepo)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			fmt.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pvzJson, _ := json.Marshal(pvzRepo)
	w.Write(pvzJson)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) DeleteByID(w http.ResponseWriter, req *http.Request) {
	key, ok := mux.Vars(req)[queryParamKey]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.Repo.DeleteByID(req.Context(), keyInt)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			fmt.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
