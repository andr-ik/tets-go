package repository

import (
	"github.com/tets-go/entity"
	"github.com/tets-go/service"
)

const SELECT_HISTORY = `
SELECT history.* FROM history
INNER JOIN flat ON flat.id = history.flat_id
INNER JOIN complex ON complex.id = flat.complex_id
WHERE complex.id = $1
`

type HistoryRepository struct {
	postgresService service.PostgresService
}

func NewHistoryRepository(postgresService service.PostgresService) HistoryRepository {
	return HistoryRepository{
		postgresService: postgresService,
	}
}

func (historyRepository *HistoryRepository) GetAllForComplexId(id string) []entity.History {
	rows, _ := historyRepository.postgresService.DB.Query(SELECT_HISTORY, id)
	defer rows.Close()

	var result []entity.History
	for rows.Next() {
		var entityHistory = entity.History{}
		if err := rows.Scan(&entityHistory.Id, &entityHistory.FlatId, &entityHistory.Date, &entityHistory.Price); err != nil {
			//
		}

		result = append(result, entityHistory)
	}

	return result
}
