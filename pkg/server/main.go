package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/b2bstytchapi"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/sessions"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/sso"
	"github.com/stytchauth/stytch-go/v12/stytch/stytcherror"
)

type StytchServerConfig struct {
	OrganizationID string
	ConnectionID   string
	PublicToken    string
}

func Serve(stytchClient *b2bstytchapi.API, conf *StytchServerConfig) {
	// create the router
	router := mux.NewRouter()

	// Register the route
	stytch := NewStytchHandler(stytchClient, conf)

	router.HandleFunc("/", stytch.home)
	router.HandleFunc("/authenticate", stytch.authenticate).Methods("GET")
	router.HandleFunc("/can-i", stytch.canI).Methods("GET")

	// Start the server
	http.ListenAndServe(":8010", router)
}

// StytchHandler implement the Backend Integration of SSO
// see https://stytch.com/docs/b2b/guides/sso/backend
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

// home is our root page, we used it to test the authentication flow
// if stytch_session is not set we redirect to SSO login
// else use stytch_session JWT to authenticate
// if the auth succeed we desplay user's metadata
func (h *StytchHandler) home(w http.ResponseWriter, r *http.Request) {
	// fetch stytch_session or redirect to SSO login
	session, err := r.Cookie("stytch_session")
	if err != nil {
		h.RedirectToSSO(w, r)
		return
	}

	// validate stytch_session or return failed authentication
	// There are three variante of this validation (authenticateSession, authenticateJWT, authenticateJWTLocal)
	// ref: https://stytch.com/docs/b2b/api/authenticate-session
	metdata, err := h.StytchClient.Sessions.AuthenticateJWT(r.Context(), &sessions.AuthenticateJWTParams{
		Body: &sessions.AuthenticateParams{
			SessionJWT: session.Value,
		},
	})
	if err != nil {
		AuthenticationFailed(w, r)
		return
	}

	// Json serialization of member metadata
	member, err := json.Marshal(metdata.Member)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(member)
}

// Authenticate handles Stytch callback after Login or Signup
// The route to this handler must match the configured RedirectURL
func (h *StytchHandler) authenticate(w http.ResponseWriter, r *http.Request) {
	resp, err := h.StytchClient.SSO.Authenticate(r.Context(), &sso.AuthenticateParams{
		SSOToken: r.URL.Query().Get("token"),
	})

	if err != nil {
		AuthenticationFailed(w, r)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "stytch_session",
		Value: resp.SessionJWT,
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// canI allows us to test the AuthorizationCheck feature
// we can ask if we can perform a given action on a resource
func (h *StytchHandler) canI(w http.ResponseWriter, r *http.Request) {
	// fetch stytch_session or redirect to SSO login
	session, err := r.Cookie("stytch_session")
	if err != nil {
		h.RedirectToSSO(w, r)
		return
	}

	// Read Query Parameters
	resource := r.URL.Query().Get("resource")
	if resource == "" {
		BadRequestHandler(w, r)
		return
	}
	action := r.URL.Query().Get("action")
	if action == "" {
		BadRequestHandler(w, r)
		return
	}

	// Perform authentication with AuthorizationCheck
	metadata, err := h.StytchClient.Sessions.AuthenticateJWT(r.Context(), &sessions.AuthenticateJWTParams{
		Body: &sessions.AuthenticateParams{
			SessionJWT: session.Value,

			AuthorizationCheck: &sessions.AuthorizationCheck{
				OrganizationID: h.Configs.OrganizationID,
				ResourceID:     resource,
				Action:         action,
			},
		},
	})

	// handle Auntentication or Authorization errors
	if err != nil {
		stytchErr := err.(stytcherror.Error)
		w.WriteHeader(stytchErr.StatusCode)
		w.Write([]byte(stytchErr.ErrorMessage))
		return
	}

	// Json serialization of verdict metadata
	verdict, err := json.Marshal(metadata.Verdict)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(verdict)
}

func (h *StytchHandler) RedirectToSSO(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/v1/public/sso/start?connection_id=%s&public_token=%s",
		h.StytchClient.SSO.C.GetConfig().BaseURI, h.Configs.ConnectionID, h.Configs.PublicToken,
	)

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}

func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("400 Bad Request"))
}

func AuthorisationFailed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("403 Forbidden: invalid role or permissions"))
}

func AuthenticationFailed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("401 invalid authentication credentials"))
}
