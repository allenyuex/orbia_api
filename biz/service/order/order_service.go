package order

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"orbia_api/biz/dal/mysql"
	orderModel "orbia_api/biz/model/order/order"
	"orbia_api/biz/utils"
)

var (
	orderRepo mysql.OrderRepository
	kolRepo   mysql.KolRepository
)

// InitOrderService 初始化订单服务
func InitOrderService() {
	orderRepo = mysql.NewOrderRepository(mysql.DB)
	kolRepo = mysql.NewKolRepository(mysql.DB)
}

// CreateOrder 创建订单
func CreateOrder(userID int64, req *orderModel.CreateOrderReq) (*orderModel.CreateOrderResp, error) {
	resp := &orderModel.CreateOrderResp{}

	// 1. 验证 KOL 是否存在且已审核通过
	kol, err := kolRepo.GetKolByID(req.KolID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("KOL 不存在")
		}
		return nil, fmt.Errorf("获取 KOL 信息失败: %w", err)
	}
	if kol.Status != "approved" {
		return nil, fmt.Errorf("该 KOL 尚未通过审核")
	}

	// 2. 验证 Plan 是否存在
	plan, err := kolRepo.GetKolPlanByID(req.PlanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("报价计划不存在")
		}
		return nil, fmt.Errorf("获取报价计划失败: %w", err)
	}

	// 3. 验证 Plan 是否属于该 KOL
	if plan.KolID != req.KolID {
		return nil, fmt.Errorf("报价计划不属于该 KOL")
	}

	// 4. 如果指定了团队ID，验证用户是否属于该团队
	if req.TeamID != nil {
		// TODO: 验证用户是否属于该团队（需要团队仓储支持）
	}

	// 5. 生成订单ID
	orderID, err := utils.GenerateOrderID()
	if err != nil {
		return nil, fmt.Errorf("生成订单ID失败: %w", err)
	}

	// 6. 创建订单（保存 Plan 快照）
	order := &mysql.KolOrder{
		OrderID:         orderID,
		UserID:          userID,
		TeamID:          req.TeamID,
		KolID:           req.KolID,
		PlanID:          req.PlanID,
		PlanTitle:       plan.Title,
		PlanDescription: plan.Description,
		PlanPrice:       plan.Price,
		PlanType:        plan.PlanType,
		Description:     req.Description,
		Status:          "pending",
	}

	if err := orderRepo.CreateOrder(order); err != nil {
		return nil, fmt.Errorf("创建订单失败: %w", err)
	}

	resp.OrderID = &orderID
	return resp, nil
}

// GetOrder 获取订单详情
func GetOrder(userID int64, req *orderModel.GetOrderReq) (*orderModel.GetOrderResp, error) {
	resp := &orderModel.GetOrderResp{}

	// 1. 获取订单（包含 KOL 信息）
	orderWithKol, err := orderRepo.GetOrderWithKolInfo(req.OrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("订单不存在")
		}
		return nil, fmt.Errorf("获取订单失败: %w", err)
	}

	// 2. 权限检查：只能查看自己的订单或收到的订单
	kol, _ := kolRepo.GetKolByUserID(userID)
	if orderWithKol.UserID != userID && (kol == nil || orderWithKol.KolID != kol.ID) {
		return nil, fmt.Errorf("无权查看该订单")
	}

	// 3. 转换为响应格式
	resp.Order = convertToOrderInfo(orderWithKol)
	return resp, nil
}

// GetOrderList 获取用户的订单列表
func GetOrderList(userID int64, req *orderModel.GetOrderListReq) (*orderModel.GetOrderListResp, error) {
	resp := &orderModel.GetOrderListResp{}

	page := int32(1)
	if req.Page != nil {
		page = *req.Page
	}
	pageSize := int32(10)
	if req.PageSize != nil {
		pageSize = *req.PageSize
	}

	offset := int((page - 1) * pageSize)
	limit := int(pageSize)

	var orders []*mysql.OrderWithKolInfo
	var total int64
	var err error

	// 根据是否指定 KOL ID 或 Team ID 来查询不同的订单列表
	if req.KolID != nil {
		// 查询指定 KOL 的订单（需要验证权限）
		kol, kolErr := kolRepo.GetKolByID(*req.KolID)
		if kolErr != nil {
			return nil, fmt.Errorf("KOL 不存在")
		}
		if kol.UserID != userID {
			return nil, fmt.Errorf("无权查看该 KOL 的订单")
		}
		orders, total, err = orderRepo.GetKolOrders(*req.KolID, req.Status, offset, limit)
	} else if req.TeamID != nil {
		// 查询团队的订单
		// TODO: 验证用户是否属于该团队
		orders, total, err = orderRepo.GetTeamOrders(*req.TeamID, req.Status, offset, limit)
	} else {
		// 查询用户自己的订单
		orders, total, err = orderRepo.GetUserOrders(userID, req.Status, offset, limit)
	}

	if err != nil {
		return nil, fmt.Errorf("获取订单列表失败: %w", err)
	}

	// 转换为响应格式
	resp.Orders = make([]*orderModel.OrderInfo, 0, len(orders))
	for _, order := range orders {
		resp.Orders = append(resp.Orders, convertToOrderInfo(order))
	}
	resp.Total = total

	return resp, nil
}

// CancelOrder 取消订单
func CancelOrder(userID int64, req *orderModel.CancelOrderReq) (*orderModel.CancelOrderResp, error) {
	resp := &orderModel.CancelOrderResp{}

	// 1. 获取订单
	order, err := orderRepo.GetOrderByID(req.OrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("订单不存在")
		}
		return nil, fmt.Errorf("获取订单失败: %w", err)
	}

	// 2. 权限检查：只有下单用户可以取消
	if order.UserID != userID {
		return nil, fmt.Errorf("无权取消该订单")
	}

	// 3. 状态检查：只有 pending 状态的订单可以取消
	if order.Status != "pending" {
		return nil, fmt.Errorf("订单状态不允许取消")
	}

	// 4. 更新订单状态
	if err := orderRepo.UpdateOrderStatus(req.OrderID, "cancelled", &req.Reason); err != nil {
		return nil, fmt.Errorf("取消订单失败: %w", err)
	}

	return resp, nil
}

// GetKolOrders 获取 KOL 收到的订单列表
func GetKolOrders(userID int64, req *orderModel.GetKolOrdersReq) (*orderModel.GetKolOrdersResp, error) {
	resp := &orderModel.GetKolOrdersResp{}

	// 1. 验证用户是否是 KOL
	kol, err := kolRepo.GetKolByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("您还不是 KOL")
		}
		return nil, fmt.Errorf("获取 KOL 信息失败: %w", err)
	}

	// 2. 获取订单列表
	page := int32(1)
	if req.Page != nil {
		page = *req.Page
	}
	pageSize := int32(10)
	if req.PageSize != nil {
		pageSize = *req.PageSize
	}

	offset := int((page - 1) * pageSize)
	limit := int(pageSize)

	orders, total, err := orderRepo.GetKolOrders(kol.ID, req.Status, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("获取订单列表失败: %w", err)
	}

	// 3. 转换为响应格式
	resp.Orders = make([]*orderModel.OrderInfo, 0, len(orders))
	for _, order := range orders {
		resp.Orders = append(resp.Orders, convertToOrderInfo(order))
	}
	resp.Total = total

	return resp, nil
}

// UpdateOrderStatus 更新订单状态（KOL 端使用）
func UpdateOrderStatus(userID int64, req *orderModel.UpdateOrderStatusReq) (*orderModel.UpdateOrderStatusResp, error) {
	resp := &orderModel.UpdateOrderStatusResp{}

	// 1. 获取订单
	order, err := orderRepo.GetOrderByID(req.OrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("订单不存在")
		}
		return nil, fmt.Errorf("获取订单失败: %w", err)
	}

	// 2. 验证用户是否是该订单的 KOL
	kol, err := kolRepo.GetKolByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("您还不是 KOL")
	}
	if order.KolID != kol.ID {
		return nil, fmt.Errorf("无权操作该订单")
	}

	// 3. 状态转换验证
	if err := validateStatusTransition(order.Status, req.Status); err != nil {
		return nil, err
	}

	// 4. 更新订单状态
	if err := orderRepo.UpdateOrderStatus(req.OrderID, req.Status, req.RejectReason); err != nil {
		return nil, fmt.Errorf("更新订单状态失败: %w", err)
	}

	return resp, nil
}

// convertToOrderInfo 将数据库模型转换为 API 模型
func convertToOrderInfo(order *mysql.OrderWithKolInfo) *orderModel.OrderInfo {
	info := &orderModel.OrderInfo{
		ID:          order.ID,
		OrderID:     order.OrderID,
		UserID:      order.UserID,
		KolID:       order.KolID,
		PlanID:      order.PlanID,
		PlanTitle:   order.PlanTitle,
		PlanPrice:   order.PlanPrice,
		PlanType:    order.PlanType,
		Description: order.Description,
		Status:      order.Status,
	}

	if order.TeamID != nil {
		info.TeamID = order.TeamID
	}

	if order.PlanDescription != nil {
		info.PlanDescription = *order.PlanDescription
	}

	if order.KolDisplayName != nil {
		info.KolDisplayName = *order.KolDisplayName
	}

	if order.KolAvatarURL != nil {
		info.KolAvatarURL = *order.KolAvatarURL
	}

	if order.RejectReason != nil {
		info.RejectReason = order.RejectReason
	}

	if order.ConfirmedAt != nil {
		confirmedAt := order.ConfirmedAt.Format(time.RFC3339)
		info.ConfirmedAt = &confirmedAt
	}

	if order.CompletedAt != nil {
		completedAt := order.CompletedAt.Format(time.RFC3339)
		info.CompletedAt = &completedAt
	}

	if order.CancelledAt != nil {
		cancelledAt := order.CancelledAt.Format(time.RFC3339)
		info.CancelledAt = &cancelledAt
	}

	info.CreatedAt = order.CreatedAt.Format(time.RFC3339)
	info.UpdatedAt = order.UpdatedAt.Format(time.RFC3339)

	return info
}

// validateStatusTransition 验证状态转换是否合法
func validateStatusTransition(currentStatus, newStatus string) error {
	// 定义允许的状态转换
	allowedTransitions := map[string][]string{
		"pending":     {"confirmed", "cancelled"},
		"confirmed":   {"in_progress", "cancelled"},
		"in_progress": {"completed", "cancelled"},
		"completed":   {"refunded"},
		"cancelled":   {},
		"refunded":    {},
	}

	allowed, exists := allowedTransitions[currentStatus]
	if !exists {
		return fmt.Errorf("无效的当前状态: %s", currentStatus)
	}

	for _, status := range allowed {
		if status == newStatus {
			return nil
		}
	}

	return fmt.Errorf("不允许从 %s 转换到 %s", currentStatus, newStatus)
}
