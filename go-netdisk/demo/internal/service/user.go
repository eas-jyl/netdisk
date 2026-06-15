package service

import (
	"errors"
	"fmt"

	"go-netdisk/internal/model"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// 接口：通过用户id ， 查询数据库的用户字段
func (s *UserService) GetByID(id uint) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		// 没有记录 & 查询失败
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("query user failed: %w", err)
	}
	return &user, nil
}
