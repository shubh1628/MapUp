// main.go

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

type RequestPayload struct {
	ToSort [][]int `json:"to_sort"`
}

type ResponsePayload struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNs       int64   `json:"time_ns"`
}

func processSingle(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := json.NewDecoder(r.Body).Decode(&requestPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	sortedArrays := make([][]int, len(requestPayload.ToSort))
	for i, arr := range requestPayload.ToSort {
		sorted := make([]int, len(arr))
		copy(sorted, arr)
		sort.Ints(sorted)
		sortedArrays[i] = sorted
	}
	timeTaken := time.Since(startTime)

	responsePayload := ResponsePayload{
		SortedArrays: sortedArrays,
		TimeNs:       timeTaken.Nanoseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responsePayload)
}

func processConcurrent(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := json.NewDecoder(r.Body).Decode(&requestPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	var wg sync.WaitGroup
	sortedArrays := make([][]int, len(requestPayload.ToSort))
	for i, arr := range requestPayload.ToSort {
		wg.Add(1)
		go func(i int, arr []int) {
			defer wg.Done()
			sorted := make([]int, len(arr))
			copy(sorted, arr)
			sort.Ints(sorted)
			sortedArrays[i] = sorted
		}(i, arr)
	}
	wg.Wait()
	timeTaken := time.Since(startTime)

	responsePayload := ResponsePayload{
		SortedArrays: sortedArrays,
		TimeNs:       timeTaken.Nanoseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responsePayload)
}

func main() {
	http.HandleFunc("/process-single", processSingle)
	http.HandleFunc("/process-concurrent", processConcurrent)

	fmt.Println("Server listening on :8000")
	http.ListenAndServe(":8000", nil)
}
