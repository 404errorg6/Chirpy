package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	erroHandler := func(msg string) {
		w.WriteHeader(400)
		data, _ := json.Marshal(struct {
			Error string `json:"error"`
		}{
			Error: msg,
		})
		w.Write(data)
	}
	var input struct {
		Body string `json:"body"`
	}
	var output struct {
		CleanedBody string `json:"cleaned_body"`
		Valid       bool   `json:"valid"`
	}

	body, err := io.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		erroHandler("error reading request")
		return
	}

	err = json.Unmarshal(body, &input)
	if err != nil {
		erroHandler("error unmarhsalling data")
		return
	}

	if len(body) > 140 {
		erroHandler("Chirp is too long")
		return
	}

	w.WriteHeader(http.StatusOK)
	output.Valid = true
	output.CleanedBody = cleanBody(input.Body)
	data, err := json.Marshal(output)
	if err != nil {
		erroHandler("error marshalling data")
		return
	}
	w.Write(data)
	w.Write([]byte("\n"))
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiCfg) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	req.Header.Set("Content-Type", "text/html")
	html := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileServerHits.Load())
	w.Write([]byte(html))
}

func (cfg *apiCfg) handlerReset(w http.ResponseWriter, req *http.Request) {
	cfg.fileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0\n"))
}
