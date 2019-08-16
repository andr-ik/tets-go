package service

import (
	"github.com/tets-go/entity"
)

const InsertComplex = `
INSERT INTO complex(title, url, external_id)
	VALUES($1, $2, $3)
ON CONFLICT (external_id) 
DO UPDATE
	SET title = EXCLUDED.title,
		url = EXCLUDED.url
RETURNING id
`

type ComplexService struct {
	postgresService PostgresService
}

func NewComplexService(postgresService PostgresService) ComplexService {
	return ComplexService{
		postgresService: postgresService,
	}
}

func (complexService *ComplexService) Add(complex *entity.Complex) {
	var id int
	err := complexService.postgresService.DB.QueryRow(InsertComplex, complex.Title, complex.Code, complex.ExternalId).Scan(&id)
	if err != nil {
		panic(err)
	}

	complex.Id = int(id)
}

func (complexService *ComplexService) Load() {
	complexes := map[int]string{
		812: "pulse-na-naberezhnoi",
		219: "chistoe-nebo",
		795: "strizhi",
		797: "artline-v-primorskom",
	}

	for externalId, code := range complexes {
		complexEntity := entity.Complex{
			Title:      code,
			Code:       code,
			ExternalId: externalId,
		}

		complexService.Add(&complexEntity)
	}
}
