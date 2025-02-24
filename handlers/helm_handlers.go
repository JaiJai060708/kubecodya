package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"github.com/gorilla/mux"
)



func ListHelmCharts(w http.ResponseWriter, r *http.Request) {
	// List charts
	listCmd := exec.Command("helm", "list", "-o", "json")
	output, err := listCmd.CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Listing failed: %s\n%s", err, output), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

func GetHelmChart(w http.ResponseWriter, r *http.Request) {
	// if GET, return the chart
	if r.Method == http.MethodGet {
		// Get parameters using mux
		vars := mux.Vars(r)
		chartName := vars["name"]

		// Get query parameters
		queryParams := r.URL.Query()
		chartOutput := queryParams.Get("output")
		chartType := queryParams.Get("type")

		if chartType == "" {
			chartType = "all"
		}

		var output []byte
		var err error

    if chartOutput == "" {
			readCmd := exec.Command("helm", "get", chartType, chartName)
			output, err = readCmd.CombinedOutput()
		
			if err != nil {
				log.Printf("Error: %s\n%s", err, output)
				http.Error(w, fmt.Sprintf("Installation failed: %s\n%s", err, output), http.StatusInternalServerError)
				return
			}
		
		} else if chartOutput == "json" {
			// Helm + yq processing
			helmCmd := exec.Command("helm", "get", chartType, chartName)
			helmOutput, err := helmCmd.CombinedOutput()
			if err != nil {
				 
				return
			}

			// Process with yq
			yqCmd := exec.Command("yq", "-o", chartOutput)
			yqCmd.Stdin = bytes.NewReader(helmOutput)
			yqOutput, err := yqCmd.CombinedOutput()
			if err != nil {
					return
			}

			// Process with jq
			jqCmd := exec.Command("jq",  "-s", ".")
			jqCmd.Stdin = bytes.NewReader(yqOutput)
			jqOutput, err := jqCmd.CombinedOutput()
			if err != nil {
					return
			}
			output = jqOutput

		} else {
			// Helm + yq processing
			helmCmd := exec.Command("helm", "get", chartType, chartName)
			helmOutput, err := helmCmd.CombinedOutput()
			if err != nil {
				 
				return
			}

			// Process with yq
			yqCmd := exec.Command("yq", "-o", chartOutput)
			yqCmd.Stdin = bytes.NewReader(helmOutput)
			yqOutput, err := yqCmd.CombinedOutput()
			if err != nil {
					return
			}
			output = yqOutput
		
		}
			
		// Set headers and respond
		if chartOutput == "json" {
				w.Header().Set("Content-Type", "application/json")

				 //Validate JSON if requested
				if !json.Valid(output) {
						http.Error(w, fmt.Sprintf("invalid JSON output"), http.StatusInternalServerError)
						return
				}
		}
		

		w.WriteHeader(http.StatusOK)
		w.Write(output)
  }
}

