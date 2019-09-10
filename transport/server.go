package transport

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"savingDeposits"
	"savingDeposits/auth"
	"time"
)

// Implemented by models that only expose certain fields.
// This method returns a struct with json tags used by
// the transports
type Public interface {
	// The public method returns a struct
	// with json tags.
	Public() interface{}
}

type Server struct {
	Db              *gorm.DB
	router          *mux.Router
	authn           auth.AuthnService
	authz           *auth.AuthzService
	DepositsService savingDeposits.DepositsService
	userService     savingDeposits.UserService
}

// Creates an http server and serves it in the specified address
func (s *Server) ServeHTTP(addr string) error {
	srv := &http.Server{
		Handler:      setCors(s.router),
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv.ListenAndServe()
}

func (s *Server) setupAuthorization() {
	s.authz.AddPermission(savingDeposits.ADMIN, savingDeposits.USERS,
		auth.Create, auth.Read, auth.Update, auth.Delete)
	s.authz.AddPermission(savingDeposits.MANAGER, savingDeposits.USERS,
		auth.Create, auth.Read, auth.Update, auth.Delete)
	s.authz.AddPermission(savingDeposits.ADMIN, savingDeposits.DEPOSITS,
		auth.Create, auth.Read, auth.Update, auth.Delete)
	s.authz.AddPermission(savingDeposits.REGULAR, savingDeposits.DEPOSITS,
		auth.Create, auth.Read, auth.Update, auth.Delete)
}

// Creates GET, POST, PATH and DELETE user handlers.
func (s *Server) AddDepositsHandlers(basePath string) {
	url := "/" + basePath
	urlWithId := url + "/{id:[0-9]+}"

	fmt.Println(url, urlWithId)

	s.router.HandleFunc(url, postSavingsHandler(s.DepositsService)).Methods("POST")
	s.router.HandleFunc(url, getAllSavingsHandler(s.DepositsService)).Methods("GET")
	s.router.HandleFunc(urlWithId, getSavingsHandler(s.DepositsService)).Methods("GET")
	s.router.HandleFunc(urlWithId, patchSavingsHandler(s.DepositsService)).Methods("PATCH")
	s.router.HandleFunc(urlWithId, deleteSavingsHandler(s.DepositsService)).Methods("DELETE")
}

// Creates GET, POST, PATH and DELETE user handlers.
func (s *Server) AddUsersHandlers(basePath string) {
	url := fmt.Sprintf("/%s", basePath)
	urlWithId := fmt.Sprintf("%s/{id:[0-9]+}", url)

	s.router.HandleFunc(url, postUsersHandler(s.userService)).Methods("POST")
	s.router.HandleFunc(url, getAllUsersHandler(s.userService)).Methods("GET")
	s.router.HandleFunc(urlWithId, getUsersHandler(s.userService)).Methods("GET")
	s.router.HandleFunc(urlWithId, patchUsersHandler(s.userService)).Methods("PATCH")
	s.router.HandleFunc(urlWithId, deleteUsersHandler(s.userService)).Methods("DELETE")
}

func NewServer(db *gorm.DB, authNService auth.AuthnService, authZService *auth.AuthzService,
	depositsService savingDeposits.DepositsService, userService savingDeposits.UserService) (*Server, error) {
	router := mux.NewRouter()

	s := &Server{
		Db:              db,
		router:          router,
		authn:           authNService,
		authz:           authZService,
		DepositsService: depositsService,
		userService:     userService,
	}

	// Adds POST, GET, PATCH, DELETE for users
	s.AddUsersHandlers("users")

	// Adds POST, GET, PATCH, DELETE for deposits
	fmt.Println("Adding deposit handlers")
	s.AddDepositsHandlers("deposits")

	// Add other handlers
	router.HandleFunc("/login", s.LoginHandler()).Methods("POST")
	router.HandleFunc("/profile", s.profileHandler()).Methods("GET")
	router.HandleFunc("/newClient", s.newClientHandler()).Methods("POST")

	// Add Authentication/Authorization middleware
	router.Use(s.AuthMiddleware)

	// Add content-type=application/json middleware
	router.Use(s.ContentTypeJsonMiddleware)

	// Log all things
	router.Use(s.LoggingMiddleware)

	// Initialize roles' permissions
	s.setupAuthorization()

	return s, nil
}

func setCors(router *mux.Router) http.Handler {
	allOrigins := handlers.AllowedOrigins([]string{"*"})
	allMethods := handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"})
	allHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	return handlers.CORS(allOrigins, allMethods, allHeaders)(router)
}
