package service

import (
	"github.com/tets-go/entity"
	"time"
)

const InsertHistory = `
INSERT INTO history(flat_id, date, price)
	VALUES($1, $2, $3)
ON CONFLICT (flat_id, date) 
DO UPDATE
	SET price = EXCLUDED.price
RETURNING id
`

type HistoryService struct {
	postgresService PostgresService
}

func NewHistoryService(postgresService PostgresService) HistoryService {
	return HistoryService{
		postgresService: postgresService,
	}
}

func (historyService *HistoryService) Add(flat *entity.Flat) {
	date := time.Now() //time.Parse("2006-01-02", "2019-07-11")
	_, err := historyService.postgresService.DB.Exec(InsertHistory, flat.Id, date, flat.Price)
	if err != nil {
		panic(err)
	}
}
