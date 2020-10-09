package repository

type Repository interface {
	FindAll() ([]interface{}, error)
	Create(interface{}) error
	Update(interface{}) error
	Delete(interface{}) error
}
