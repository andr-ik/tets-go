package repository

import (
	"github.com/tets-go/entity"
	"github.com/tets-go/service"
)

const SELECT_COMPLEX_BY_ID = `
SELECT * FROM complex WHERE id = $1
`

const SELECT_COMPLEX = `
SELECT * FROM complex
`

type ComplexRepository struct {
	postgresService service.PostgresService
}

func NewComplexRepository(postgresService service.PostgresService) ComplexRepository {
	return ComplexRepository{
		postgresService: postgresService,
	}
}

func (complexRepository *ComplexRepository) Get(id int) entity.Complex {
	rows, _ := complexRepository.postgresService.DB.Query(SELECT_COMPLEX_BY_ID, id)
	defer rows.Close()

	var entityComplex = entity.Complex{}
	for rows.Next() {
		if err := rows.Scan(&entityComplex.Id, &entityComplex.Title, &entityComplex.Code, &entityComplex.ExternalId); err != nil {
			//
		}
	}

	return entityComplex
}

func (complexRepository *ComplexRepository) GetAll() []entity.Complex {
	rows, _ := complexRepository.postgresService.DB.Query(SELECT_COMPLEX)
	defer rows.Close()

	var result []entity.Complex
	for rows.Next() {
		var entityComplex = entity.Complex{}
		if err := rows.Scan(&entityComplex.Id, &entityComplex.Title, &entityComplex.Code, &entityComplex.ExternalId); err != nil {
			//
		}

		result = append(result, entityComplex)
	}

	return result
}
