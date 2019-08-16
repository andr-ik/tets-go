package parser

import "net/url"

const START_URL = "https://www.spbrealty.ru/"
const LOGIN_URL = "https://lk.spbrealty.ru/auth?client_id=spbrealty&response_type=code&redirect_uri=https://www.spbrealty.ru/auth/cabinet"

const METHOD_GET = "GET"
const METHOD_POST = "POST"

const USER_AGENT = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36"

func Login(guzzle Guzzle) {
	document := guzzle.LoadPage(LOGIN_URL, METHOD_POST, url.Values{
		"type":     {"phone"},
		"phone":    {"+7(911) 842-0815"},
		"email":    {"dz"},
		"password": {"dzMV47!ghm"},
	})

	a, _ := document.Html()

	_ = a
}

func LoadBoobaleh(guzzle Guzzle) {
	document := guzzle.LoadPage("https://boobaleh.ru", METHOD_GET, url.Values{})
	a, _ := document.Html()
	_ = a
}
