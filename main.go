package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type University struct {
	Name     string   `json:"name"`
	Country  string   `json:"country"`
	WebPages []string `json:"web_pages"`
	Domains  []string `json:"domains"`
}

// Ambil data dari API publik
func fetchUniversities() ([]University, error) {
	url := "http://universities.hipolabs.com/search?country=Indonesia"

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	var data []University
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

// Handler endpoint API lokal
func listUniHandler(w http.ResponseWriter, r *http.Request) {
	data, err := fetchUniversities()
	if err != nil {
		http.Error(w, "Gagal mengambil data universitas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

func main() {
	r := mux.NewRouter()

	// Endpoint API
	r.HandleFunc("/api/universitas", listUniHandler).Methods("GET")

	// Serve file statis (HTML, CSS, JS)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./templates/")))

	addr := ":8080"
	fmt.Println("Server berjalan di http://localhost" + addr)
	log.Fatal(http.ListenAndServe(addr, r))
}