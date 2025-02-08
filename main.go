package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var images []string
var logger = logrus.New()
var port string

// Okay I admit, this is bad. Just a workaround for now, until I figure out a clean
// way of displaying all images in a grid.
var tmpl = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Image Gallery</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; }
        .gallery {
            column-count: 5;
            column-gap: 8px;
        }
        .gallery a {
            display: inline-block;
            width: 100%;
            margin-bottom: 8px;
        }
        .gallery img {
            width: 100%;
            height: auto;
            border-radius: 8px;
            display: block;
        }
        @media (max-width: 1024px) { .gallery { column-count: 4; } }
        @media (max-width: 768px) { .gallery { column-count: 3; } }
        @media (max-width: 480px) { .gallery { column-count: 2; } }
    </style>
</head>
<body>
    <h1>Image Gallery</h1>
    <div class="gallery">
        {{range $index, $img := .Images}}
        <a href="/api/id?id={{$index}}">
            <img src="/api/id?id={{$index}}" alt="Image {{$index}}">
        </a>
        {{end}}
    </div>
</body>
</html>`))

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

	port = viper.GetString("server.port")

	flag.Parse()
	images = getImages()

	// Add request logging middleware
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/api/id", idHandler)
	mux.HandleFunc("/api/list", listHandler)
	mux.HandleFunc("/api/random", randomHandler)
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid API path", http.StatusNotFound)
	})

	// Wrap the mux with the logging middleware
	http.Handle("/", logRequest(mux))

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getImages() []string {
	files, err := os.ReadDir("images/")
	if err != nil {
		logger.WithError(err).Fatal("Error reading images directory")
	}

	if len(files) == 0 {
		logger.Warn("No images found in the images directory")
	}

	var images []string
	for _, file := range files {
		images = append(images, file.Name())
		logger.Info("Loaded image:", file.Name())
	}
	return images
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, struct {
		Images []string
	}{Images: images})
}

func idHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	i, err := strconv.Atoi(id)
	if err != nil || i < 0 || i >= len(images) {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	imagePath := "images/" + images[i]
	if !isValidImagePath(imagePath) {
		http.Error(w, "Invalid image path", http.StatusBadRequest)
		return
	}

	http.ServeFile(w, r, imagePath)
}

func isValidImagePath(path string) bool {
	if !strings.HasPrefix(path, "images/") {
		return false
	}
	return true
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

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger.Infof("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		logger.Infof("Completed %s %s in %v", r.Method, r.URL.Path, duration)
	})
}

func randomHandler(w http.ResponseWriter, r *http.Request) {
	source := rand.NewSource(time.Now().UnixNano())
	rand.New(source)

	i := rand.Intn(len(images))
	http.Redirect(w, r, "/api/id?id="+strconv.Itoa(i), http.StatusSeeOther)
}
