package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
)

type Guzzle struct {
	client  *http.Client
	referer string
}

func NewGuzzle(startReferer string) Guzzle {
	var guzzle Guzzle

	cookieJar, _ := cookiejar.New(nil)
	guzzle.client = &http.Client{
		Jar: cookieJar,
	}
	guzzle.referer = startReferer

	return guzzle
}

func (guzzle *Guzzle) LoadPage(url, method string, data url.Values) *goquery.Document {
	var response *http.Response
	request, _ := http.NewRequest(method, url, strings.NewReader(data.Encode()))

	fmt.Printf("%s: %s", method, url)
	fmt.Println()

	request.Header.Set("User-Agent", USER_AGENT)
	request.Header.Set("referer", guzzle.referer)

	guzzle.referer = url

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

func (guzzle *Guzzle) SaveCookies(domain *url.URL) error {
	f, err := os.Create(fmt.Sprintf("cookies-%s.json", domain.Host))
	if err != nil {
		return err
	}
	defer f.Close()

	cookies := guzzle.client.Jar.Cookies(domain)
	r, err := json.Marshal(cookies)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, bytes.NewReader(r))
	return err
}

func (guzzle *Guzzle) LoadCookies(domain *url.URL) error {
	f, err := os.Open(fmt.Sprintf("cookies-%s.json", domain.Host))
	if err != nil {
		return err
	}
	defer f.Close()

	var cookies []*http.Cookie

	jsonRaw, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonRaw, &cookies)
	if err != nil {
		return err
	}

	guzzle.client.Jar.SetCookies(domain, cookies)

	return nil
}

func (guzzle *Guzzle) GetJar() http.CookieJar {
	return guzzle.client.Jar
}
