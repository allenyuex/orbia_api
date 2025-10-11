package service

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"orbia_api/biz/consts"
	"orbia_api/biz/dal/mysql"
	"orbia_api/biz/model/api"
)

type UserService struct {
	userDAO *mysql.UserDAO
}

func NewUserService() *UserService {
	return &UserService{
		userDAO: mysql.NewUserDAO(),
	}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(req *api.CreateUserReq) (*api.CreateUserResp, error) {
	// 检查邮箱是否已存在
	exists, err := s.userDAO.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return &api.CreateUserResp{
			BaseResp: &api.BaseResp{
				Code:    consts.UserAlreadyExistCode,
				Message: consts.UserAlreadyExistMsg,
			},
		}, nil
	}

	// 创建用户
	user := &mysql.User{
		Name:  req.Name,
		Email: req.Email,
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}

	if err := s.userDAO.Create(user); err != nil {
		return nil, err
	}

	return &api.CreateUserResp{
		UserID: int64(user.ID),
		BaseResp: &api.BaseResp{
			Code:    consts.SuccessCode,
			Message: consts.SuccessMsg,
		},
	}, nil
}

// GetUser 获取用户信息
func (s *UserService) GetUser(req *api.GetUserReq) (*api.GetUserResp, error) {
	user, err := s.userDAO.GetByID(uint(req.UserID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &api.GetUserResp{
				BaseResp: &api.BaseResp{
					Code:    consts.UserNotFoundCode,
					Message: consts.UserNotFoundMsg,
				},
			}, nil
		}
		return nil, err
	}

	return &api.GetUserResp{
		User: &api.User{
			ID:        int64(user.ID),
			Name:      user.Name,
			Email:     user.Email,
			Phone:     &user.Phone,
			CreatedAt: user.CreatedAt.Format(consts.DateTimeFormat),
			UpdatedAt: user.UpdatedAt.Format(consts.DateTimeFormat),
		},
		BaseResp: &api.BaseResp{
			Code:    consts.SuccessCode,
			Message: consts.SuccessMsg,
		},
	}, nil
}

// ListUsers 获取用户列表
func (s *UserService) ListUsers(req *api.ListUsersReq) (*api.ListUsersResp, error) {
	page := consts.DefaultPage
	pageSize := consts.DefaultPageSize

	if req.Page != nil && *req.Page > 0 {
		page = int(*req.Page)
	}
	if req.PageSize != nil && *req.PageSize > 0 {
		pageSize = int(*req.PageSize)
		if pageSize > consts.MaxPageSize {
			pageSize = consts.MaxPageSize
		}
	}

	users, total, err := s.userDAO.List(page, pageSize)
	if err != nil {
		return nil, err
	}

	userList := make([]*api.User, 0, len(users))
	for _, user := range users {
		userList = append(userList, &api.User{
			ID:        int64(user.ID),
			Name:      user.Name,
			Email:     user.Email,
			Phone:     &user.Phone,
			CreatedAt: user.CreatedAt.Format(consts.DateTimeFormat),
			UpdatedAt: user.UpdatedAt.Format(consts.DateTimeFormat),
		})
	}

	return &api.ListUsersResp{
		Users: userList,
		Total: int32(total),
		BaseResp: &api.BaseResp{
			Code:    consts.SuccessCode,
			Message: consts.SuccessMsg,
		},
	}, nil
}

// Hello Demo 服务
func (s *UserService) Hello(req *api.HelloReq) (*api.HelloResp, error) {
	name := req.Name
	if name == "" {
		name = "World"
	}

	return &api.HelloResp{
		Message:   "Hello, " + name + "! Welcome to Orbia API",
		Timestamp: time.Now().Unix(),
		BaseResp: &api.BaseResp{
			Code:    consts.SuccessCode,
			Message: consts.SuccessMsg,
		},
	}, nil
}
