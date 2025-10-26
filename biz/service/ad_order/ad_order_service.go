package ad_order

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"orbia_api/biz/dal/mysql"
	adOrderModel "orbia_api/biz/model/ad_order"
	"orbia_api/biz/utils"
)

var (
	adOrderRepo mysql.AdOrderRepository
)

// InitAdOrderService 初始化广告订单服务
func InitAdOrderService() {
	adOrderRepo = mysql.NewAdOrderRepository(mysql.DB)
}

// CreateAdOrder 创建广告订单
func CreateAdOrder(userID int64, req *adOrderModel.CreateAdOrderReq) (*adOrderModel.CreateAdOrderResp, error) {
	resp := &adOrderModel.CreateAdOrderResp{}

	// 1. 如果指定了团队ID，验证用户是否属于该团队
	if req.TeamID != nil {
		// TODO: 验证用户是否属于该团队（需要团队仓储支持）
	}

	// 2. 生成订单ID (ADORD_ 前缀表示 Ad Order)
	orderID := utils.GenerateAdOrderID()

	// 3. 创建订单
	order := &mysql.AdOrder{
		OrderID:        orderID,
		UserID:         userID,
		TeamID:         req.TeamID,
		Title:          req.Title,
		Description:    req.Description,
		Budget:         req.Budget,
		AdType:         req.AdType,
		TargetAudience: req.TargetAudience,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		Status:         "pending", // 默认状态为待审核
	}

	if err := adOrderRepo.CreateAdOrder(order); err != nil {
		return nil, fmt.Errorf("创建广告订单失败: %w", err)
	}

	resp.OrderID = &orderID
	return resp, nil
}

// GetAdOrder 获取广告订单详情
func GetAdOrder(userID int64, isAdmin bool, req *adOrderModel.GetAdOrderReq) (*adOrderModel.GetAdOrderResp, error) {
	resp := &adOrderModel.GetAdOrderResp{}

	// 1. 获取订单（包含用户信息）
	orderWithUser, err := adOrderRepo.GetAdOrderWithUserInfo(req.OrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("订单不存在")
		}
		return nil, fmt.Errorf("获取订单失败: %w", err)
	}

	// 2. 权限验证：只有订单创建者、管理员可以查看
	if !isAdmin && orderWithUser.UserID != userID {
		return nil, fmt.Errorf("无权查看该订单")
	}

	// 3. 转换为响应模型
	resp.Order = convertToAdOrderInfo(orderWithUser)
	return resp, nil
}

// GetUserAdOrderList 获取用户自己的广告订单列表
func GetUserAdOrderList(userID int64, req *adOrderModel.GetUserAdOrderListReq) (*adOrderModel.GetUserAdOrderListResp, error) {
	resp := &adOrderModel.GetUserAdOrderListResp{}

	// 计算分页参数
	page, pageSize := utils.GetPageParams(req.Page, req.PageSize)
	offset := (page - 1) * pageSize

	// 获取订单列表
	orders, total, err := adOrderRepo.GetUserAdOrders(userID, req.Status, req.Keyword, req.AdType, offset, pageSize)
	if err != nil {
		return nil, fmt.Errorf("获取订单列表失败: %w", err)
	}

	// 转换为响应模型
	resp.Orders = make([]*adOrderModel.AdOrderInfo, 0, len(orders))
	for _, order := range orders {
		resp.Orders = append(resp.Orders, convertToAdOrderInfo(order))
	}
	resp.Total = total

	return resp, nil
}

// GetAdOrderList 获取所有广告订单列表（管理员使用）
func GetAdOrderList(req *adOrderModel.GetAdOrderListReq) (*adOrderModel.GetAdOrderListResp, error) {
	resp := &adOrderModel.GetAdOrderListResp{}

	// 计算分页参数
	page, pageSize := utils.GetPageParams(req.Page, req.PageSize)
	offset := (page - 1) * pageSize

	// 获取订单列表
	orders, total, err := adOrderRepo.GetAllAdOrders(req.Status, req.Keyword, req.AdType, offset, pageSize)
	if err != nil {
		return nil, fmt.Errorf("获取订单列表失败: %w", err)
	}

	// 转换为响应模型
	resp.Orders = make([]*adOrderModel.AdOrderInfo, 0, len(orders))
	for _, order := range orders {
		resp.Orders = append(resp.Orders, convertToAdOrderInfo(order))
	}
	resp.Total = total

	return resp, nil
}

// UpdateAdOrderStatus 更新广告订单状态（管理员使用）
func UpdateAdOrderStatus(req *adOrderModel.UpdateAdOrderStatusReq) (*adOrderModel.UpdateAdOrderStatusResp, error) {
	resp := &adOrderModel.UpdateAdOrderStatusResp{}

	// 1. 获取订单
	order, err := adOrderRepo.GetAdOrderByID(req.OrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("订单不存在")
		}
		return nil, fmt.Errorf("获取订单失败: %w", err)
	}

	// 2. 验证状态转换是否合法
	validTransitions := map[string][]string{
		"pending":     {"approved", "cancelled"},
		"approved":    {"in_progress", "cancelled"},
		"in_progress": {"completed", "cancelled"},
	}

	allowedStatuses, exists := validTransitions[order.Status]
	if !exists {
		return nil, fmt.Errorf("当前订单状态无法变更")
	}

	isAllowed := false
	for _, allowed := range allowedStatuses {
		if req.Status == allowed {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		return nil, fmt.Errorf("不允许从 %s 状态转换到 %s 状态", order.Status, req.Status)
	}

	// 3. 更新订单状态
	if err := adOrderRepo.UpdateAdOrderStatus(req.OrderID, req.Status, req.RejectReason); err != nil {
		return nil, fmt.Errorf("更新订单状态失败: %w", err)
	}

	return resp, nil
}

// CancelAdOrder 取消广告订单（用户使用）
func CancelAdOrder(userID int64, req *adOrderModel.CancelAdOrderReq) (*adOrderModel.CancelAdOrderResp, error) {
	resp := &adOrderModel.CancelAdOrderResp{}

	// 1. 获取订单
	order, err := adOrderRepo.GetAdOrderByID(req.OrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("订单不存在")
		}
		return nil, fmt.Errorf("获取订单失败: %w", err)
	}

	// 2. 验证订单是否属于该用户
	if order.UserID != userID {
		return nil, fmt.Errorf("无权操作该订单")
	}

	// 3. 验证订单状态是否可以取消
	if order.Status == "completed" || order.Status == "cancelled" {
		return nil, fmt.Errorf("该订单无法取消")
	}

	// 4. 取消订单
	if err := adOrderRepo.UpdateAdOrderStatus(req.OrderID, "cancelled", &req.Reason); err != nil {
		return nil, fmt.Errorf("取消订单失败: %w", err)
	}

	return resp, nil
}

// convertToAdOrderInfo 转换为广告订单信息模型
func convertToAdOrderInfo(order *mysql.AdOrderWithUserInfo) *adOrderModel.AdOrderInfo {
	info := &adOrderModel.AdOrderInfo{
		OrderID:        order.OrderID,
		UserID:         order.UserID,
		Title:          order.Title,
		Description:    order.Description,
		Budget:         order.Budget,
		AdType:         order.AdType,
		TargetAudience: order.TargetAudience,
		StartDate:      order.StartDate,
		EndDate:        order.EndDate,
		Status:         order.Status,
		RejectReason:   strPtrToOptional(order.RejectReason),
		CreatedAt:      order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      order.UpdatedAt.Format(time.RFC3339),
	}

	if order.TeamID != nil {
		info.TeamID = order.TeamID
	}

	if order.UserNickname != nil {
		info.UserNickname = *order.UserNickname
	}

	if order.TeamName != nil {
		info.TeamName = strPtrToOptional(order.TeamName)
	}

	if order.ApprovedAt != nil {
		approvedAt := order.ApprovedAt.Format(time.RFC3339)
		info.ApprovedAt = &approvedAt
	}

	if order.CompletedAt != nil {
		completedAt := order.CompletedAt.Format(time.RFC3339)
		info.CompletedAt = &completedAt
	}

	if order.CancelledAt != nil {
		cancelledAt := order.CancelledAt.Format(time.RFC3339)
		info.CancelledAt = &cancelledAt
	}

	return info
}

// strPtrToOptional 字符串指针转可选字符串
func strPtrToOptional(s *string) *string {
	if s == nil || *s == "" {
		return nil
	}
	return s
}
