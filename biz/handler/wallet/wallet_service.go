package wallet

import (
	"context"
	"strconv"

	"orbia_api/biz/dal/model"
	"orbia_api/biz/dal/mysql"
	walletModel "orbia_api/biz/model/wallet"
	"orbia_api/biz/mw"
	walletService "orbia_api/biz/service/wallet"
	"orbia_api/biz/utils"

	"github.com/cloudwego/hertz/pkg/app"
)

var (
	walletSvc walletService.WalletService
)

// InitWalletHandler 初始化钱包 handler
func InitWalletHandler() {
	db := mysql.DB
	walletRepo := mysql.NewWalletRepository(db)
	txRepo := mysql.NewTransactionRepository(db)
	walletSvc = walletService.NewWalletService(db, walletRepo, txRepo)
}

// GetWalletInfo 获取钱包信息
// @router /api/v1/wallet/info [POST]
func GetWalletInfo(ctx context.Context, c *app.RequestContext) {
	var req walletModel.GetWalletInfoReq
	if err := c.BindAndValidate(&req); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	// 从上下文获取用户ID
	userID, exists := mw.GetAuthUserID(c)
	if !exists {
		utils.ErrorResponse(c, 401, "unauthorized")
		return
	}

	// 获取钱包信息
	walletData, err := walletSvc.GetWalletInfo(userID)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	// 构建响应
	walletInfo := &walletModel.WalletInfo{
		ID:            walletData.ID,
		UserID:        walletData.UserID,
		Balance:       formatAmount(walletData.Balance),
		FrozenBalance: formatAmount(walletData.FrozenBalance),
		TotalRecharge: formatAmount(walletData.TotalRecharge),
		TotalConsume:  formatAmount(walletData.TotalConsume),
		CreatedAt:     utils.FormatTime(walletData.CreatedAt),
		UpdatedAt:     utils.FormatTime(walletData.UpdatedAt),
	}

	utils.SuccessResponse(c, map[string]interface{}{
		"wallet": walletInfo,
	})
}

// CryptoRecharge 加密货币充值
// @router /api/v1/wallet/recharge/crypto [POST]
func CryptoRecharge(ctx context.Context, c *app.RequestContext) {
	var req walletModel.CryptoRechargeReq
	if err := c.BindAndValidate(&req); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	// 从上下文获取用户ID
	userID, exists := mw.GetAuthUserID(c)
	if !exists {
		utils.ErrorResponse(c, 401, "unauthorized")
		return
	}

	// 解析金额
	amount, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid amount format")
		return
	}

	// 创建充值交易
	transaction, err := walletSvc.CryptoRecharge(
		userID,
		amount,
		req.CryptoCurrency,
		req.CryptoChain,
		req.CryptoAddress,
	)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	// 构建响应
	txResp := buildTransactionResponse(transaction)
	utils.SuccessResponse(c, map[string]interface{}{
		"transaction": txResp,
	})
}

// OnlineRecharge 在线支付充值
// @router /api/v1/wallet/recharge/online [POST]
func OnlineRecharge(ctx context.Context, c *app.RequestContext) {
	var req walletModel.OnlineRechargeReq
	if err := c.BindAndValidate(&req); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	// 从上下文获取用户ID
	userID, exists := mw.GetAuthUserID(c)
	if !exists {
		utils.ErrorResponse(c, 401, "unauthorized")
		return
	}

	// 解析金额
	amount, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid amount format")
		return
	}

	// 创建充值交易
	transaction, paymentURL, err := walletSvc.OnlineRecharge(
		userID,
		amount,
		req.Platform,
	)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	// 构建响应
	txResp := buildTransactionResponse(transaction)
	utils.SuccessResponse(c, map[string]interface{}{
		"transaction": txResp,
		"payment_url": paymentURL,
	})
}

// ConfirmCryptoRecharge 确认加密货币充值
// @router /api/v1/wallet/recharge/crypto/confirm [POST]
func ConfirmCryptoRecharge(ctx context.Context, c *app.RequestContext) {
	var req walletModel.ConfirmCryptoRechargeReq
	if err := c.BindAndValidate(&req); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	// 从上下文获取用户ID
	userID, exists := mw.GetAuthUserID(c)
	if !exists {
		utils.ErrorResponse(c, 401, "unauthorized")
		return
	}

	// 确认充值
	transaction, err := walletSvc.ConfirmCryptoRecharge(
		userID,
		req.TransactionID,
		req.CryptoTxHash,
	)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	// 构建响应
	txResp := buildTransactionResponse(transaction)
	utils.SuccessResponse(c, map[string]interface{}{
		"transaction": txResp,
	})
}

// GetTransactionList 获取交易记录列表
// @router /api/v1/wallet/transactions [POST]
func GetTransactionList(ctx context.Context, c *app.RequestContext) {
	var req walletModel.GetTransactionListReq
	if err := c.BindAndValidate(&req); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	// 从上下文获取用户ID
	userID, exists := mw.GetAuthUserID(c)
	if !exists {
		utils.ErrorResponse(c, 401, "unauthorized")
		return
	}

	// 设置默认值
	page := int(1)
	if req.Page != nil && *req.Page > 0 {
		page = int(*req.Page)
	}

	pageSize := int(20)
	if req.PageSize != nil && *req.PageSize > 0 {
		pageSize = int(*req.PageSize)
	}

	// 获取交易列表
	transactions, total, err := walletSvc.GetTransactionList(
		userID,
		req.Type,
		req.Status,
		page,
		pageSize,
	)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	// 构建响应
	txList := make([]*walletModel.Transaction, 0, len(transactions))
	for _, tx := range transactions {
		txList = append(txList, buildTransactionResponse(tx))
	}

	utils.SuccessResponse(c, map[string]interface{}{
		"transactions": txList,
		"total":        total,
		"page":         page,
		"page_size":    pageSize,
	})
}

// GetTransactionDetail 获取交易详情
// @router /api/v1/wallet/transaction/:transaction_id [POST]
func GetTransactionDetail(ctx context.Context, c *app.RequestContext) {
	var req walletModel.GetTransactionDetailReq
	if err := c.BindAndValidate(&req); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	// 从上下文获取用户ID
	userID, exists := mw.GetAuthUserID(c)
	if !exists {
		utils.ErrorResponse(c, 401, "unauthorized")
		return
	}

	// 获取交易详情
	transaction, err := walletSvc.GetTransactionDetail(userID, req.TransactionID)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	// 构建响应
	txResp := buildTransactionResponse(transaction)
	utils.SuccessResponse(c, map[string]interface{}{
		"transaction": txResp,
	})
}

// 辅助函数：构建交易响应
func buildTransactionResponse(tx *model.OrbiaTransaction) *walletModel.Transaction {
	resp := &walletModel.Transaction{
		ID:            tx.ID,
		TransactionID: tx.TransactionID,
		UserID:        tx.UserID,
		Type:          tx.Type,
		Amount:        formatAmount(tx.Amount),
		BalanceBefore: formatAmount(tx.BalanceBefore),
		BalanceAfter:  formatAmount(tx.BalanceAfter),
		Status:        tx.Status,
		CreatedAt:     utils.FormatTime(tx.CreatedAt),
		UpdatedAt:     utils.FormatTime(tx.UpdatedAt),
	}

	// Transaction表已简化，只保留消费相关字段
	if tx.RelatedOrderType != nil {
		relatedOrderType := *tx.RelatedOrderType
		resp.PaymentMethod = &relatedOrderType // 复用PaymentMethod字段存储订单类型
	}
	if tx.RelatedOrderID != nil {
		resp.RelatedOrderID = tx.RelatedOrderID
	}
	if tx.Remark != nil {
		resp.Remark = tx.Remark
	}
	if tx.CompletedAt != nil {
		completedAt := utils.FormatTime(tx.CompletedAt)
		resp.CompletedAt = &completedAt
	}

	return resp
}

// 辅助函数：格式化金额
func formatAmount(amount float64) string {
	return strconv.FormatFloat(amount, 'f', 2, 64)
}
