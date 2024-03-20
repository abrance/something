package main

import "net/http"

type Response struct {
	Value int `json:"value"`
}

var response Response = Response{}

//func main() {
//	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Content-Type", "application/json")
//		json.NewEncoder(w).Encode(response)
//	})
//
//	http.ListenAndServe("0.0.0.0:7000", nil)
//}

func main() {
	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(""))
	})
	http.ListenAndServe("0.0.0.0:7000", nil)
}
