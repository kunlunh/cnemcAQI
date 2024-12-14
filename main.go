/*
Created on Dec 14 17:44:00 2024
@author: Kunlun HUANG
github.com/kunlunh
*/
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type AQIData struct {
	PositionName string `json:"PositionName"`
	Quality      string `json:"Quality"`
	AQI          string `json:"AQI"`
	O3           string `json:"O3"`
	NO2          string `json:"NO2"`
	PM10         string `json:"PM10"`
	PM25         string `json:"PM2_5"`
	SO2          string `json:"SO2"`
	CO           string `json:"CO"`
	Latitude     string `json:"Latitude"`
	Longitude    string `json:"Longitude"`
}

const url = "https://air.cnemc.cn:18007/CityData/GetAQIDataPublishLive?cityName=%E5%B9%BF%E5%B7%9E%E5%B8%82"

func fetchAQIData() ([]AQIData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data []AQIData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data, err := fetchAQIData()
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		log.Println("Error fetching data:", err)
		return
	}

	tmpl := `<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>Air Quality Data</title>
</head>
<body>
	<h1>Guangzhou Air Quality Data</h1>
	<table border="1" style="border-collapse: collapse; width: 100%; text-align: center;">
		<tr>
			<th>Position Name</th>
			<th>Quality</th>
			<th>AQI</th>
			<th>O3</th>
			<th>NO2</th>
			<th>PM10</th>
			<th>PM2.5</th>
			<th>SO2</th>
			<th>CO</th>
			<th>Latitude</th>
			<th>Longitude</th>
		</tr>
		{{range .}}
		<tr>
			<td>{{.PositionName}}</td>
			<td>{{.Quality}}</td>
			<td>{{.AQI}}</td>
			<td>{{.O3}}</td>
			<td>{{.NO2}}</td>
			<td>{{.PM10}}</td>
			<td>{{.PM25}}</td>
			<td>{{.SO2}}</td>
			<td>{{.CO}}</td>
			<td>{{.Latitude}}</td>
			<td>{{.Longitude}}</td>
		</tr>
		{{end}}
	</table>
</body>
</html>`

	parsedTemplate, err := template.New("webpage").Parse(tmpl)
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		log.Println("Error parsing template:", err)
		return
	}

	if err := parsedTemplate.Execute(w, data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Println("Error executing template:", err)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	fmt.Println("Starting server on 127.0.0.1:50001...")
	if err := http.ListenAndServe("127.0.0.1:50001", nil); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
