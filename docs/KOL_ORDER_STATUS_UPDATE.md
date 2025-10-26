# KOL订单状态更新说明

## 概述
本次更新为KOL订单新增了"待支付"状态，优化了订单流程，使其更符合实际业务场景。

## 订单状态变更

### 旧状态流程
1. pending - 待确认
2. confirmed - 已确认
3. in_progress - 进行中
4. completed - 已完成
5. cancelled - 已取消
6. refunded - 已退款

### 新状态流程
1. **pending_payment - 待支付** （新增）
2. pending - 待确认
3. confirmed - 已确认
4. in_progress - 进行中
5. completed - 已完成
6. cancelled - 已取消
7. refunded - 已退款

## 状态流转说明

### 1. 下单流程
```
用户创建订单 -> pending_payment（待支付）
   ↓
用户完成支付 -> pending（待确认）
   ↓
KOL确认订单 -> confirmed（已确认）
   ↓
KOL开始制作 -> in_progress（进行中）
   ↓
KOL完成交付 -> completed（已完成）
```

### 2. 取消/退款流程
- **用户可以取消的状态**：pending_payment、pending、confirmed、in_progress
- **KOL可以拒绝的状态**：pending（拒绝后变为cancelled）
- **退款状态**：任何状态都可能变为refunded（需要管理员操作）

## 新增API接口

### 确认订单支付
**接口**: `POST /api/v1/kol-order/payment/confirm`

**说明**: 用户支付完成后调用此接口，将订单状态从"待支付"变更为"待确认"

**请求参数**:
```json
{
  "order_id": "KORD_1234567890_abc123"
}
```

**响应**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

**错误说明**:
- 订单不存在
- 无权操作该订单
- 订单状态不是待支付，无法确认支付

## 数据库变更

### orbia_kol_order 表
- `status` 字段新增 `pending_payment` 状态
- 默认值从 `pending` 改为 `pending_payment`

### 迁移说明
```sql
-- 修改 status 字段的 ENUM 类型
ALTER TABLE orbia_kol_order 
MODIFY COLUMN status ENUM('pending_payment', 'pending', 'confirmed', 'in_progress', 'completed', 'cancelled', 'refunded') 
NOT NULL DEFAULT 'pending_payment' 
COMMENT '订单状态：pending_payment-待支付，pending-待确认，confirmed-已确认，in_progress-进行中，completed-已完成，cancelled-已取消，refunded-已退款';
```

## 业务逻辑变更

### 1. CreateKolOrder (创建订单)
- **旧行为**: 创建订单后状态为 `pending`
- **新行为**: 创建订单后状态为 `pending_payment`

### 2. CancelKolOrder (取消订单)
- **不变**: 所有非终态（completed、cancelled、refunded）的订单都可以取消
- **包括**: 现在 `pending_payment` 状态也可以取消

### 3. UpdateKolOrderStatus (KOL更新订单状态)
- **不变**: KOL只能操作 `pending` 状态及之后的订单
- **说明**: `pending_payment` 状态的订单对KOL不可见，必须先完成支付

## 前端集成建议

### 1. 订单创建流程
```javascript
// Step 1: 创建订单
const createResp = await createKolOrder({
  kol_id: 123,
  plan_id: 456,
  // ... 其他参数
});

const orderId = createResp.data.order_id;

// Step 2: 跳转到支付页面
// 用户完成支付后...

// Step 3: 确认支付
await confirmKolOrderPayment({
  order_id: orderId
});

// Step 4: 跳转到订单详情页
```

### 2. 订单状态显示
```javascript
const statusMap = {
  'pending_payment': '待支付',
  'pending': '待确认',
  'confirmed': '已确认',
  'in_progress': '进行中',
  'completed': '已完成',
  'cancelled': '已取消',
  'refunded': '已退款'
};
```

### 3. 订单操作权限
```javascript
// 用户可以取消订单
const canCancel = !['completed', 'cancelled', 'refunded'].includes(order.status);

// 用户可以支付订单
const canPay = order.status === 'pending_payment';

// KOL可以操作订单（确认/拒绝/进行中/完成）
const kolCanOperate = ['pending', 'confirmed', 'in_progress'].includes(order.status);
```

## 注意事项

1. **向后兼容性**: 
   - 已有的 `pending` 状态订单不受影响
   - 新创建的订单将从 `pending_payment` 开始

2. **支付集成**: 
   - 需要与支付系统集成，在支付成功回调中调用 `ConfirmKolOrderPayment` 接口
   - 支付失败时订单保持 `pending_payment` 状态，用户可以重新支付或取消

3. **KOL订单列表**: 
   - KOL的订单列表不会显示 `pending_payment` 状态的订单
   - 只有完成支付变为 `pending` 后，KOL才能看到该订单

4. **订单超时处理**（建议）: 
   - 可以设置 `pending_payment` 状态的超时时间（如30分钟）
   - 超时后自动取消订单

## 测试建议

### 1. 正常流程测试
- 创建订单 -> 确认支付 -> KOL确认 -> 进行中 -> 完成

### 2. 异常流程测试
- 创建订单后直接取消
- 创建订单但不支付
- 重复确认支付
- 不同用户尝试确认支付

### 3. 权限测试
- 非订单所有者尝试确认支付
- KOL尝试查看 pending_payment 状态的订单

