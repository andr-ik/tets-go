package cookie

import (
	"github.com/tets-go/parser"
	"net/url"
)

func domains() []string {
	domains := []string{
		"https://www.spbrealty.ru",
		"https://lk.spbrealty.ru",
	}

	return domains
}

func Save(guzzle parser.Guzzle) {
	parser.Login(guzzle)

	for _, domain := range domains() {
		host, _ := url.Parse(domain)
		_ = guzzle.SaveCookies(host)
	}
}

func Load(guzzle parser.Guzzle) {
	for _, domain := range domains() {
		host, _ := url.Parse(domain)
		_ = guzzle.LoadCookies(host)
	}
}
