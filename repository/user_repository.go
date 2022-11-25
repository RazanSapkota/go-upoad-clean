package repository

import (
	"example/go-api/infrastructure"
)

// UserRepository database structure
type UserRepository struct {
	infrastructure.Database

}

// NewUserRepository creates a new user repository
func NewUserRepository(db infrastructure.Database) UserRepository {
	
	return UserRepository{
		Database: db,
	}
}