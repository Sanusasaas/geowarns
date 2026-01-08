package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	port := flag.String("port", "9090", "port to listen on")
	flag.Parse()

	mux := http.NewServeMux()

	// Обработчик вебхуков
	mux.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		fmt.Printf("\n[Webhook Received]\nTime: %s\nPayload: %+v\n",
			time.Now().Format(time.RFC3339),
			payload)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "received"}`))
	})

	// Обработчик проверки состояния
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	addr := fmt.Sprintf(":%s", *port)
	fmt.Printf("Webhook mock server running on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}