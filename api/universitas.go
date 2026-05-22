package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type University struct {
	Name     string   `json:"name"`
	Country  string   `json:"country"`
	WebPages []string `json:"web_pages"`
	Domains  []string `json:"domains"`
}

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

// Handler is the serverless function entrypoint for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	data, err := fetchUniversities()
	if err != nil {
		http.Error(w, "Gagal mengambil data universitas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// Allow CORS if needed
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(data)
}
