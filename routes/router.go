package routes

import (
	"net/http"
	
	"kubecodya/config"
	"kubecodya/handlers"
	
	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	box := packr.New("myBox", config.StaticFilesPath)

	// API routes
	r.HandleFunc("/api/helm-chart/{name}", handlers.GetHelmChart).Methods("GET")
	r.HandleFunc("/api/helm-chart-list", handlers.ListHelmCharts).Methods("GET")

	// Static file routes
	r.PathPrefix("/about").Handler(http.StripPrefix("/about", http.FileServer(box)))
	r.PathPrefix("/charts").Handler(http.StripPrefix("/charts", http.FileServer(box)))
	r.PathPrefix("/").Handler(http.FileServer(box))

	return r
}
