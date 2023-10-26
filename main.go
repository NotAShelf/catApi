package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3005"
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := getRandomImagePath()
		http.ServeFile(w, r, path)
	})

	http.HandleFunc("/api/random", func(w http.ResponseWriter, r *http.Request) {
		path := getRandomImagePath()
		id := strings.Split(path, "-")[1][:len(path)-5]
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":` + id + `, "url":"http://localhost:` + port + `/` + id + `"}`))
	})

	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/")
		path := getImageById(id)
		if path != "" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"id":` + id + `, "url":"http://localhost:` + port + `/` + id + `"}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "Image not found"}`))
		}
	})

	http.HandleFunc("/api/list", func(w http.ResponseWriter, r *http.Request) {
		ids := getAllImageIds()
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[` + strings.Trim(strings.Join(strings.Split(strings.Trim(fmt.Sprint(ids), "[]"), " "), ","), ",") + `]`))
	})

	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getRandomImagePath() string {
	files, _ := os.ReadDir("./images")
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(files))
	return "./images/" + files[randomIndex].Name()
}

func getImageById(id string) string {
	files, _ := os.ReadDir("./images")
	for _, file := range files {
		if strings.HasPrefix(file.Name(), id) {
			return "./images/" + file.Name()
		}
	}
	return ""
}

func getAllImageIds() []int {
	files, _ := os.ReadDir("./images")
	ids := make([]int, 0, len(files))
	for _, file := range files {
		id := strings.Split(file.Name(), "-")[1][:len(file.Name())-5]
		idInt, _ := strconv.Atoi(id)
		ids = append(ids, idInt)
	}
	return ids
}
