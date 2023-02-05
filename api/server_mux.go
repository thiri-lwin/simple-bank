package api

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	db "github.com/thiri-lwin/thiri-bank/db/sqlc"
)

type MuxServer struct {
	store  *db.Store
	router *mux.Router
}

func NewMuxServer(store *db.Store) *MuxServer {
	server := &MuxServer{
		store: store,
	}
	router := mux.NewRouter()
	server.router = router
	server.setRoutes()
	return server
}

func (s *MuxServer) setRoutes() {
	s.router.HandleFunc("/accounts", s.createAccount).Methods("POST")
}

func (s *MuxServer) StartMuxServer(address string) error {
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	return http.ListenAndServe(address, handlers.CORS(allowedOrigins)(s.router))
}

type response struct {
	Error string
}

func respondJson(err error) *response {
	return &response{Error: err.Error()}
}
