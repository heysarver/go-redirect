package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

type RedirectConfig struct {
	Hosts []struct {
		Hostname       string `yaml:"hostname"`
		DestinationURL string `yaml:"destinationURL"`
		StatusCode     int    `yaml:"statusCode"`
	} `yaml:"hosts"`
}

func readConfig() (*RedirectConfig, error) {
	var config RedirectConfig

	// Try to read from config.yaml
	data, err := ioutil.ReadFile("config.yaml")
	if err == nil {
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, err
		}
		return &config, nil
	}

	// If config.yaml is not found, use environment variables
	statusCode, convErr := strconv.Atoi(os.Getenv("STATUS_CODE"))
	if convErr != nil {
		statusCode = 302 // Default to 302 if STATUS_CODE is not a valid integer
	}

	config.Hosts = append(config.Hosts, struct {
		Hostname       string `yaml:"hostname"`
		DestinationURL string `yaml:"destinationURL"`
		StatusCode     int    `yaml:"statusCode"`
	}{
		Hostname:       os.Getenv("HOSTNAME"),
		DestinationURL: os.Getenv("DESTINATION_URL"),
		StatusCode:     statusCode,
	})

	return &config, nil
}

func redirectHandler(config *RedirectConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, host := range config.Hosts {
			if r.Host == host.Hostname {
				http.Redirect(w, r, host.DestinationURL, host.StatusCode)
				return
			}
		}
		http.Error(w, "Host not found", http.StatusNotFound)
	}
}

func main() {
	config, err := readConfig()
	if err != nil {
		fmt.Println("Error reading configuration:", err)
		return
	}

	http.HandleFunc("/", redirectHandler(config))
	fmt.Println("Server is running on port 8300...")
	if err := http.ListenAndServe(":8300", nil); err != nil {
		panic(err)
	}
}
