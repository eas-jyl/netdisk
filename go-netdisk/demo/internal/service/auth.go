package service

import (
	"errors"
	"fmt"
	"time"

	"go-netdisk/internal/config"
	"go-netdisk/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrInvalidCredentials = errors.New("invalid username or password")

type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

// jwt claims结构体
type TokenClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 创建AuthService对象
func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

// 接口：注册用户
func (s *AuthService) Register(username, passwd string) (*model.User, error) {
	// 查询数据库是否有该用户
	var count int64
	if err := s.db.Model(&model.User{}).Where("username = ?", username).Count(&count); err != nil {
		return nil, fmt.Errorf("check username exists failed: %w", err)
	}
	if count > 0 {
		return nil, fmt.Errorf("username already exists")
	}

	// 生成加密后的密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password failed: %w", err)
	}

	// 创建用户存入数据库
	user := &model.User{
		UUID:     uuid.NewString(), // 	全局唯一的字符串id
		Username: username,
		Password: string(hashedPassword),
	}
	if err := s.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("create user failed: %w", err)
	}

	return user, nil
}

// 接口：登陆用户（用户名和密码正确，生成jwt token）
func (s *AuthService) Login(username, passwd string) (*model.User, string, error) {
	var user model.User

	// 验证用户名
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		// 数据库没有这个用户名
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", ErrInvalidCredentials
		}
		// 查询数据库出现问题
		return nil, "", fmt.Errorf("query user failed: %w", err)
	}

	// 验证密码（数据库加密的密码，明文密码）
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwd)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// 生成token
	token, err := s.generateToken(&user)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

// 接口：产生token
func (s *AuthService) generateToken(user *model.User) (string, error) {
	// 计算截止时间
	expiresAt := time.Now().Add(time.Duration(s.cfg.JWT.ExpireHours) * time.Hour)

	// 创建claims
	claims := TokenClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("sign token failed: %w", err)
	}
	return signedToken, nil
}
