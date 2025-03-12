package sqlite

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

type Storage struct {
	db *gorm.DB
}

type Costume struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Title       string
	Description string
	Image       string
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
		PrepareStmt: true,
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

}
