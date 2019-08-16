package old

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

const START_URL = "https://lk.uk-dial.ru"
const LOGIN_URL = "https://lk.uk-dial.ru/login"
const GET_RESULTS_URL = "https://lk.uk-dial.ru/accruals/accruals?month=%s&year=%s"

const METHOD_GET = "GET"
const METHOD_POST = "POST"

const USER_AGENT = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36"

type Guzzle struct {
	client *http.Client
}

func NewGuzzle() Guzzle {
	var guzzle Guzzle

	cookieJar, _ := cookiejar.New(nil)
	guzzle.client = &http.Client{
		Jar: cookieJar,
	}

	return guzzle
}

func (guzzle *Guzzle) loadPage(url, method string, data url.Values) *goquery.Document {
	var response *http.Response
	request, _ := http.NewRequest(method, url, strings.NewReader(data.Encode()))

	request.Header.Set("User-Agent", USER_AGENT)
	request.Header.Set("referer", START_URL)

	if method == METHOD_POST {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	response, _ = guzzle.client.Do(request)
	defer response.Body.Close()

	document, _ := goquery.NewDocumentFromReader(response.Body)
	a, _ := document.Html()
	_ = a

	return document
}

type Row struct {
	Name  string
	Total float64
}

func getCSRF(guzzle Guzzle) string {
	document := guzzle.loadPage(START_URL, METHOD_GET, nil)

	if csrf, ok := document.Find("form input").FilterFunction(func(i int, selection *goquery.Selection) bool {
		if name, ok := selection.Attr("name"); ok {
			return name == "YII_CSRF_TOKEN"
		}

		return false
	}).Attr("value"); ok {
		return csrf
	}

	return ""
}

func login(guzzle Guzzle) {
	csrf := getCSRF(guzzle)
	guzzle.loadPage(LOGIN_URL, METHOD_POST, url.Values{
		"YII_CSRF_TOKEN":      {csrf},
		"LoginForm[email]":    {"mikhailov.andery@yandex.ru"},
		"LoginForm[password]": {"5540c338b8"},
	})
}

func loadData(guzzle Guzzle) map[string]Row {
	document := guzzle.loadPage(prepareLoadDataUri("2", "2019"), METHOD_GET, nil)

	var headers []string
	document.Find("table thead tr th").Each(func(i int, header *goquery.Selection) {
		headers = append(headers, header.Text())
	})

	values := make(map[string]Row)
	document.Find("table tbody tr").Each(func(i int, row *goquery.Selection) {
		tds := row.Find("td")
		name := tds.FilterFunction(func(i int, selection *goquery.Selection) bool {
			return headers[i] == "Услуга ЖКХ"
		}).Text()
		total := tds.FilterFunction(func(i int, selection *goquery.Selection) bool {
			return headers[i] == "Итого"
		}).Text()
		total = strings.Replace(total, "р.", "", -1)
		total = strings.Replace(total, "-", "0", -1)
		total = strings.Replace(total, ",", ".", -1)
		total = strings.Replace(total, " ", "", -1)
		totalValue, _ := strconv.ParseFloat(total, 64)

		values[name] = Row{
			Name:  name,
			Total: totalValue,
		}
	})

	return values
}

func prepareLoadDataUri(month, year string) string {
	return fmt.Sprintf(GET_RESULTS_URL, month, year)
}

func main() {
	guzzle := NewGuzzle()
	login(guzzle)
	headers := loadData(guzzle)

	fmt.Println(headers)
}
