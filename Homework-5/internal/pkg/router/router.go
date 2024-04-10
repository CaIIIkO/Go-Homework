package router

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"homework-3/config"
	"homework-3/internal/pkg/repository"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const queryParamKey = "key"

type Server struct {
	Repo       repository.PvzRepo
	AuthConfig config.AuthConfigS
}

type ServerI interface {
	Create(w http.ResponseWriter, req *http.Request)
	GetByID(w http.ResponseWriter, req *http.Request)
	Update(w http.ResponseWriter, req *http.Request)
	DeleteByID(w http.ResponseWriter, req *http.Request)
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

	data, status := s.crt(req.Context(), pvzRepo)
	w.Write(data)
	w.WriteHeader(status)
}

func (s *Server) crt(ctx context.Context, repo *repository.Pvz) ([]byte, int) {
	id, err := s.Repo.Add(ctx, repo)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	resp := &addPvzResponse{
		ID:      id,
		Name:    repo.Name,
		Address: repo.Address,
		Contact: repo.Contact,
	}
	pvzJson, err := json.Marshal(resp)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return pvzJson, http.StatusOK
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

	data, status := s.get(req.Context(), keyInt)
	w.WriteHeader(status)
	w.Write(data)
}

func (s *Server) get(ctx context.Context, key int64) ([]byte, int) {
	pvz, err := s.Repo.GetByID(ctx, key)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, http.StatusNotFound
		}
		return nil, http.StatusInternalServerError
	}
	pvzJson, err := json.Marshal(pvz)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return pvzJson, http.StatusOK
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

	data, status := s.upd(req.Context(), &pvzRepo)
	w.Write(data)
	w.WriteHeader(status)
}

func (s *Server) upd(ctx context.Context, repo *repository.Pvz) ([]byte, int) {
	err := s.Repo.Update(ctx, repo)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, http.StatusNotFound
		}
		return nil, http.StatusInternalServerError
	}
	pvzJson, err := json.Marshal(&repo)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return pvzJson, http.StatusOK
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

	status := s.del(req.Context(), keyInt)
	w.WriteHeader(status)
}

func (s *Server) del(ctx context.Context, key int64) int {
	err := s.Repo.DeleteByID(ctx, key)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return http.StatusNotFound
		}
		return http.StatusInternalServerError
	}
	return http.StatusOK
}
