package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "18080"
	}
	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		log.Fatalf("invalid SERVER_PORT: %q", port)
	}

	config := map[string]string{
		"DISPLAY_TITLE":   os.Getenv("DISPLAY_TITLE"),
		"DISPLAY_MESSAGE": os.Getenv("DISPLAY_MESSAGE"),
		"DISPLAY_COUNT":   os.Getenv("DISPLAY_COUNT"),
		"DISPLAY_MODE":    os.Getenv("DISPLAY_MODE"),
		"DISPLAY_TAGS":    os.Getenv("DISPLAY_TAGS"),
		"DISPLAY_COLOR":   os.Getenv("DISPLAY_COLOR"),
		"DISPLAY_NOTE":    os.Getenv("DISPLAY_NOTE"),
		"CONTACTS":        os.Getenv("CONTACTS"),
		"SERVER_PORT":     port,
	}

	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServer(http.Dir("asset")))
	mux.HandleFunc("GET /api/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		if err := json.NewEncoder(w).Encode(config); err != nil {
			log.Printf("write test response: %v", err)
		}
	})
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
