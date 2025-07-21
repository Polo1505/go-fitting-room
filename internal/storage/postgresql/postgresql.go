package postgresql

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

type Costume struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Title       string    `gorm:"index"`
	Description string
	Image       string
}

func New(host, port, user, password, dbName string) (*Storage, error) {
	const op = "storage.postgres.New"

	// Формируем строку подключения (DSN) для PostgreSQL
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	// Открываем соединение с PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // Отключаем проверку внешних ключей
	})
	if err != nil {
		return nil, fmt.Errorf("%s: failed to open database: %w", op, err)
	}

	// Включаем расширение uuid-ossp для поддержки uuid_generate_v4()
	err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		return nil, fmt.Errorf("%s: failed to enable uuid-ossp extension: %w", op, err)
	}

	// Автоматическая миграция схемы
	if err := db.AutoMigrate(&Costume{}); err != nil {
		return nil, fmt.Errorf("%s: failed to migrate database schema: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// Close закрывает соединение с базой данных
func (s *Storage) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	return sqlDB.Close()
}

// CreateCostume создаёт новый костюм
func (s *Storage) CreateCostume(costume *Costume) error {
	const op = "storage.postgres.CreateCostume"

	if costume == nil {
		return fmt.Errorf("%s: costume is nil", op)
	}

	// ID будет автоматически сгенерирован PostgreSQL (uuid_generate_v4)
	if err := s.db.Create(costume).Error; err != nil {
		return fmt.Errorf("%s: failed to create costume: %w", op, err)
	}
	return nil
}

// GetCostume получает костюм по ID
func (s *Storage) GetCostume(id uuid.UUID) (*Costume, error) {
	const op = "storage.postgres.GetCostume"

	var costume Costume
	if err := s.db.Where("id = ?", id).First(&costume).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("%s: costume not found: %w", op, err)
		}
		return nil, fmt.Errorf("%s: failed to get costume: %w", op, err)
	}
	return &costume, nil
}

// GetAllCostumes возвращает все костюмы
func (s *Storage) GetAllCostumes() ([]Costume, error) {
	const op = "storage.postgres.GetAllCostumes"

	var costumes []Costume
	if err := s.db.Find(&costumes).Error; err != nil {
		return nil, fmt.Errorf("%s: failed to get all costumes: %w", op, err)
	}
	return costumes, nil
}

// UpdateCostume обновляет существующий костюм
func (s *Storage) UpdateCostume(id uuid.UUID, updatedCostume *Costume) error {
	const op = "storage.postgres.UpdateCostume"

	if updatedCostume == nil {
		return fmt.Errorf("%s: updated costume is nil", op)
	}

	// Проверяем, существует ли костюм
	var costume Costume
	if err := s.db.Where("id = ?", id).First(&costume).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("%s: costume not found: %w", op, err)
		}
		return fmt.Errorf("%s: failed to find costume: %w", op, err)
	}

	// Обновляем поля
	if err := s.db.Model(&costume).Updates(updatedCostume).Error; err != nil {
		return fmt.Errorf("%s: failed to update costume: %w", op, err)
	}
	return nil
}

// DeleteCostume удаляет костюм по ID
func (s *Storage) DeleteCostume(id uuid.UUID) error {
	const op = "storage.postgres.DeleteCostume"

	result := s.db.Where("id = ?", id).Delete(&Costume{})
	if result.Error != nil {
		return fmt.Errorf("%s: failed to delete costume: %w", op, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%s: costume not found", op)
	}
	return nil
}
