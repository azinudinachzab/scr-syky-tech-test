package delivery

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/azinudinachzab/scr-syky-tech-test/model"
	"github.com/azinudinachzab/scr-syky-tech-test/pkg/errs"
	"github.com/azinudinachzab/scr-syky-tech-test/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

type HttpServer struct {
	service service.Service
}

func NewHttpServer(svc service.Service) http.Handler {
	r := chi.NewRouter()
	d := &HttpServer{
		service: svc,
	}

	/* ***** ***** *****
	 * init middleware
	 * ***** ***** *****/
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"*"}, // "True-Client-IP", "X-Forwarded-For", "X-Real-IP", "X-Request-Id", "Origin", "Accept", "Content-Type", "Authorization", "Token"
		AllowCredentials: true,
		MaxAge:           86400,
	}))
	r.Use(httprate.LimitByIP(80, 1*time.Minute))
	r.Use(middleware.CleanPath)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	/* ***** ***** *****
	 * init custom error for 404 and 405
	 * ***** ***** *****/
	r.NotFound(d.Custom404)
	r.MethodNotAllowed(d.Custom405)

	/* ***** ***** *****
	 * init path route
	 * ***** ***** *****/
	r.Get("/", d.Home)

	// onboarding
	r.Post("/registration", d.Registration)

	return r
}

func (d *HttpServer) Home(w http.ResponseWriter, r *http.Request) {
	responseData(w, r, httpResponse{Message: "Hello World : " + time.Now().Format(time.RFC3339)})
}

func (d *HttpServer) Registration(w http.ResponseWriter, r *http.Request) {
	var req model.RegistrationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = errs.New(model.ECodeBadRequest, "failed to decode request body")
		responseError(w, r, err)
		return
	}

	err := d.service.Registration(r.Context(), req)
	if err != nil {
		responseError(w, r, err)
		return
	}

	responseData(w, r, httpResponse{Message: "Registration Success"})
}

func (d *HttpServer) Custom404(w http.ResponseWriter, r *http.Request) {
	err := errs.New(model.ECodeNotFound, "route does not exist")
	responseError(w, r, err)
}

func (d *HttpServer) Custom405(w http.ResponseWriter, r *http.Request) {
	err := errs.New(model.ECodeMethodFail, "method is not valid")
	responseError(w, r, err)
}
