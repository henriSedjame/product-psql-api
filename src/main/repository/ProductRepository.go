package repository

import (
	"github.com/go-pg/pg/v10"
	"github.com/hsedjame/product-psql-api/src/main/models"
)

type ProductRepository struct {
	DB *pg.DB
}

func NewProductRepository(DB *pg.DB) *ProductRepository {
	return &ProductRepository{DB: DB}
}

func (p *ProductRepository) FindAll() ([]models.Product, error) {
	var products []models.Product

	if err := p.DB.Model(&products).Select(); err != nil {
		return nil, err
	}

	return products, nil
}

func (p *ProductRepository) Create(product *models.Product) error {
	_, err := p.DB.Model(product).Insert()
	return err
}

func (p *ProductRepository) Update(product *models.Product) error {
	_, err := p.DB.Model(product).WherePK().Update()
	return err
}

func (p *ProductRepository) Delete(id int) error {
	product := &models.Product{}
	_, err := p.DB.Model(product).Where("id = ?", id).Delete()
	return err
}
