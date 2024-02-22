package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Site struct {
	ID             int
	SiteName       string
	SiteNo         string
	Location       interface{}
	TSEstimateInMM float64
	MKZ            float64
}

type SitesResponse struct {
	Sites []Site
}

func (handler *ServerHandler) SitesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var response SitesResponse

	rows, err := handler.DB.Query(context.Background(), "SELECT id, site_name, site_no, ST_AsGeoJSON(location) as location, ts_estimate_in_mm, mk_z FROM sites;")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var site Site
		if err := rows.Scan(&site.ID, &site.SiteName, &site.SiteNo, &site.Location, &site.TSEstimateInMM, &site.MKZ); err != nil {
			fmt.Println(err)
			continue
		}

		response.Sites = append(response.Sites, site)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
