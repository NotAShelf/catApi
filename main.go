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
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // path to look for the config file in
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading configuration file: %v", err)
	}

	port := viper.GetString("server.port")
	flag.Parse()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading configuration file: %v", err)
	}

	images = getImages()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/id", idHandler)
	http.HandleFunc("/api/list", listHandler)
	http.HandleFunc("/api/random", randomHandler)

	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid API path", http.StatusNotFound)
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getImages() []string {
	files, err := os.ReadDir("images/")
	if err != nil {
		logger.WithError(err).Fatal("Error reading images directory")
	}
	var images []string
	for _, file := range files {
		images = append(images, file.Name())
		logger.Info("Loaded image:", file.Name())
	}
	return images
}

func sanitizeInput(input string) string {
	return template.HTMLEscapeString(input)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><body><div style=\"display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); grid-gap: 10px;\">")
	for i := range images {
		io.WriteString(w, `<a href="/api/id?id=`+strconv.Itoa(i)+`">`)
		io.WriteString(w, `<img src="/api/id?id=`+strconv.Itoa(i)+`" style="width: 100%; height: auto;"/>`)
		io.WriteString(w, `</a>`)
	}
	io.WriteString(w, "</div></body></html>")
}

func idHandler(w http.ResponseWriter, r *http.Request) {
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
	http.ServeFile(w, r, "images/"+images[i])
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	imageList := []map[string]string{}
	for i := range images {
		imageInfo := map[string]string{
			"id":  strconv.Itoa(i),
			"url": "/api/id?id=" + strconv.Itoa(i),
		}
		imageList = append(imageList, imageInfo)
	}
	jsonData, err := json.Marshal(imageList)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func randomHandler(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(images))
	http.Redirect(w, r, "/api/id?id="+strconv.Itoa(i), http.StatusSeeOther)
}
