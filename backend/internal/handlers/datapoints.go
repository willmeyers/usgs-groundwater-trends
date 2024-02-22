package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Datapoint struct {
	Year  int
	Value float32
}

type DatapointResponse struct {
	Datapoints []Datapoint
}

func (handler *ServerHandler) DatapointsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	siteID := r.URL.Query().Get("site_id")
	if siteID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "site_id is required"})
		return
	}

	var response DatapointResponse

	query := `SELECT EXTRACT(YEAR FROM ts) AS year, AVG(value) AS average_value FROM datapoints WHERE site_id = $1 AND value > 0 GROUP BY EXTRACT(YEAR FROM ts) ORDER BY year ASC;`
	rows, err := handler.DB.Query(context.Background(), query, siteID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var datapoint Datapoint
		if err := rows.Scan(&datapoint.Year, &datapoint.Value); err != nil {
			fmt.Println(err)
			continue
		}

		response.Datapoints = append(response.Datapoints, datapoint)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
