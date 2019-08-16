package main

import (
	"flag"
	"fmt"
	"github.com/tets-go/cookie"
	"github.com/tets-go/entity"
	"github.com/tets-go/parser"
	"github.com/tets-go/repository"
	"github.com/tets-go/service"
)

func main() {
	reload := flag.Bool("reload", false, "reload page")
	flag.Parse()

	guzzle := parser.NewGuzzle(parser.START_URL)

	//cookie.Save(guzzle)
	cookie.Load(guzzle)

	postgresService := service.NewPostgresService("127.0.0.1", 5440, "root", "", "postgres")
	defer postgresService.Close()

	complexes := parser.ParseComplexes(guzzle, *reload)

	complexService := service.NewComplexService(postgresService)
	for _, complexEntity := range complexes {
		complexService.Add(&complexEntity)
	}
	complexService.Load()

	complexRepository := repository.NewComplexRepository(postgresService)
	complexForParsing := []entity.Complex{
		complexRepository.Get(36),
		complexRepository.Get(37),
		complexRepository.Get(34),
	}

	flatService := service.NewFlatService(postgresService)
	historyService := service.NewHistoryService(postgresService)

	for _, complexEntity := range complexForParsing {
		flats := parser.ParseComplex(guzzle, complexEntity, *reload)

		for _, flat := range flats {
			flatService.Add(&flat)
			historyService.Add(&flat)
		}

		for _, flat := range flats {
			fmt.Println(flat)
		}
	}
}
