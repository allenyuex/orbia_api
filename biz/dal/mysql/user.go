package mysql

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID               int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	WalletAddress    *string        `gorm:"uniqueIndex;column:wallet_address;size:42" json:"wallet_address"`
	Email            *string        `gorm:"uniqueIndex;column:email;size:255" json:"email"`
	PasswordHash     *string        `gorm:"column:password_hash;size:255" json:"-"`
	VerificationCode *string        `gorm:"column:verification_code;size:10" json:"-"`
	CodeExpiry       *time.Time     `gorm:"column:code_expiry" json:"-"`
	Nickname         *string        `gorm:"column:nickname;size:100" json:"nickname"`
	AvatarURL        *string        `gorm:"column:avatar_url;size:500" json:"avatar_url"`
	Role             string         `gorm:"column:role;type:enum('user','admin');default:'user';not null" json:"role"`
	CurrentTeamID    *int64         `gorm:"column:current_team_id" json:"current_team_id"`
	CreatedAt        time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index;column:deleted_at" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "orbia_user"
}

// UserRepository 用户仓储接口
type UserRepository interface {
	CreateUser(user *User) error
	GetUserByWalletAddress(walletAddress string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int64) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id int64) error
}

// userRepository 用户仓储实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// CreateUser 创建用户
func (r *userRepository) CreateUser(user *User) error {
	return r.db.Create(user).Error
}

// GetUserByWalletAddress 根据钱包地址获取用户
func (r *userRepository) GetUserByWalletAddress(walletAddress string) (*User, error) {
	var user User
	err := r.db.Where("wallet_address = ?", walletAddress).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail 根据邮箱获取用户
func (r *userRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据ID获取用户
func (r *userRepository) GetUserByID(id int64) (*User, error) {
	var user User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (r *userRepository) UpdateUser(user *User) error {
	return r.db.Save(user).Error
}

// DeleteUser 删除用户（软删除）
func (r *userRepository) DeleteUser(id int64) error {
	return r.db.Delete(&User{}, id).Error
}
