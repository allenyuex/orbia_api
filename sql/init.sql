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
    role ENUM('user', 'admin') NOT NULL DEFAULT 'user' COMMENT '用户角色：user-普通用户，admin-管理员',
    status ENUM('normal', 'disabled', 'deleted') NOT NULL DEFAULT 'normal' COMMENT '用户状态：normal-正常，disabled-禁用，deleted-已删除',
    kol_id BIGINT COMMENT '关联的KOL ID',
    current_team_id BIGINT COMMENT '当前团队ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_wallet_address (wallet_address),
    INDEX idx_email (email),
    INDEX idx_role (role),
    INDEX idx_status (status),
    INDEX idx_kol_id (kol_id),
    INDEX idx_current_team_id (current_team_id),
    INDEX idx_created_at (created_at)
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
    total_followers BIGINT DEFAULT 0 COMMENT '全社交网站粉丝总数',
    tiktok_followers BIGINT DEFAULT 0 COMMENT 'TikTok粉丝数',
    youtube_subscribers BIGINT DEFAULT 0 COMMENT 'Youtube订阅数',
    x_followers BIGINT DEFAULT 0 COMMENT 'X粉丝数',
    discord_members BIGINT DEFAULT 0 COMMENT 'Discord成员数',
    tiktok_avg_views BIGINT DEFAULT 0 COMMENT 'TikTok视频平均观看数',
    engagement_rate DECIMAL(10, 2) DEFAULT 0 COMMENT '订阅指数（Engagement Rate）',
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
    status ENUM('pending', 'confirmed', 'in_progress', 'completed', 'cancelled', 'refunded') NOT NULL DEFAULT 'pending' COMMENT '订单状态：pending-待确认，confirmed-已确认，in_progress-进行中，completed-已完成，cancelled-已取消，refunded-已退款',
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

-- 交易记录表
CREATE TABLE IF NOT EXISTS orbia_transaction (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID（内部使用）',
    transaction_id VARCHAR(64) NOT NULL UNIQUE COMMENT '交易ID（业务唯一ID，格式：TXN{snowflake_id}）',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    type ENUM('recharge', 'consume', 'refund', 'freeze', 'unfreeze') NOT NULL COMMENT '交易类型：recharge-充值，consume-消费，refund-退款，freeze-冻结，unfreeze-解冻',
    amount DECIMAL(12, 2) NOT NULL COMMENT '交易金额（美元）',
    balance_before DECIMAL(12, 2) NOT NULL COMMENT '交易前余额（美元）',
    balance_after DECIMAL(12, 2) NOT NULL COMMENT '交易后余额（美元）',
    status ENUM('pending', 'processing', 'completed', 'failed', 'cancelled') NOT NULL DEFAULT 'pending' COMMENT '交易状态：pending-待处理，processing-处理中，completed-已完成，failed-失败，cancelled-已取消',
    payment_method ENUM('crypto', 'online') COMMENT '支付方式：crypto-加密货币，online-在线支付',
    crypto_currency VARCHAR(20) COMMENT '加密货币类型：USDT, USDC',
    crypto_chain VARCHAR(50) COMMENT '加密货币链：ETH, BSC, POLYGON, TRON, ARBITRUM, OPTIMISM',
    crypto_address VARCHAR(500) COMMENT '加密货币支付地址',
    crypto_tx_hash VARCHAR(500) COMMENT '加密货币交易哈希',
    online_payment_platform VARCHAR(50) COMMENT '在线支付平台：stripe, paypal',
    online_payment_order_id VARCHAR(200) COMMENT '在线支付平台订单ID',
    online_payment_url TEXT COMMENT '在线支付URL',
    related_order_id VARCHAR(64) COMMENT '关联订单ID（如果是消费/退款类型）',
    remark TEXT COMMENT '备注说明',
    failed_reason TEXT COMMENT '失败原因',
    completed_at TIMESTAMP NULL COMMENT '完成时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_transaction_id (transaction_id),
    INDEX idx_user_id (user_id),
    INDEX idx_type (type),
    INDEX idx_status (status),
    INDEX idx_payment_method (payment_method),
    INDEX idx_related_order_id (related_order_id),
    INDEX idx_created_at (created_at),
    INDEX idx_crypto_tx_hash (crypto_tx_hash),
    FOREIGN KEY (user_id) REFERENCES orbia_user(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='交易记录表';