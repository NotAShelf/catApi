package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var images []string
var logger = logrus.New()
var port string

func init() {
	// Log as JSON instead of the default ASCII formatter
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout (or any other output you prefer)
	logger.SetOutput(os.Stdout)

	// Set the log level (info, warning, error, etc.)
	logger.SetLevel(logrus.InfoLevel)
}

func main() {
	flag.StringVar(&port, "port", "3000", "Port to run the server on")
	flag.Parse()

	// Initialize Viper for configuration management
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv() // Allow environment variables to override config settings

	// Read the configuration file (config.yaml in this example)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading configuration file: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		images := getImages(r) // Call getImages with the request parameter
		io.WriteString(w, "<html><body><div style=\"display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); grid-gap: 10px;\">")
		for i, image := range images {
			io.WriteString(w, `<a href="/api/id?id=`+strconv.Itoa(i)+`">`)
			io.WriteString(w, `<img src="`+filepath.Base(image)+`" style="width: 100%; height: auto;"/>`)
			io.WriteString(w, `</a>`)
		}
		io.WriteString(w, "</div></body></html>")
	})

	http.HandleFunc("/api/id", func(w http.ResponseWriter, r *http.Request) {
		id := sanitizeInput(r.URL.Query().Get("id"))
		if id == "" {
			http.Error(w, "Missing id", http.StatusBadRequest)
			return
		}
		i, err := strconv.Atoi(id)
		if err != nil || i < 0 || i >= len(images) {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}
		http.ServeFile(w, r, images[i])
	})

	http.HandleFunc("/api/list", func(w http.ResponseWriter, r *http.Request) {
		// Create a slice to store image information
		imageList := []map[string]string{}

		for _, image := range images {
			imageInfo := map[string]string{
				"image": image,
				"url":   "/" + filepath.Base(image),
			}
			imageList = append(imageList, imageInfo)
		}

		// Convert the slice to JSON
		jsonData, err := json.Marshal(imageList)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Set the content type to JSON
		w.Header().Set("Content-Type", "application/json")

		// Write the JSON response
		w.Write(jsonData)
	})

	http.HandleFunc("/api/random", func(w http.ResponseWriter, r *http.Request) {
		// Reseed the random number generator to make it truly random on each request
		rand.Seed(time.Now().UnixNano())

		i := rand.Intn(len(images))
		http.ServeFile(w, r, images[i])
	})

	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid API path", http.StatusNotFound)
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getImages(r *http.Request) []string {
	files, err := os.ReadDir("images/")
	if err != nil {
		logger.WithError(err).Fatal("Error reading images directory")
	}
	var images []string
	serverAddress := r.Host // Get the server address from the request
	for _, file := range files {
		imagePath := "http://" + serverAddress + "/images/" + file.Name()
		images = append(images, imagePath)
		logger.Info("Loaded image:", imagePath)
	}
	return images
}

func sanitizeInput(input string) string {
	return template.HTMLEscapeString(input)
}
