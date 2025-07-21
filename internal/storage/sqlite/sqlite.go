package sqlite

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

type Costume struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Title       string    `gorm:"index"` // Индекс на поле Title
	Description string
	Image       string
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	// Открытие соединения с базой данных SQLite
	db, err := gorm.Open(sqlite.Open(storagePath), &gorm.Config{
		PrepareStmt:                              true, // Включение подготовленных запросов
		DisableForeignKeyConstraintWhenMigrating: true, // Отключение проверки внешних ключей при миграции
	})
	if err != nil {
		return nil, fmt.Errorf("%s: failed to open database: %w", op, err)
	}

	// Автоматическая миграция схемы (создание таблицы и индексов)
	if err := db.AutoMigrate(&Costume{}); err != nil {
		return nil, fmt.Errorf("%s: failed to migrate database schema: %w", op, err)
	}

	return &Storage{db: db}, nil
}
