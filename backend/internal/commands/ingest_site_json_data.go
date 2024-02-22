package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
	"usgs_tracker/internal/trends"

	"github.com/jackc/pgx/v4/pgxpool"
)

type IngestSiteJSONDataCommand struct {
	db        *pgxpool.Pool
	filenames []string
}

type SiteData struct {
	Value interface{}
}

type SiteRow struct {
	ID int
}

func (cmd *IngestSiteJSONDataCommand) process(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("no file found")
		return
	}
	defer file.Close()

	var siteData SiteData
	if err := json.NewDecoder(file).Decode(&siteData); err != nil {
		fmt.Println("malformed json")
		return
	}

	timeSeries := siteData.Value.(map[string]interface{})["timeSeries"].([]interface{})
	if len(timeSeries) == 0 {
		fmt.Println("could not parse", filename)
		return
	}

	siteName := timeSeries[0].(map[string]interface{})["sourceInfo"].(map[string]interface{})["siteName"].(string)
	siteCode := timeSeries[0].(map[string]interface{})["sourceInfo"].(map[string]interface{})["siteCode"].([]interface{})[0].(map[string]interface{})["value"]
	latitude := timeSeries[0].(map[string]interface{})["sourceInfo"].(map[string]interface{})["geoLocation"].(map[string]interface{})["geogLocation"].(map[string]interface{})["latitude"]
	longitude := timeSeries[0].(map[string]interface{})["sourceInfo"].(map[string]interface{})["geoLocation"].(map[string]interface{})["geogLocation"].(map[string]interface{})["longitude"]
	values := timeSeries[0].(map[string]interface{})["values"].([]interface{})[0].(map[string]interface{})["value"].([]interface{})

	dayAverages := make(map[string][]float64)
	levelMeasurements := []float64{}
	for _, value := range values {
		dateTime := value.(map[string]interface{})["dateTime"].(string)
		parsedTime, err := time.Parse(time.RFC3339, dateTime)
		if err != nil {
			fmt.Println("error parsing date", err)
		}
		dateKey := parsedTime.Format("2006-01-02")

		measurement := value.(map[string]interface{})["value"].(string)
		v, _ := strconv.ParseFloat(measurement, 64)
		v /= 3.28084

		dayAverages[dateKey] = append(dayAverages[dateKey], v)
		levelMeasurements = append(levelMeasurements, v)
	}

	dailyAverageMeasurements := make(map[string]float64)
	for date, values := range dayAverages {
		sum := 0.0
		for _, v := range values {
			sum += v
		}
		dailyAverageMeasurements[date] = sum / float64(len(values))
	}

	sortedDates := make([]string, 0, len(dailyAverageMeasurements))
	for date := range dailyAverageMeasurements {
		sortedDates = append(sortedDates, date)
	}
	sort.Strings(sortedDates)

	xs := make([]float64, 0, len(dailyAverageMeasurements))
	averagesSlice := make([]float64, 0, len(dailyAverageMeasurements))
	for idx, date := range sortedDates {
		xs = append(xs, float64(idx))
		averagesSlice = append(averagesSlice, dailyAverageMeasurements[date])
	}

	medianSlope := trends.ApproxMedianSlopeInMM(xs, averagesSlice, 7)
	zValue, _ := trends.MannKendall(averagesSlice)

	statement := fmt.Sprintf("INSERT INTO sites (site_name, site_no, location, ts_estimate_in_mm, mk_z) VALUES ('%s', '%s', ST_GeomFromText('POINT(%f %f)', 4326), %f, %f);", siteName, siteCode, longitude, latitude, medianSlope, zValue)
	_, errr := cmd.db.Exec(context.Background(), statement)
	if errr != nil {
		fmt.Println(err)
	}

	row := cmd.db.QueryRow(context.Background(), fmt.Sprintf("SELECT id FROM sites WHERE site_no = '%s'", siteCode))
	r := SiteRow{}
	row.Scan(&r.ID)

	for day, value := range dailyAverageMeasurements {
		statement = fmt.Sprintf("INSERT INTO datapoints (site_id, value, ts) VALUES (%d, %f, '%s');", r.ID, value, day)
		_, err := cmd.db.Exec(context.Background(), statement)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("finished", siteCode)
}

func (cmd IngestSiteJSONDataCommand) Run(args []string) {
	var wg sync.WaitGroup
	workerCount := 5
	fileChan := make(chan string, workerCount)

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filename := range fileChan {
				cmd.process(filename)
			}
		}()
	}

	for _, filename := range cmd.filenames {
		fileChan <- filename
	}
	close(fileChan)

	wg.Wait()
}
