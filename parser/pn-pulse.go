package parser

import (
	"crypto/md5"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tets-go/entity"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const ComplexesUrl = "https://www.spbrealty.ru/buildings/ajax/selectLivecomplexNew?city=spb"
const ComplexUrl = "https://www.spbrealty.ru/buildings/%s"
const FlatsUrl = "https://www.spbrealty.ru/buildings/ajax/showFlatCardLayout?id=%s"

func SaveComplexesPage(guzzle Guzzle) {
	SavePage(guzzle, prepareComplexesUrl(), fmt.Sprintf("tmp/complexes-%x.html", md5.Sum([]byte(prepareComplexesUrl()))))
}

func SaveComplexPage(guzzle Guzzle, complex entity.Complex) {
	SavePage(guzzle, prepareComplexUrl(complex), fmt.Sprintf("tmp/%s-%x.html", complex.Code, md5.Sum([]byte(prepareComplexUrl(complex)))))
}

func SaveFlatPage(guzzle Guzzle, complex entity.Complex, index string) {
	SavePage(guzzle, prepareFlatUrl(index), fmt.Sprintf("tmp/%s-flat-%x.html", complex.Code, md5.Sum([]byte(prepareFlatUrl(index)))))
}

func SavePage(guzzle Guzzle, link, file string) {
	document := guzzle.LoadPage(link, METHOD_GET, url.Values{})

	f, _ := os.Create(file)
	defer f.Close()

	data, _ := document.Html()
	_, err := io.Copy(f, strings.NewReader(data))
	if err != nil {
		panic(err)
	}
}

func LoadComplexesPage() *goquery.Document {
	return LoadPage(fmt.Sprintf("tmp/complexes-%x.html", md5.Sum([]byte(prepareComplexesUrl()))))
}

func LoadComplexPage(complex entity.Complex) *goquery.Document {
	return LoadPage(fmt.Sprintf("tmp/%s-%x.html", complex.Code, md5.Sum([]byte(prepareComplexUrl(complex)))))
}

func LoadFlatPage(complex entity.Complex, index string) *goquery.Document {
	return LoadPage(fmt.Sprintf("tmp/%s-flat-%x.html", complex.Code, md5.Sum([]byte(prepareFlatUrl(index)))))
}

func LoadPage(file string) *goquery.Document {
	f, _ := os.Open(file)
	defer f.Close()

	document, _ := goquery.NewDocumentFromReader(f)

	return document
}

func IsEmptyPage(document *goquery.Document) bool {
	html, _ := document.Html()

	return strings.EqualFold(html, "<html><head></head><body></body></html>")
}

func ParseComplexes(guzzle Guzzle, reloadCache bool) []entity.Complex {
	var complexes []entity.Complex

	if reloadCache {
		SaveComplexesPage(guzzle)
	}

	document := LoadComplexesPage()
	document.Find(".items li").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		var complexTitle string
		var complexCode string
		var complexExternalId int

		complexTitle = parseComplexTitle(selection)
		complexCode = parseComplexCode(selection)
		complexExternalId = parseComplexExternalId(selection)

		complexEntity := entity.Complex{
			Title:      complexTitle,
			Code:       complexCode,
			ExternalId: complexExternalId,
		}

		complexes = append(complexes, complexEntity)

		return true
	})

	return complexes
}

func ParseComplex(guzzle Guzzle, complex entity.Complex, reloadCache bool) []entity.Flat {
	var flats []entity.Flat

	if reloadCache {
		SaveComplexPage(guzzle, complex)
	}

	document := LoadComplexPage(complex)
	document.Find("#flats-table .sticky-wrapper").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		flatType := parseFlatType(selection)

		var flatFloor int
		var flatExternalId int
		var flatImage string
		var flatDeadline string
		var flatArea float64
		var flatPrice float64

		selection.Find("table tr").EachWithBreak(func(i int, selection *goquery.Selection) bool {
			index := parseFlatIndex(selection)
			if reloadCache {
				SaveFlatPage(guzzle, complex, index)
			}

			flatDeadline = parseFlatDeadline(selection)
			flatArea = parseFlatArea(selection)

			documentSub := LoadFlatPage(complex, index)
			if IsEmptyPage(documentSub) {
				fmt.Println("FAIL LOAD URL")
				return true
			}

			items := documentSub.Find(".plan .overlay-body ul li")

			flatImage = parseFlatImage(documentSub.Find(".plan"))

			if items.Length() > 0 {
				items.Each(func(i int, selection *goquery.Selection) {
					flatFloor = parseFlatFloorFromList(selection)
					flatExternalId = parseFlatExternalIdFromList(selection)
					flatPrice = parseFlatPriceFromList(selection)

					flat := entity.Flat{
						Type:       flatType,
						Complex:    complex,
						Floor:      flatFloor,
						Image:      flatImage,
						ExternalId: flatExternalId,
						Deadline:   flatDeadline,
						Area:       flatArea,
						Price:      flatPrice,
					}

					flats = append(flats, flat)
				})
			} else {
				/*data, _ := documentSub.Find(".info .title span").Html()
				data2, _ := documentSub.Find(".info .title").Find("span").Html()
				data3 := documentSub.Find(".info .title").Find("span").Text()
				_ = data
				_ = data2
				_ = data3*/
				flatFloor = parseFlatFloorFromCard(documentSub.Find(".info .title"))
				flatExternalId = parseFlatExternalIdFromCard(documentSub.Find(".info .title"))
				flatPrice = parseFlatPriceFromCard(documentSub.Find(".info .title"))

				flat := entity.Flat{
					Type:       flatType,
					Complex:    complex,
					Floor:      flatFloor,
					Image:      flatImage,
					ExternalId: flatExternalId,
					Deadline:   flatDeadline,
					Area:       flatArea,
					Price:      flatPrice,
				}

				flats = append(flats, flat)
			}

			return true
		})

		return true
	})

	return flats
}

func parseComplexTitle(selection *goquery.Selection) string {
	result, _ := selection.Attr("title")

	return strings.TrimSpace(result)
}

func parseComplexCode(selection *goquery.Selection) string {
	result, _ := selection.Attr("title")
	//return strings.TrimSpace("")

	return strings.TrimSpace(result)
}

func parseComplexExternalId(selection *goquery.Selection) int {
	resultString, _ := selection.Attr("data-bid")
	result, _ := strconv.Atoi(strings.TrimSpace(resultString))

	return result
}

func parseFlatIndex(selection *goquery.Selection) string {
	result, _ := selection.Attr("data-index")

	return result
}

func parseFlatType(selection *goquery.Selection) string {
	return strings.TrimSpace(selection.Find(".divider .name").Text())
}

func parseFlatDeadline(selection *goquery.Selection) string {
	return strings.TrimSpace(selection.Find(".deadline").Text())
}

func parseFlatImage(selection *goquery.Selection) string {
	image, _ := selection.Find("img").Attr("src")

	return strings.TrimSpace(image)
}

func parseFlatArea(selection *goquery.Selection) float64 {
	infoHtml, _ := selection.Find(".info").Html()
	infoArray := strings.Split(infoHtml, "<br/>")
	areaInfo := strings.Split(infoArray[0], " ")
	area, _ := strconv.ParseFloat(strings.TrimSpace(areaInfo[0]), 64)

	return area
}

func parseFlatFloorFromList(selection *goquery.Selection) int {
	result, _ := strconv.Atoi(strings.TrimSpace(selection.Find(".floor").Text()))

	return result
}

func parseFlatFloorFromCard(selection *goquery.Selection) int {
	return 0
}

func parseFlatExternalIdFromList(selection *goquery.Selection) int {
	js, _ := selection.Find("a").Attr("href")
	id := strings.ReplaceAll(strings.ReplaceAll(js, "javascript:showFlatLayout(", ""), ",'card');", "")
	result, _ := strconv.Atoi(strings.TrimSpace(id))

	return result
}

func parseFlatExternalIdFromCard(selection *goquery.Selection) int {
	id, _ := selection.Find("a").Attr("rel")
	result, _ := strconv.Atoi(strings.TrimSpace(id))

	return result
}

func parseFlatPriceFromList(selection *goquery.Selection) float64 {
	price := selection.Find(".price").Text()
	price = strings.ReplaceAll(strings.ReplaceAll(price, "руб.", ""), " ", "")
	result, _ := strconv.ParseFloat(strings.TrimSpace(price), 64)

	return result
}

func parseFlatPriceFromCard(selection *goquery.Selection) float64 {
	price, _ := selection.Find("span").Html()
	price = strings.ReplaceAll(strings.ReplaceAll(price, "₽", ""), " ", "")
	result, _ := strconv.ParseFloat(strings.TrimSpace(price), 64)

	return result
}

func prepareComplexesUrl() string {
	return ComplexesUrl
}

func prepareComplexUrl(complex entity.Complex) string {
	return fmt.Sprintf(ComplexUrl, complex.Code)
}

func prepareFlatUrl(index string) string {
	return fmt.Sprintf(FlatsUrl, index)
}
