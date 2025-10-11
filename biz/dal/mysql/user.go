package mysql

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Email     string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Phone     string         `gorm:"size:20" json:"phone"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// UserDAO 用户数据访问对象
type UserDAO struct {
	db *gorm.DB
}

// NewUserDAO 创建用户 DAO
func NewUserDAO() *UserDAO {
	return &UserDAO{db: DB}
}

// Create 创建用户
func (dao *UserDAO) Create(user *User) error {
	return dao.db.Create(user).Error
}

// GetByID 根据 ID 获取用户
func (dao *UserDAO) GetByID(id uint) (*User, error) {
	var user User
	err := dao.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (dao *UserDAO) GetByEmail(email string) (*User, error) {
	var user User
	err := dao.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// List 获取用户列表
func (dao *UserDAO) List(page, pageSize int) ([]*User, int64, error) {
	var users []*User
	var total int64

	// 计算总数
	if err := dao.db.Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := dao.db.Offset(offset).Limit(pageSize).Order("id DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Update 更新用户
func (dao *UserDAO) Update(user *User) error {
	return dao.db.Save(user).Error
}

// Delete 删除用户（软删除）
func (dao *UserDAO) Delete(id uint) error {
	return dao.db.Delete(&User{}, id).Error
}

// ExistsByEmail 检查邮箱是否已存在
func (dao *UserDAO) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := dao.db.Model(&User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
