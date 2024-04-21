package server

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	cook "github.com/gorilla/sessions"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/b2bstytchapi"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/magiclinks"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/magiclinks/email"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/sessions"
)

var store = cook.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

type StytchServerConfig struct {
	OrganizationID string
}

func Serve(stytchClient *b2bstytchapi.API, conf *StytchServerConfig) {
	// create the router
	router := mux.NewRouter()

	// Register the route
	stytch := NewStytchHandler(stytchClient, conf)

	router.HandleFunc("/", stytch.home)
	router.HandleFunc("/login-or-signup", stytch.LoginOrSignUp).Methods("POST")
	router.HandleFunc("/authenticate", stytch.Authenticate).Methods("GET")

	// Start the server
	http.ListenAndServe(":8010", router)
}

type StytchHandler struct {
	StytchClient *b2bstytchapi.API
	Configs      *StytchServerConfig
}

func NewStytchHandler(s *b2bstytchapi.API, conf *StytchServerConfig) *StytchHandler {
	return &StytchHandler{
		StytchClient: s,
		Configs:      conf,
	}
}

func (h *StytchHandler) home(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "stytch_session")
	if err != nil {
		AuthenticationFailed(w, r)
		return
	}

	JWT := session.Values["jwt"].(string)

	_, err = h.StytchClient.Sessions.AuthenticateJWT(r.Context(), &sessions.AuthenticateJWTParams{
		Body: &sessions.AuthenticateParams{
			SessionJWT: JWT,
		},
	})

	if err != nil {
		AuthenticationFailed(w, r)
		return
	}

	w.Write([]byte("This is my home page"))
}

type LoginOrSignUpRequestParams struct {
	Email string
}

func (h *StytchHandler) LoginOrSignUp(w http.ResponseWriter, r *http.Request) {
	var params LoginOrSignUpRequestParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	_, err := h.StytchClient.MagicLinks.Email.LoginOrSignup(r.Context(), &email.LoginOrSignupParams{
		OrganizationID: h.Configs.OrganizationID,
		EmailAddress:   params.Email,
	})

	if err != nil {
		AuthenticationFailed(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Complete the authentication flow which mints the session
func (h *StytchHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	resp, err := h.StytchClient.MagicLinks.Authenticate(r.Context(), &magiclinks.AuthenticateParams{})

	if err != nil {
		AuthenticationFailed(w, r)
		return
	}

	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, _ := store.Get(r, "stytch_session")
	session.Values["jwt"] = resp.SessionJWT

	err = session.Save(r, w)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}

func AuthenticationFailed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("400 Autentication failed"))
}

func AuthenticationUnauthorized(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("401 Autentication failed"))
}
