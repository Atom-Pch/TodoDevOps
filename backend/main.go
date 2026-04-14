package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func setJSONHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
}

func main() {
	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		setJSONHeaders(w)

		if r.Method == http.MethodOptions {
			return
		}

		resp := Response{Message: "Hello from Go!"}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonResp)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		setJSONHeaders(w)

		resp := Response{Message: "Health Check OK"}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonResp)
	})

	fmt.Println("Server listening on port 8888")
	http.ListenAndServe(":8888", nil)
}
