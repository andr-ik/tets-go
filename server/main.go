package main

import (
	"encoding/json"
	"fmt"
	"github.com/tets-go/repository"
	"github.com/tets-go/service"
	"net/http"
)

func main() {
	postgresService := service.NewPostgresService("127.0.0.1", 5440, "root", "", "postgres")
	defer postgresService.Close()

	//start := time.Date(2019, 7, 11, 0, 0, 0, 0, time.UTC)
	//end := time.Now()
	//interval := service.NewDateIntervalService(start, 24 * time.Hour, end)

	http.HandleFunc("/api/price", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		complexId, ok := query["complexId"]
		if !ok || len(complexId[0]) < 1 {
			return
		}

		id := complexId[0]

		historyRepository := repository.NewHistoryRepository(postgresService)
		histories := historyRepository.GetAllForComplexId(id)

		historiesMap := make(map[int]map[string]float64)
		for _, history := range histories {
			if _, ok := historiesMap[history.FlatId]; !ok {
				historiesMap[history.FlatId] = make(map[string]float64)
			}

			if _, ok := historiesMap[history.FlatId]["id"]; !ok {
				historiesMap[history.FlatId]["id"] = float64(history.FlatId)
			}

			historiesMap[history.FlatId][history.Date.Format("02-01")] = history.Price
		}

		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		jsonRaw, _ := json.Marshal(historiesMap)
		_, _ = fmt.Fprintf(w, string(jsonRaw))
	})
	fmt.Println("Run server: http://127.0.0.1:1081/api/price")
	http.ListenAndServe(":1081", nil)
}
