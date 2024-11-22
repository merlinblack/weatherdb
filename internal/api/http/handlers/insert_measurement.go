package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/merlinblack/weatherdb/internal/config"
	"github.com/merlinblack/weatherdb/internal/repository/weather"
)

type POSTInsertData struct {
	Key         *string  `json:"key"`
	Temperature *float64 `json:"temperature"`
	Humidity    *float64 `json:"humidity"`
	Pressure    *float64 `json:"pressure"`
	RecordedAt  *string  `json:"recordedAt"`
	Location    *string  `json:"location"`

	ParsedRecordedAt time.Time
}

func InsertMeasurement(w http.ResponseWriter, r *http.Request, cfg *config.Config, repo *weather.Queries) {

	data, err := decodeData(r)

	if err != nil {
		internal400(w, `There was a problem parsing the POST data`, err)
		return
	}

	err = validateData(data)

	if err != nil {
		internal400(w, `Validation error`, err)
		return
	}

	if *data.Key != cfg.API.WritePassword {
		notAllowed(w)
		return
	}

	if data.Location == nil {
		data.Location = &cfg.API.DefaultLocation
	}

	log.Printf("Inserting data: %s, %f, %f, %f %s\n", *data.RecordedAt, *data.Temperature, *data.Humidity, *data.Pressure, *data.Location)

	record, err := repo.InsertMeasurement(context.Background(),
		weather.InsertMeasurementParams{
			RecordedAt:  data.ParsedRecordedAt,
			Temperature: *data.Temperature,
			Humidity:    *data.Humidity,
			Pressure:    *data.Pressure,
			Location:    *data.Location,
		},
	)

	if err != nil {
		internal500(w, `Problem inserting row into DB`, err)
		return
	}

	writeResponse(w, record)
}

func decodeData(r *http.Request) (*POSTInsertData, error) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	data := &POSTInsertData{}

	err := decoder.Decode(data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func validateData(data *POSTInsertData) error {

	if data.Key == nil {
		return errors.New(`missing requried field: key`)
	}

	if data.Temperature == nil {
		return errors.New(`missing requried field: temperature`)
	}

	if data.Humidity == nil {
		return errors.New(`missing requried field: humidity`)
	}

	if data.Pressure == nil {
		return errors.New(`missing requried field: pressure`)
	}

	if data.RecordedAt == nil {
		return errors.New(`missing requried field: recordedAt`)
	}

	time, err := time.Parse(timeJSONLayout, *data.RecordedAt)
	if err != nil {
		return err
	}

	data.ParsedRecordedAt = time

	return nil
}

func notAllowed(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	resp := make(map[string]string)

	resp["message"] = `Not allowed`
	resp["status"] = "400"

	jsonResp, err := json.Marshal(resp)

	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(jsonResp)

}

func writeResponse(w http.ResponseWriter, data weather.Measurement) {

	resp := make(map[string]interface{})

	resp[`recordedAt`] = data.RecordedAt.Format(timeJSONLayout)
	resp[`temperature`] = data.Temperature
	resp[`humidity`] = data.Humidity
	resp[`pressure`] = data.Pressure
	resp[`location`] = data.Location

	jsonResp, err := json.Marshal(resp)

	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonResp)

}