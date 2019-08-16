package service

import (
	"github.com/tets-go/entity"
)

const InsertFlat = `
INSERT INTO flat(complex_id, area, external_id, floor, deadline, type, image)
	VALUES($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (external_id) 
DO UPDATE
	SET complex_id = EXCLUDED.complex_id,
		area = EXCLUDED.area,
		floor = EXCLUDED.floor,
		deadline = EXCLUDED.deadline,
		type = EXCLUDED.type,
		image = EXCLUDED.image
RETURNING id
`

type FlatService struct {
	postgresService PostgresService
}

func NewFlatService(postgresService PostgresService) FlatService {
	return FlatService{
		postgresService: postgresService,
	}
}

func (flatService *FlatService) Add(flat *entity.Flat) {
	var id int
	err := flatService.postgresService.DB.QueryRow(InsertFlat, flat.Complex.Id, flat.Area, flat.ExternalId, flat.Floor, flat.Deadline, flat.Type, flat.Image).Scan(&id)
	if err != nil {
		panic(err)
	}

	flat.Id = int(id)
}
