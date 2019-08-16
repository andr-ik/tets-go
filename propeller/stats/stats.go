package stats

import (
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Metrics struct {
	Requests      string
	Impressions   string
	Clicks        string
	Conversions   string
	Subscriptions string
	Revenue       string
	Costs         string
}

func (m *Metrics) getValue(field string) string {
	requests, _ := strconv.Atoi(m.Requests)
	impressions, _ := strconv.Atoi(m.Impressions)
	clicks, _ := strconv.Atoi(m.Clicks)
	conversions, _ := strconv.Atoi(m.Conversions)
	subscriptions, _ := strconv.Atoi(m.Subscriptions)
	revenue, _ := strconv.ParseFloat(m.Revenue, 64)
	cost, _ := strconv.ParseFloat(m.Costs, 64)

	array := map[string]string{
		"Requests":      fmt.Sprintf("%d", requests),
		"Impressions":   fmt.Sprintf("%d", impressions),
		"Clicks":        fmt.Sprintf("%d", clicks),
		"Conversions":   fmt.Sprintf("%d", conversions),
		"Subscriptions": fmt.Sprintf("%d", subscriptions),
		"Revenue":       fmt.Sprintf("%.2f", revenue),
		"Costs":         fmt.Sprintf("%.2f", cost),
		"Profit":        fmt.Sprintf("%.2f", revenue-cost),
	}

	return array[field]
}

func Start() {
	periods := map[string][]string{
		"2017":       []string{"2017-01-01", "2017-12-31"},
		"2018":       []string{"2018-01-01", "2018-12-31"},
		"Q1+Q2 2019": []string{"2019-01-01", "2019-06-30"},
	}
	headers := []string{"Direction", "RateModel"}
	metrics := []string{"Requests", "Impressions", "Clicks", "Conversions", "Subscriptions", "Revenue", "Costs", "Profit"}

	report := make(map[string]map[string]map[string]Metrics)
	for namePeriod, period := range periods {
		params := map[string]string{
			"day_from": period[0],
			"day_to":   period[1],
			"group":    "direction,revenue_type",
			"metrics":  strings.ToLower(strings.Join(metrics, ",")),
		}

		report[namePeriod] = loadStats(params)
	}

	saveReport(headers, metrics, periods, report)
}

func saveReport(headers, metrics []string, periods map[string][]string, report map[string]map[string]map[string]Metrics) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Report")
	if err != nil {
		fmt.Printf(err.Error())
	}
	row = sheet.AddRow()

	for _, nameHeader := range headers {
		cell = row.AddCell()
		cell.Value = nameHeader
	}

	for _, nameMetric := range metrics {
		for _, indexPeriod := range sortMapString(periods) {
			headers = append(headers, fmt.Sprintf("%s %s", nameMetric, indexPeriod))

			cell = row.AddCell()
			cell.Value = fmt.Sprintf("%s %s", nameMetric, indexPeriod)
		}
	}

	directions := getDirections()
	priceModels := getPriceModel()
	for _, indexDirection := range sortMap(directions) {
		for _, indexPriceModel := range sortMap(priceModels) {
			nameDirection := directions[indexDirection]
			namePriceModel := priceModels[indexPriceModel]

			row = sheet.AddRow()
			cell = row.AddCell()
			cell.Value = nameDirection
			cell = row.AddCell()
			cell.Value = namePriceModel

			for _, nameMetric := range metrics {
				for _, indexPeriod := range sortMapString(periods) {
					metrics := report[indexPeriod][nameDirection][namePriceModel]
					cell = row.AddCell()
					cell.Value = metrics.getValue(nameMetric)
				}
			}
		}
	}

	err = file.Save("Report.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func loadStats(params map[string]string) map[string]map[string]Metrics {
	baseUrl := "http://1stats-api.rtty.in/api/v1/stats/"
	baseParams := map[string]string{
		"_format": "",
		"token":   "eV257kaFXskvk6xFBMJ3ZQK97ZKkDJAA",
	}

	url := baseUrl
	resultParams := []string{}
	for key, value := range baseParams {
		resultParams = append(resultParams, fmt.Sprintf("%s=%s", key, value))
	}
	for key, value := range params {
		resultParams = append(resultParams, fmt.Sprintf("%s=%s", key, value))
	}
	url += "?" + strings.Join(resultParams, "&")

	fmt.Println(url)

	response, _ := http.Get(url)
	defer response.Body.Close()
	jsonRaw, _ := ioutil.ReadAll(response.Body)

	statsResponse := make(map[string]Metrics)
	_ = json.Unmarshal(jsonRaw, &statsResponse)

	report := make(map[string]map[string]Metrics)
	for slice, metrics := range statsResponse {
		direction, rateModel := parseSlice(slice)

		if direction == "" || rateModel == "" {
			continue
		}

		if _, ok := report[direction]; !ok {
			report[direction] = make(map[string]Metrics)
		}

		report[direction][rateModel] = metrics
	}

	return report
}

func parseSlice(slice string) (string, string) {
	sliceArray := strings.Split(slice, "|")

	directions := getDirections()
	priceModels := getPriceModel()

	return directions[sliceArray[0]], priceModels[sliceArray[1]]
}

func sortMap(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

func sortMapString(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

func getPriceModel() map[string]string {
	return map[string]string{
		"":  "null",
		"1": "cpm",
		"2": "cpc",
		"3": "cpa",
		"4": "parsed",
		"5": "rtb",
		"6": "scpm",
		"7": "scpa",
	}
}

func getDirections() map[string]string {
	return map[string]string{
		"":   "null",
		"1":  "onclick",
		"2":  "video",
		"3":  "banner",
		"4":  "notice",
		"5":  "direct",
		"6":  "pushup",
		"7":  "mobnotice",
		"8":  "mobbanner",
		"11": "adex",
		"12": "syndi",
		"13": "mediabuy",
		"14": "brand",
		"15": "interstitial",
		"46": "onclick5",
		"50": "other",
		"51": "inapp",
		"52": "discount",
		"53": "nativeads",
		"54": "rtb",
		"55": "native",
		"57": "pusherpps",
		"58": "pushservice",
	}
}
