-- 创建数据库
CREATE DATABASE IF NOT EXISTS orbia CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE orbia;

-- 用户表
CREATE TABLE IF NOT EXISTS orbia_user (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    wallet_address VARCHAR(42) UNIQUE COMMENT '钱包地址',
    email VARCHAR(255) UNIQUE COMMENT '邮箱',
    password_hash VARCHAR(255) COMMENT '密码哈希',
    verification_code VARCHAR(10) COMMENT '验证码',
    code_expiry TIMESTAMP NULL COMMENT '验证码过期时间',
    nickname VARCHAR(100) COMMENT '昵称',
    avatar_url VARCHAR(500) COMMENT '头像URL',
    role ENUM('normal', 'admin') NOT NULL DEFAULT 'normal' COMMENT '用户角色：normal-普通用户，admin-管理员',
    status ENUM('normal', 'disabled', 'deleted') NOT NULL DEFAULT 'normal' COMMENT '用户状态：normal-正常，disabled-禁用，deleted-已删除',
    kol_id BIGINT COMMENT '关联的KOL ID',
    current_team_id BIGINT COMMENT '当前团队ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    INDEX idx_wallet_address (wallet_address),
    INDEX idx_email (email),
    INDEX idx_role (role),
    INDEX idx_status (status),
    INDEX idx_kol_id (kol_id),
    INDEX idx_current_team_id (current_team_id),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 团队表
CREATE TABLE IF NOT EXISTS orbia_team (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '团队ID',
    name VARCHAR(20) NOT NULL COMMENT '团队名称',
    icon_url VARCHAR(500) COMMENT '团队图标URL',
    creator_id BIGINT NOT NULL COMMENT '创建者ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    INDEX idx_creator_id (creator_id),
    INDEX idx_deleted_at (deleted_at),

    FOREIGN KEY (creator_id) REFERENCES orbia_user(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='团队表';

-- 团队成员表
CREATE TABLE IF NOT EXISTS orbia_team_member (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '成员ID',
    team_id BIGINT NOT NULL COMMENT '团队ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    role ENUM('creator', 'owner', 'member') NOT NULL DEFAULT 'member' COMMENT '角色：creator-创建者，owner-拥有者，member-成员',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
    UNIQUE KEY uk_team_user (team_id, user_id),
    INDEX idx_team_id (team_id),
    INDEX idx_user_id (user_id),
    INDEX idx_role (role),
    FOREIGN KEY (team_id) REFERENCES orbia_team(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES orbia_user(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='团队成员表';

-- 团队邀请表
CREATE TABLE IF NOT EXISTS orbia_team_invitation (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '邀请ID',
    team_id BIGINT NOT NULL COMMENT '团队ID',
    inviter_id BIGINT NOT NULL COMMENT '邀请者ID',
    invitee_email VARCHAR(255) COMMENT '被邀请者邮箱',
    invitee_wallet VARCHAR(42) COMMENT '被邀请者钱包地址',
    role ENUM('owner', 'member') NOT NULL DEFAULT 'member' COMMENT '邀请角色',
    status ENUM('pending', 'accepted', 'rejected', 'expired') NOT NULL DEFAULT 'pending' COMMENT '邀请状态',
    invitation_code VARCHAR(32) NOT NULL COMMENT '邀请码',
    expires_at TIMESTAMP NOT NULL COMMENT '过期时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    INDEX idx_team_id (team_id),
    INDEX idx_inviter_id (inviter_id),
    INDEX idx_invitee_email (invitee_email),
    INDEX idx_invitee_wallet (invitee_wallet),
    INDEX idx_invitation_code (invitation_code),
    INDEX idx_status (status),
    INDEX idx_expires_at (expires_at),
    INDEX idx_deleted_at (deleted_at),
    FOREIGN KEY (team_id) REFERENCES orbia_team(id) ON DELETE CASCADE,
    FOREIGN KEY (inviter_id) REFERENCES orbia_user(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='团队邀请表';

-- KOL信息表
CREATE TABLE IF NOT EXISTS orbia_kol (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT 'KOL ID',
    user_id BIGINT NOT NULL UNIQUE COMMENT '关联用户ID',
    avatar_url VARCHAR(500) COMMENT 'KOL头像URL',
    display_name VARCHAR(100) COMMENT 'KOL显示名称',
    description TEXT COMMENT 'KOL描述',
    country VARCHAR(50) COMMENT '所在国家',
    tiktok_url VARCHAR(500) COMMENT 'TikTok地址',
    youtube_url VARCHAR(500) COMMENT 'Youtube地址',
    x_url VARCHAR(500) COMMENT 'X(Twitter)地址',
    discord_url VARCHAR(500) COMMENT 'Discord地址',
    status ENUM('pending', 'approved', 'rejected') NOT NULL DEFAULT 'pending' COMMENT '审核状态：pending-待审核，approved-已通过，rejected-已拒绝',
    reject_reason TEXT COMMENT '拒绝原因',
    approved_at TIMESTAMP NULL COMMENT '审核通过时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at),
    FOREIGN KEY (user_id) REFERENCES orbia_user(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='KOL信息表';

-- KOL语言表
CREATE TABLE IF NOT EXISTS orbia_kol_language (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT 'ID',
    kol_id BIGINT NOT NULL COMMENT 'KOL ID',
    language_code VARCHAR(10) NOT NULL COMMENT '语言代码，如：en, zh, ja',
    language_name VARCHAR(50) NOT NULL COMMENT '语言名称，如：English, 中文, 日本語',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_kol_language (kol_id, language_code),
    INDEX idx_kol_id (kol_id),
    FOREIGN KEY (kol_id) REFERENCES orbia_kol(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='KOL语言表';

-- KOL标签表
CREATE TABLE IF NOT EXISTS orbia_kol_tag (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT 'ID',
    kol_id BIGINT NOT NULL COMMENT 'KOL ID',
    tag VARCHAR(50) NOT NULL COMMENT '标签，如：Defi, Web3, NFT',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_kol_tag (kol_id, tag),
    INDEX idx_kol_id (kol_id),
    INDEX idx_tag (tag),
    FOREIGN KEY (kol_id) REFERENCES orbia_kol(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='KOL标签表';

-- KOL数据统计表
CREATE TABLE IF NOT EXISTS orbia_kol_stats (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT 'ID',
    kol_id BIGINT NOT NULL UNIQUE COMMENT 'KOL ID',
    total_followers BIGINT COMMENT '全社交网站粉丝总数',
    tiktok_followers BIGINT COMMENT 'TikTok粉丝数',
    youtube_subscribers BIGINT COMMENT 'Youtube订阅数',
    x_followers BIGINT COMMENT 'X粉丝数',
    discord_members BIGINT COMMENT 'Discord成员数',
    tiktok_avg_views BIGINT COMMENT 'TikTok视频平均观看数',
    engagement_rate DECIMAL(10, 2) DEFAULT 0.00 COMMENT '订阅指数（Engagement Rate）',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_kol_id (kol_id),
    INDEX idx_total_followers (total_followers),
    INDEX idx_engagement_rate (engagement_rate),
    FOREIGN KEY (kol_id) REFERENCES orbia_kol(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='KOL数据统计表';

-- KOL报价Plans表
CREATE TABLE IF NOT EXISTS orbia_kol_plan (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT 'Plan ID',
    kol_id BIGINT NOT NULL COMMENT 'KOL ID',
    title VARCHAR(200) NOT NULL COMMENT '报价标题',
    description TEXT COMMENT '报价描述',
    price DECIMAL(10, 2) NOT NULL COMMENT '价格（美元）',
    plan_type ENUM('basic', 'standard', 'premium') NOT NULL COMMENT 'Plan类型：basic-基础，standard-标准，premium-高级',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    INDEX idx_kol_id (kol_id),
    INDEX idx_plan_type (plan_type),
    INDEX idx_deleted_at (deleted_at),
    FOREIGN KEY (kol_id) REFERENCES orbia_kol(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='KOL报价Plans表';

-- KOL视频表
DROP TABLE IF EXISTS orbia_kol_video;
CREATE TABLE orbia_kol_video (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '视频ID',
    kol_id BIGINT NOT NULL COMMENT 'KOL ID',
    embed_code TEXT NOT NULL COMMENT '视频嵌入代码',
    cover_url VARCHAR(500) COMMENT '视频封面URL',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    INDEX idx_kol_id (kol_id),
    INDEX idx_deleted_at (deleted_at),
    FOREIGN KEY (kol_id) REFERENCES orbia_kol(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='KOL视频表';

-- KOL订单表
DROP TABLE IF EXISTS orbia_kol_order;
CREATE TABLE orbia_kol_order (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID（内部使用）',
    order_id VARCHAR(64) NOT NULL UNIQUE COMMENT '订单ID（业务唯一ID，格式：KORD_{timestamp}_{random}）',
    user_id BIGINT NOT NULL COMMENT '下单用户ID',
    team_id BIGINT COMMENT '下单团队ID（如果是团队下单）',
    kol_id BIGINT NOT NULL COMMENT 'KOL ID',
    plan_id BIGINT NOT NULL COMMENT 'KOL报价Plan ID',
    plan_title VARCHAR(200) NOT NULL COMMENT 'Plan标题（快照）',
    plan_description TEXT COMMENT 'Plan描述（快照）',
    plan_price DECIMAL(10, 2) NOT NULL COMMENT 'Plan价格（快照，美元）',
    plan_type VARCHAR(20) NOT NULL COMMENT 'Plan类型（快照）：basic, standard, premium',
    title VARCHAR(200) NOT NULL COMMENT '订单标题',
    requirement_description TEXT NOT NULL COMMENT '合作需求描述',
    video_type VARCHAR(100) NOT NULL COMMENT '视频类型（用户手动输入）',
    video_duration INT NOT NULL COMMENT '视频预计时长（秒数）',
    target_audience VARCHAR(500) NOT NULL COMMENT '目标受众',
    expected_delivery_date DATE NOT NULL COMMENT '期望交付日期',
    additional_requirements TEXT COMMENT '额外要求',
    conversation_id VARCHAR(64) COMMENT '关联的会话ID（引用orbia_conversation.conversation_id）',
    status ENUM('pending_payment', 'pending', 'confirmed', 'in_progress', 'completed', 'cancelled', 'refunded') NOT NULL DEFAULT 'pending_payment' COMMENT '订单状态：pending_payment-待支付，pending-待确认，confirmed-已确认，in_progress-进行中，completed-已完成，cancelled-已取消，refunded-已退款',
    reject_reason TEXT COMMENT '拒绝/取消原因',
    confirmed_at TIMESTAMP NULL COMMENT '确认时间',
    completed_at TIMESTAMP NULL COMMENT '完成时间',
    cancelled_at TIMESTAMP NULL COMMENT '取消时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    INDEX idx_order_id (order_id),
    INDEX idx_user_id (user_id),
    INDEX idx_team_id (team_id),
    INDEX idx_kol_id (kol_id),
    INDEX idx_plan_id (plan_id),
    INDEX idx_conversation_id (conversation_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_title (title),
    INDEX idx_expected_delivery_date (expected_delivery_date),
    FOREIGN KEY (user_id) REFERENCES orbia_user(id) ON DELETE CASCADE,
    FOREIGN KEY (team_id) REFERENCES orbia_team(id) ON DELETE SET NULL,
    FOREIGN KEY (kol_id) REFERENCES orbia_kol(id) ON DELETE CASCADE,
    FOREIGN KEY (plan_id) REFERENCES orbia_kol_plan(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='KOL订单表';

-- 广告订单表
DROP TABLE IF EXISTS orbia_ad_order;
CREATE TABLE orbia_ad_order (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID（内部使用）',
    order_id VARCHAR(64) NOT NULL UNIQUE COMMENT '订单ID（业务唯一ID，格式：ADORD_{timestamp}_{random}）',
    user_id BIGINT NOT NULL COMMENT '下单用户ID',
    team_id BIGINT COMMENT '下单团队ID（如果是团队下单）',
    title VARCHAR(200) NOT NULL COMMENT '广告订单标题',
    description TEXT NOT NULL COMMENT '广告订单描述',
    budget DECIMAL(12, 2) NOT NULL COMMENT '广告预算（美元）',
    ad_type VARCHAR(50) NOT NULL COMMENT '广告类型：banner, video, social_media, influencer',
    target_audience VARCHAR(500) NOT NULL COMMENT '目标受众',
    start_date DATE NOT NULL COMMENT '开始日期',
    end_date DATE NOT NULL COMMENT '结束日期',
    status ENUM('pending', 'approved', 'in_progress', 'completed', 'cancelled') NOT NULL DEFAULT 'pending' COMMENT '订单状态：pending-待审核，approved-已批准，in_progress-进行中，completed-已完成，cancelled-已取消',
    reject_reason TEXT COMMENT '拒绝/取消原因',
    approved_at TIMESTAMP NULL COMMENT '批准时间',
    completed_at TIMESTAMP NULL COMMENT '完成时间',
    cancelled_at TIMESTAMP NULL COMMENT '取消时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    INDEX idx_order_id (order_id),
    INDEX idx_user_id (user_id),
    INDEX idx_team_id (team_id),
    INDEX idx_status (status),
    INDEX idx_ad_type (ad_type),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_title (title),
    INDEX idx_start_date (start_date),
    INDEX idx_end_date (end_date),
    FOREIGN KEY (user_id) REFERENCES orbia_user(id) ON DELETE CASCADE,
    FOREIGN KEY (team_id) REFERENCES orbia_team(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='广告订单表';

-- 用户钱包表
CREATE TABLE IF NOT EXISTS orbia_wallet (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '钱包ID',
    user_id BIGINT NOT NULL UNIQUE COMMENT '用户ID',
    balance DECIMAL(12, 2) NOT NULL DEFAULT 0.00 COMMENT '余额（美元）',
    frozen_balance DECIMAL(12, 2) NOT NULL DEFAULT 0.00 COMMENT '冻结余额（美元）',
    total_recharge DECIMAL(12, 2) NOT NULL DEFAULT 0.00 COMMENT '累计充值金额（美元）',
    total_consume DECIMAL(12, 2) NOT NULL DEFAULT 0.00 COMMENT '累计消费金额（美元）',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_id (user_id),
    INDEX idx_balance (balance),
    FOREIGN KEY (user_id) REFERENCES orbia_user(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户钱包表';

-- 交易记录表（仅用于记录支出账单）
DROP TABLE IF EXISTS orbia_transaction;
CREATE TABLE orbia_transaction (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID（内部使用）',
    transaction_id VARCHAR(64) NOT NULL UNIQUE COMMENT '交易ID（业务唯一ID，格式：TXN{snowflake_id}）',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    type ENUM('consume', 'refund', 'freeze', 'unfreeze') NOT NULL COMMENT '交易类型：consume-消费，refund-退款，freeze-冻结，unfreeze-解冻',
    amount DECIMAL(12, 2) NOT NULL COMMENT '交易金额（美元）',
    balance_before DECIMAL(12, 2) NOT NULL COMMENT '交易前余额（美元）',
    balance_after DECIMAL(12, 2) NOT NULL COMMENT '交易后余额（美元）',
    status ENUM('pending', 'processing', 'completed', 'failed', 'cancelled') NOT NULL DEFAULT 'pending' COMMENT '交易状态：pending-待处理，processing-处理中，completed-已完成，failed-失败，cancelled-已取消',
    related_order_type VARCHAR(50) COMMENT '关联订单类型：kol_order-KOL订单，ad_order-广告订单',
    related_order_id VARCHAR(64) COMMENT '关联订单ID（如果是消费/退款类型）',
    remark TEXT COMMENT '备注说明',
    completed_at TIMESTAMP NULL COMMENT '完成时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_transaction_id (transaction_id),
    INDEX idx_user_id (user_id),
    INDEX idx_type (type),
    INDEX idx_status (status),
    INDEX idx_related_order_type (related_order_type),
    INDEX idx_related_order_id (related_order_id),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (user_id) REFERENCES orbia_user(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='交易记录表（仅记录支出账单）';


-- 创建数据字典表
CREATE TABLE orbia_dictionary (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '字典ID',
    code VARCHAR(100) NOT NULL UNIQUE COMMENT '字典编码（唯一，只能大小写字母）',
    name VARCHAR(100) NOT NULL COMMENT '字典名称',
    description TEXT COMMENT '字典描述',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    INDEX idx_code (code),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据字典表';

-- 创建数据字典项表
CREATE TABLE orbia_dictionary_item (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '字典项ID',
    dictionary_id BIGINT NOT NULL COMMENT '字典ID',
    parent_id BIGINT NOT NULL DEFAULT 0 COMMENT '父级ID（0表示根节点）',
    code VARCHAR(100) NOT NULL COMMENT '字典项编码',
    name VARCHAR(200) NOT NULL COMMENT '字典项名称',
    description TEXT COMMENT '字典项描述',
    icon_url VARCHAR(500) COMMENT '图标URL',
    sort_order INT NOT NULL DEFAULT 0 COMMENT '排序序号（升序）',
    level INT NOT NULL DEFAULT 1 COMMENT '层级（1开始）',
    path VARCHAR(1000) NOT NULL COMMENT '路径（如: 1/2/3）',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    UNIQUE KEY uk_dict_parent_code (dictionary_id, parent_id, code),
    INDEX idx_dictionary_id (dictionary_id),
    INDEX idx_parent_id (parent_id),
    INDEX idx_code (code),
    INDEX idx_status (status),
    INDEX idx_sort_order (sort_order),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_path (path(255)),
    FOREIGN KEY (dictionary_id) REFERENCES orbia_dictionary(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据字典项表';

-- 充值订单表（需要先删除，因为有外键约束）
DROP TABLE IF EXISTS orbia_recharge_order;

-- 收款钱包设置表
DROP TABLE IF EXISTS orbia_payment_setting;
CREATE TABLE orbia_payment_setting (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '设置ID',
    network VARCHAR(100) NOT NULL COMMENT '区块链网络（如：TRC-20 - TRON Network (TRC-20)）',
    address VARCHAR(500) NOT NULL COMMENT '钱包地址',
    label VARCHAR(200) NOT NULL COMMENT '钱包标签（如：USDT-TRC20 主钱包）',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    INDEX idx_network (network),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='收款钱包设置表';

-- 创建充值订单表
CREATE TABLE orbia_recharge_order (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID（内部使用）',
    order_id VARCHAR(64) NOT NULL UNIQUE COMMENT '订单ID（业务唯一ID，格式：RCHORD_{timestamp}_{random}）',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    amount DECIMAL(12, 2) NOT NULL COMMENT '充值金额（美元）',
    payment_type ENUM('crypto', 'online') NOT NULL COMMENT '支付类型：crypto-加密货币，online-在线支付',
    payment_setting_id BIGINT COMMENT '关联的payment_setting ID（用于加密货币支付）',
    payment_network VARCHAR(100) COMMENT '快照-区块链网络',
    payment_address VARCHAR(500) COMMENT '快照-钱包地址',
    payment_label VARCHAR(200) COMMENT '快照-钱包标签',
    user_crypto_address VARCHAR(500) COMMENT '用户的转账钱包地址（仅加密货币支付）',
    crypto_tx_hash VARCHAR(500) COMMENT '加密货币交易哈希（用户或管理员填写）',
    online_payment_platform VARCHAR(50) COMMENT '在线支付平台：stripe, paypal',
    online_payment_order_id VARCHAR(200) COMMENT '在线支付平台订单ID',
    online_payment_url TEXT COMMENT '在线支付URL',
    status ENUM('pending', 'confirmed', 'failed', 'cancelled') NOT NULL DEFAULT 'pending' COMMENT '订单状态：pending-待确认，confirmed-已确认，failed-失败，cancelled-已取消',
    confirmed_by BIGINT COMMENT '确认人ID（管理员）',
    confirmed_at TIMESTAMP NULL COMMENT '确认时间',
    failed_reason TEXT COMMENT '失败原因',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    INDEX idx_order_id (order_id),
    INDEX idx_user_id (user_id),
    INDEX idx_payment_type (payment_type),
    INDEX idx_payment_setting_id (payment_setting_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_confirmed_at (confirmed_at),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_crypto_tx_hash (crypto_tx_hash),
    FOREIGN KEY (user_id) REFERENCES orbia_user(id) ON DELETE CASCADE,
    FOREIGN KEY (payment_setting_id) REFERENCES orbia_payment_setting(id) ON DELETE SET NULL,
    FOREIGN KEY (confirmed_by) REFERENCES orbia_user(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='充值订单表';

-- 验证码表
DROP TABLE IF EXISTS orbia_verification_code;
CREATE TABLE orbia_verification_code (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID',
    email VARCHAR(255) NOT NULL COMMENT '邮箱地址',
    code VARCHAR(10) NOT NULL COMMENT '验证码',
    code_type ENUM('login', 'register', 'reset_password') NOT NULL DEFAULT 'login' COMMENT '验证码类型：login-登录，register-注册，reset_password-重置密码',
    status ENUM('unused', 'used', 'expired') NOT NULL DEFAULT 'unused' COMMENT '状态：unused-未使用，used-已使用，expired-已过期',
    used_at TIMESTAMP NULL COMMENT '使用时间',
    expires_at TIMESTAMP NOT NULL COMMENT '过期时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_email (email),
    INDEX idx_code (code),
    INDEX idx_code_type (code_type),
    INDEX idx_status (status),
    INDEX idx_expires_at (expires_at),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='验证码表';

-- 会话表
DROP TABLE IF EXISTS orbia_message;
DROP TABLE IF EXISTS orbia_conversation_member;
DROP TABLE IF EXISTS orbia_conversation;

CREATE TABLE orbia_conversation (
    conversation_id VARCHAR(64) PRIMARY KEY COMMENT '会话ID（业务唯一ID，格式：CONV_{timestamp}_{random}）',
    title VARCHAR(200) COMMENT '会话标题',
    type ENUM('kol_order', 'ad_order', 'general', 'support') NOT NULL DEFAULT 'general' COMMENT '会话类型：kol_order-KOL订单会话，ad_order-广告订单会话，general-普通会话，support-客服会话',
    related_order_type VARCHAR(50) COMMENT '关联订单类型：kol_order, ad_order',
    related_order_id VARCHAR(64) COMMENT '关联订单ID',
    status ENUM('active', 'archived', 'closed') NOT NULL DEFAULT 'active' COMMENT '会话状态：active-活跃，archived-已归档，closed-已关闭',
    last_message_at TIMESTAMP NULL COMMENT '最后消息时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    INDEX idx_type (type),
    INDEX idx_related_order_type (related_order_type),
    INDEX idx_related_order_id (related_order_id),
    INDEX idx_status (status),
    INDEX idx_last_message_at (last_message_at),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会话表';

-- 会话成员表
CREATE TABLE orbia_conversation_member (
    conversation_id VARCHAR(64) NOT NULL COMMENT '会话ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    role ENUM('creator', 'member', 'admin') NOT NULL DEFAULT 'member' COMMENT '成员角色：creator-创建者，member-成员，admin-管理员',
    unread_count INT NOT NULL DEFAULT 0 COMMENT '未读消息数',
    last_read_at TIMESTAMP NULL COMMENT '最后已读时间',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
    PRIMARY KEY (conversation_id, user_id),
    INDEX idx_user_id (user_id),
    INDEX idx_role (role),
    FOREIGN KEY (conversation_id) REFERENCES orbia_conversation(conversation_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES orbia_user(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会话成员表';

-- 消息表
CREATE TABLE orbia_message (
    message_id VARCHAR(64) PRIMARY KEY COMMENT '消息ID（业务唯一ID，格式：MSG_{timestamp}_{random}）',
    conversation_id VARCHAR(64) NOT NULL COMMENT '会话ID',
    sender_id BIGINT NOT NULL COMMENT '发送者用户ID',
    message_type ENUM('text', 'image', 'file', 'video', 'audio', 'system') NOT NULL DEFAULT 'text' COMMENT '消息类型：text-文本，image-图片，file-文件，video-视频，audio-音频，system-系统消息',
    content TEXT NOT NULL COMMENT '消息内容（文本内容或文件URL）',
    file_name VARCHAR(500) COMMENT '文件名（如果是文件类型）',
    file_size BIGINT COMMENT '文件大小（字节）',
    file_type VARCHAR(100) COMMENT '文件MIME类型',
    status ENUM('sent', 'delivered', 'read', 'failed') NOT NULL DEFAULT 'sent' COMMENT '消息状态：sent-已发送，delivered-已送达，read-已读，failed-发送失败',
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间（毫秒精度）',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    INDEX idx_conversation_id (conversation_id),
    INDEX idx_sender_id (sender_id),
    INDEX idx_message_type (message_type),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at),
    FOREIGN KEY (conversation_id) REFERENCES orbia_conversation(conversation_id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES orbia_user(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='消息表';

-- Campaign表（广告活动表）
CREATE TABLE IF NOT EXISTS orbia_campaign (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID',
    campaign_id VARCHAR(64) NOT NULL UNIQUE COMMENT '业务唯一ID（格式：CAMPAIGN_{timestamp}_{random}）',
    user_id BIGINT NOT NULL COMMENT '创建用户ID',
    team_id BIGINT NOT NULL COMMENT '所属团队ID',
    campaign_name VARCHAR(200) NOT NULL COMMENT '活动名称',
    promotion_objective ENUM('awareness', 'consideration', 'conversion') NOT NULL COMMENT '推广目标：awareness-品牌认知，consideration-受众意向，conversion-行为转化',
    optimization_goal VARCHAR(50) NOT NULL COMMENT '优化目标：根据promotion_objective不同有不同值',
    location TEXT COMMENT '地区（JSON数组，存储数据字典ID列表）',
    age BIGINT COMMENT '年龄段（引用数据字典ID）',
    gender BIGINT COMMENT '性别（引用数据字典ID）',
    languages TEXT COMMENT '语言（JSON数组，多选数据字典ID）',
    spending_power BIGINT COMMENT '消费能力（引用数据字典ID）',
    operating_system BIGINT COMMENT '操作系统（引用数据字典ID）',
    os_versions TEXT COMMENT '系统版本（JSON数组，多选数据字典ID）',
    device_models TEXT COMMENT '设备品牌（JSON数组，多选数据字典ID）',
    connection_types TEXT COMMENT '网络情况（JSON数组，多选数据字典ID）',
    device_price_type TINYINT DEFAULT 0 COMMENT '设备价格类型：0-any，1-specific range',
    device_price_min DECIMAL(15,2) COMMENT '设备价格最小值',
    device_price_max DECIMAL(15,2) COMMENT '设备价格最大值',
    planned_start_time TIMESTAMP NOT NULL COMMENT '计划开始时间',
    planned_end_time TIMESTAMP NOT NULL COMMENT '计划结束时间',
    time_zone BIGINT COMMENT '时区（引用数据字典ID）',
    dayparting_type TINYINT DEFAULT 0 COMMENT '分时段类型：0-全天，1-特定时段',
    dayparting_schedule TEXT COMMENT '特定时段配置（JSON格式）',
    frequency_cap_type TINYINT DEFAULT 0 COMMENT '频次上限类型：0-每七天不超过三次，1-每天不超过一次，2-自定义',
    frequency_cap_times INT COMMENT '自定义频次（次数）',
    frequency_cap_days INT COMMENT '自定义频次（天数）',
    budget_type TINYINT NOT NULL COMMENT '预算类型：0-每日预算，1-总预算',
    budget_amount DECIMAL(15,2) NOT NULL COMMENT '预算金额',
    website VARCHAR(1000) COMMENT '网站链接',
    ios_download_url VARCHAR(1000) COMMENT 'iOS下载链接',
    android_download_url VARCHAR(1000) COMMENT 'Android下载链接',
    status ENUM('pending', 'active', 'paused', 'ended') NOT NULL DEFAULT 'pending' COMMENT '状态：pending-待启动，active-已启动，paused-暂停，ended-已结束',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    INDEX idx_campaign_id (campaign_id),
    INDEX idx_user_id (user_id),
    INDEX idx_team_id (team_id),
    INDEX idx_status (status),
    INDEX idx_promotion_objective (promotion_objective),
    INDEX idx_planned_start_time (planned_start_time),
    INDEX idx_planned_end_time (planned_end_time),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at),
    FOREIGN KEY (user_id) REFERENCES orbia_user(id) ON DELETE CASCADE,
    FOREIGN KEY (team_id) REFERENCES orbia_team(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Campaign表（广告活动表）';

-- Campaign附件表
CREATE TABLE IF NOT EXISTS orbia_campaign_attachment (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '附件ID',
    campaign_id BIGINT NOT NULL COMMENT '关联Campaign ID',
    file_url VARCHAR(1000) NOT NULL COMMENT '文件URL',
    file_name VARCHAR(500) NOT NULL COMMENT '文件名',
    file_type VARCHAR(100) NOT NULL COMMENT '文件类型（MIME类型）',
    file_size BIGINT COMMENT '文件大小（字节）',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    INDEX idx_campaign_id (campaign_id),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at),
    FOREIGN KEY (campaign_id) REFERENCES orbia_campaign(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Campaign附件表';

-- 优秀广告案例表
CREATE TABLE IF NOT EXISTS orbia_excellent_case (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '案例ID',
    video_url VARCHAR(1000) NOT NULL COMMENT '视频URL',
    cover_url VARCHAR(1000) NOT NULL COMMENT '封面URL',
    title VARCHAR(200) NOT NULL COMMENT '案例标题',
    description TEXT COMMENT '案例描述',
    sort_order INT NOT NULL DEFAULT 0 COMMENT '排序序号（升序，数字越小越靠前）',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    INDEX idx_sort_order (sort_order),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='优秀广告案例表';

-- 内容趋势表
CREATE TABLE IF NOT EXISTS orbia_content_trend (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '趋势ID',
    ranking INT NOT NULL COMMENT '排名（1,2,3,4,5...）',
    hot_keyword VARCHAR(200) NOT NULL COMMENT '热点词（如：defi, web3等）',
    value_level ENUM('low', 'medium', 'high') NOT NULL COMMENT '价值等级：low-低，medium-中，high-高',
    heat BIGINT NOT NULL DEFAULT 0 COMMENT '热度值',
    growth_rate DECIMAL(10, 2) NOT NULL DEFAULT 0.00 COMMENT '增长比例（百分比，如：15.5表示15.5%）',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    UNIQUE KEY uk_ranking (ranking, deleted_at),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_created_at (created_at),
    INDEX idx_ranking (ranking)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='内容趋势表';

-- 平台数据统计表
CREATE TABLE IF NOT EXISTS orbia_platform_stats (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '统计ID（此表只有一行数据）',
    active_kols BIGINT NOT NULL DEFAULT 0 COMMENT '活跃的KOLs数量',
    total_coverage BIGINT NOT NULL DEFAULT 0 COMMENT '总覆盖用户数',
    total_ad_impressions BIGINT NOT NULL DEFAULT 0 COMMENT '累计广告曝光次数',
    total_transaction_amount DECIMAL(15, 2) NOT NULL DEFAULT 0.00 COMMENT '平台总交易额（美元）',
    average_roi DECIMAL(10, 2) NOT NULL DEFAULT 0.00 COMMENT '平均ROI（百分比，如：15.5表示15.5%）',
    average_cpm DECIMAL(10, 2) NOT NULL DEFAULT 0.00 COMMENT '平均CPM（Cost Per Mille，每千次展示成本）',
    web3_brand_count BIGINT NOT NULL DEFAULT 0 COMMENT '合作Web3品牌数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_created_at (created_at),
    INDEX idx_updated_at (updated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='平台数据统计表';