package mapper

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/dom/vos"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/dom/dto"
)

func ProductToDomain(dto dto.ProductDTO) (vos.Product, error) {
	productId := vos.ProductID(dto.ID)
	return vos.NewProduct(
		productId,
		dto.Name,
		dto.Description,
		dto.Team,
		dto.Organization,
	)
}

func ProductToDTO(product vos.Product) dto.ProductDTO {
	return dto.ProductDTO{
		ID:           string(product.ID()),
		Name:         product.Name(),
		Description:  product.Description(),
		Team:         product.Team(),
		Organization: product.Organization(),
	}
}
