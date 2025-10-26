package kol_order

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"orbia_api/biz/dal/mysql"
	kolOrderModel "orbia_api/biz/model/kol_order"
	"orbia_api/biz/utils"
)

var (
	orderRepo mysql.OrderRepository
	kolRepo   mysql.KolRepository
)

// InitKolOrderService 初始化KOL订单服务
func InitKolOrderService() {
	orderRepo = mysql.NewOrderRepository(mysql.DB)
	kolRepo = mysql.NewKolRepository(mysql.DB)
}

// CreateKolOrder 创建KOL订单
func CreateKolOrder(userID int64, req *kolOrderModel.CreateKolOrderReq) (*kolOrderModel.CreateKolOrderResp, error) {
	resp := &kolOrderModel.CreateKolOrderResp{}

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

	// 5. 生成订单ID (KORD_ 前缀表示 KOL Order)
	orderID := utils.GenerateKolOrderID()

	// 6. 创建订单（保存 Plan 快照）
	order := &mysql.KolOrder{
		OrderID:                orderID,
		UserID:                 userID,
		TeamID:                 req.TeamID,
		KolID:                  req.KolID,
		PlanID:                 req.PlanID,
		PlanTitle:              plan.Title,
		PlanDescription:        plan.Description,
		PlanPrice:              plan.Price,
		PlanType:               plan.PlanType,
		Title:                  req.Title,
		RequirementDescription: req.RequirementDescription,
		VideoType:              req.VideoType,
		VideoDuration:          req.VideoDuration,
		TargetAudience:         req.TargetAudience,
		ExpectedDeliveryDate:   req.ExpectedDeliveryDate,
		AdditionalRequirements: req.AdditionalRequirements,
		Status:                 "pending_payment", // 初始状态为待支付
	}

	if err := orderRepo.CreateOrder(order); err != nil {
		return nil, fmt.Errorf("创建订单失败: %w", err)
	}

	resp.OrderID = &orderID
	return resp, nil
}

// GetKolOrder 获取KOL订单详情
func GetKolOrder(userID int64, req *kolOrderModel.GetKolOrderReq) (*kolOrderModel.GetKolOrderResp, error) {
	resp := &kolOrderModel.GetKolOrderResp{}

	// 1. 获取订单（包含 KOL 信息）
	orderWithKol, err := orderRepo.GetOrderWithKolInfo(req.OrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("订单不存在")
		}
		return nil, fmt.Errorf("获取订单失败: %w", err)
	}

	// 2. 权限验证：只有订单创建者、订单相关的KOL、管理员可以查看
	// 获取用户的 KOL 信息（如果是 KOL）
	userKol, _ := kolRepo.GetKolByUserID(userID)
	isKolOwner := userKol != nil && userKol.ID == orderWithKol.KolID
	isOrderOwner := orderWithKol.UserID == userID

	if !isOrderOwner && !isKolOwner {
		return nil, fmt.Errorf("无权查看该订单")
	}

	// 3. 转换为响应模型
	resp.Order = convertToKolOrderInfo(orderWithKol)
	return resp, nil
}

// GetUserKolOrderList 获取用户自己的KOL订单列表
func GetUserKolOrderList(userID int64, req *kolOrderModel.GetUserKolOrderListReq) (*kolOrderModel.GetUserKolOrderListResp, error) {
	resp := &kolOrderModel.GetUserKolOrderListResp{}

	// 计算分页参数
	page, pageSize := utils.GetPageParams(req.Page, req.PageSize)
	offset := (page - 1) * pageSize

	// 获取订单列表
	orders, total, err := orderRepo.GetUserOrders(userID, req.Status, req.Keyword, req.KolID, offset, pageSize)
	if err != nil {
		return nil, fmt.Errorf("获取订单列表失败: %w", err)
	}

	// 转换为响应模型
	resp.Orders = make([]*kolOrderModel.KolOrderInfo, 0, len(orders))
	for _, order := range orders {
		resp.Orders = append(resp.Orders, convertToKolOrderInfo(order))
	}
	resp.Total = total

	return resp, nil
}

// GetKolReceivedOrderList 获取KOL收到的订单列表
func GetKolReceivedOrderList(userID int64, req *kolOrderModel.GetKolReceivedOrderListReq) (*kolOrderModel.GetKolReceivedOrderListResp, error) {
	resp := &kolOrderModel.GetKolReceivedOrderListResp{}

	// 1. 验证用户是否是 KOL
	kol, err := kolRepo.GetKolByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("您还不是 KOL，无法查看订单")
		}
		return nil, fmt.Errorf("获取 KOL 信息失败: %w", err)
	}

	// 2. 计算分页参数
	page, pageSize := utils.GetPageParams(req.Page, req.PageSize)
	offset := (page - 1) * pageSize

	// 3. 获取订单列表
	orders, total, err := orderRepo.GetKolOrders(kol.ID, req.Status, req.Keyword, offset, pageSize)
	if err != nil {
		return nil, fmt.Errorf("获取订单列表失败: %w", err)
	}

	// 4. 转换为响应模型
	resp.Orders = make([]*kolOrderModel.KolOrderInfo, 0, len(orders))
	for _, order := range orders {
		resp.Orders = append(resp.Orders, convertToKolOrderInfo(order))
	}
	resp.Total = total

	return resp, nil
}

// UpdateKolOrderStatus 更新KOL订单状态（KOL使用）
func UpdateKolOrderStatus(userID int64, req *kolOrderModel.UpdateKolOrderStatusReq) (*kolOrderModel.UpdateKolOrderStatusResp, error) {
	resp := &kolOrderModel.UpdateKolOrderStatusResp{}

	// 1. 验证用户是否是 KOL
	kol, err := kolRepo.GetKolByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("您还不是 KOL，无法操作订单")
		}
		return nil, fmt.Errorf("获取 KOL 信息失败: %w", err)
	}

	// 2. 获取订单
	order, err := orderRepo.GetOrderByID(req.OrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("订单不存在")
		}
		return nil, fmt.Errorf("获取订单失败: %w", err)
	}

	// 3. 验证订单是否属于该 KOL
	if order.KolID != kol.ID {
		return nil, fmt.Errorf("无权操作该订单")
	}

	// 4. 验证状态转换是否合法
	validTransitions := map[string][]string{
		"pending":     {"confirmed", "cancelled"},   // 待确认 -> 已确认/已取消
		"confirmed":   {"in_progress", "cancelled"}, // 已确认 -> 进行中/已取消
		"in_progress": {"completed", "cancelled"},   // 进行中 -> 已完成/已取消
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

	// 5. 更新订单状态
	if err := orderRepo.UpdateOrderStatus(req.OrderID, req.Status, req.RejectReason); err != nil {
		return nil, fmt.Errorf("更新订单状态失败: %w", err)
	}

	return resp, nil
}

// ConfirmKolOrderPayment 确认KOL订单支付（用户支付完成后调用）
func ConfirmKolOrderPayment(userID int64, req *kolOrderModel.ConfirmKolOrderPaymentReq) (*kolOrderModel.ConfirmKolOrderPaymentResp, error) {
	resp := &kolOrderModel.ConfirmKolOrderPaymentResp{}

	// 1. 获取订单
	order, err := orderRepo.GetOrderByID(req.OrderID)
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

	// 3. 验证订单状态是否为待支付
	if order.Status != "pending_payment" {
		return nil, fmt.Errorf("订单状态不是待支付，无法确认支付")
	}

	// 4. 更新订单状态为待确认（等待KOL确认）
	if err := orderRepo.UpdateOrderStatus(req.OrderID, "pending", nil); err != nil {
		return nil, fmt.Errorf("更新订单状态失败: %w", err)
	}

	return resp, nil
}

// CancelKolOrder 取消KOL订单（用户使用）
func CancelKolOrder(userID int64, req *kolOrderModel.CancelKolOrderReq) (*kolOrderModel.CancelKolOrderResp, error) {
	resp := &kolOrderModel.CancelKolOrderResp{}

	// 1. 获取订单
	order, err := orderRepo.GetOrderByID(req.OrderID)
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
	// 只有已完成、已取消、已退款状态不能取消
	if order.Status == "completed" || order.Status == "cancelled" || order.Status == "refunded" {
		return nil, fmt.Errorf("该订单无法取消")
	}

	// 4. 取消订单
	if err := orderRepo.UpdateOrderStatus(req.OrderID, "cancelled", &req.Reason); err != nil {
		return nil, fmt.Errorf("取消订单失败: %w", err)
	}

	return resp, nil
}

// convertToKolOrderInfo 转换为 KOL 订单信息模型
func convertToKolOrderInfo(order *mysql.OrderWithKolInfo) *kolOrderModel.KolOrderInfo {
	planDesc := ""
	if order.PlanDescription != nil {
		planDesc = *order.PlanDescription
	}

	info := &kolOrderModel.KolOrderInfo{
		OrderID:                order.OrderID,
		UserID:                 order.UserID,
		KolID:                  order.KolID,
		PlanID:                 order.PlanID,
		PlanTitle:              order.PlanTitle,
		PlanDescription:        planDesc,
		PlanPrice:              order.PlanPrice,
		PlanType:               order.PlanType,
		Title:                  order.Title,
		RequirementDescription: order.RequirementDescription,
		VideoType:              order.VideoType,
		VideoDuration:          order.VideoDuration,
		TargetAudience:         order.TargetAudience,
		ExpectedDeliveryDate:   order.ExpectedDeliveryDate,
		AdditionalRequirements: strPtrToOptional(order.AdditionalRequirements),
		Status:                 order.Status,
		RejectReason:           strPtrToOptional(order.RejectReason),
		CreatedAt:              order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:              order.UpdatedAt.Format(time.RFC3339),
	}

	if order.TeamID != nil {
		info.TeamID = order.TeamID
	}

	if order.KolDisplayName != nil {
		info.KolDisplayName = *order.KolDisplayName
	}

	if order.KolAvatarURL != nil {
		info.KolAvatarURL = *order.KolAvatarURL
	}

	if order.UserNickname != nil {
		info.UserNickname = *order.UserNickname
	}

	if order.TeamName != nil {
		info.TeamName = strPtrToOptional(order.TeamName)
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

	return info
}

// strPtrToOptional 字符串指针转可选字符串
func strPtrToOptional(s *string) *string {
	if s == nil || *s == "" {
		return nil
	}
	return s
}
