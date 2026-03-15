package data

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"user-center/internal/conf"
	"user-center/internal/entity"
)

type DB struct {
	*gorm.DB
}

func NewDB(cfg *conf.Database) (*DB, error) {
	var db *gorm.DB
	var err error

	// Use MySQL only
	if cfg == nil || cfg.Source == "" {
		return nil, fmt.Errorf("database config is required")
	}
	
	db, err = gorm.Open(mysql.Open(cfg.Source), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect MySQL: %v", err)
	}
	fmt.Println("Using MySQL database:", cfg.Source)

	// Auto migrate
	if err := db.AutoMigrate(&entity.User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate: %v", err)
	}

	return &DB{db}, nil
}

func (d *DB) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *DB) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := d.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (d *DB) GetUserByPhone(phone string) (*entity.User, error) {
	var user entity.User
	if err := d.Where("phone = ?", phone).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (d *DB) CreateUser(user *entity.User) error {
	return d.Create(user).Error
}

func (d *DB) GetUserByID(id int64) (*entity.User, error) {
	var user entity.User
	if err := d.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (d *DB) UpdateUser(user *entity.User) error {
	return d.Save(user).Error
}

func (d *DB) DeleteUser(id int64) error {
	return d.Delete(&entity.User{}, id).Error
}
