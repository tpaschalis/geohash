package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type batchDecodeResponse struct {
	Responses []decodeResponse
}

type decodeResponse struct {
	InputGeohash string  `json:"inputGeohash"`
	Lat          float64 `json:"lat"`
	Lon          float64 `json:"lon"`
	BBox         Box     `json:"boundingBox"`
}

type encodeResponse struct {
	Query     string  `json:"query"`
	InputLat  float64 `json:"lat"`
	InputLon  float64 `json:"lon"`
	Precision int     `json:"precision"`
	Geohash   string  `json:"geohash"`
}

type encodeBatchInput struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Pre int     `json:"precision"`
}

type encodeBatchResponse struct {
	Query   string           `json:"query"`
	Results []encodeResponse `json:"results"`
}

type decodeBatchInput struct {
	Geohash string `json:"geohash"`
}

type decodeBatchResponse struct {
	Query   string           `json:"query"`
	Results []decodeResponse `json:"results"`
}

type errorResponse struct {
	Query  string `json:"query"`
	Errors string `json:"errors"`
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexPage)
	router.HandleFunc("/v1/health", healthCheckHandler).Methods("GET")
	router.HandleFunc("/v1/encode/{lat:[+-]{0,1}[0-9]+.{0,1}[0-9]+},{lon:[+-]{0,1}[0-9]+.{0,1}[0-9]+}", encodeLonLat).Queries("pre", "{pre}").Methods("GET")
	router.HandleFunc("/v1/encode/{lat:[+-]{0,1}[0-9]+.{0,1}[0-9]+},{lon:[+-]{0,1}[0-9]+.{0,1}[0-9]+}", encodeLonLat).Methods("GET")
	router.HandleFunc("/v1/decode/{hash:[0123456789bcdefghjkmnpqrstuvwxyz]{1,12}}", decodeGeohash).Methods("GET")
	router.HandleFunc("/v1/batchEncode", encodeBatchLonLat).Methods("POST")
	router.HandleFunc("/v1/batchDecode", decodeBatchGeohash).Methods("POST")
	router.NotFoundHandler = http.HandlerFunc(route404Handler)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func encodeLonLat(w http.ResponseWriter, r *http.Request) {

	var lat, lon float64
	var latErr, lonErr, err error
	var hash string
	var pre int
	lat, latErr = strconv.ParseFloat(mux.Vars(r)["lat"], 64)
	lon, lonErr = strconv.ParseFloat(mux.Vars(r)["lon"], 64)

	query, _ := r.URL.Parse(r.URL.RequestURI())
	fmt.Println(query)

	if latErr != nil && lonErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{
			query.String(),
			"Provided Latitude or Longitude could not be parsed as a float",
		})
		return
	}

	preInput, found := mux.Vars(r)["pre"]

	if !found {
		pre = 12
		hash, err = EncodeUsingPrecision(Point{lat, lon}, pre)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{
				query.String(),
				"The request returned an error : " + err.Error(),
			})
			return
		}
	} else {
		pre, err = strconv.Atoi(preInput)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{
				query.String(),
				"Could not parse provided precision as a valid integer",
			})
			return
		}
		hash, err = EncodeUsingPrecision(Point{lat, lon}, pre)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{
				query.String(),
				"The request returned an error : " + err.Error(),
			})
			return
		}
	}

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{
			query.String(),
			"The request returned an error : " + err.Error(),
		})
		return
	}
	fmt.Println(pre)

	res := &encodeResponse{query.String(), lat, lon, pre, hash}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func encodeBatchLonLat(w http.ResponseWriter, r *http.Request) {

	var input []encodeBatchInput

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{
			string(reqBody),
			"Could not parse request body as valid JSON",
		})
		return
	}

	err = json.Unmarshal(reqBody, &input)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{
			string(reqBody),
			"Could not unmarshal JSON into a valid format",
		})
		return
	}

	fmt.Println(input)

	var res encodeBatchResponse

	res.Query = string(reqBody)

	for _, p := range input {
		hash, err := EncodeUsingPrecision(Point{p.Lat, p.Lon}, p.Pre)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{
				string(reqBody),
				"The request returned an error : " + err.Error(),
			})
			return
		}
		//fmt.Println(hash)
		query := fmt.Sprintf("%f,%f,%d", p.Lat, p.Lon, p.Pre)
		//fmt.Println(query)

		res.Results = append(res.Results, encodeResponse{
			query,
			p.Lat,
			p.Lon,
			p.Pre,
			hash,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func decodeGeohash(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]

	box, err := DecodeToBox(hash)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{
			hash,
			"The request returned an error : " + err.Error(),
		})
		return
	}
	lat := (box.LatMin + box.LatMax) / 2.
	lon := (box.LonMin + box.LonMax) / 2.

	res := &decodeResponse{hash, lat, lon, box}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func decodeBatchGeohash(w http.ResponseWriter, r *http.Request) {
	var input []decodeBatchInput

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{
			string(reqBody),
			"Could not parse request body as valid JSON",
		})
		return
	}

	err = json.Unmarshal(reqBody, &input)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{
			string(reqBody),
			"Could not unmarshal JSON into a valid format",
		})
		return
	}

	fmt.Println(input)

	var res decodeBatchResponse

	res.Query = string(reqBody)

	for _, p := range input {

		box, err := DecodeToBox(p.Geohash)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{
				res.Query,
				"The request returned an error : " + err.Error(),
			})
			return
		}
		lat := (box.LatMin + box.LatMax) / 2.
		lon := (box.LonMin + box.LonMax) / 2.

		res.Results = append(res.Results, decodeResponse{
			p.Geohash,
			lat,
			lon,
			box,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func route404Handler(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	json.NewEncoder(w).Encode(errorResponse{
		string(reqBody),
		"Route not found",
	})
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{"alive": "true"})
	//fmt.Fprintf(w, `{"alive": true}`)
}
