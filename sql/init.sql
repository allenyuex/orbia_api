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
    current_team_id BIGINT COMMENT '当前团队ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_wallet_address (wallet_address),
    INDEX idx_email (email),
    INDEX idx_role (role),
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
CREATE TABLE IF NOT EXISTS orbia_kol_video (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '视频ID',
    kol_id BIGINT NOT NULL COMMENT 'KOL ID',
    title VARCHAR(500) NOT NULL COMMENT '视频标题',
    content TEXT COMMENT '视频内容/描述',
    cover_url VARCHAR(500) COMMENT '视频封面URL',
    video_url VARCHAR(500) COMMENT '视频链接',
    platform VARCHAR(50) NOT NULL COMMENT '平台：tiktok, youtube',
    platform_video_id VARCHAR(200) COMMENT '平台视频ID',
    likes_count BIGINT DEFAULT 0 COMMENT '点赞数',
    views_count BIGINT DEFAULT 0 COMMENT '观看数',
    comments_count BIGINT DEFAULT 0 COMMENT '评论数',
    shares_count BIGINT DEFAULT 0 COMMENT '分享数',
    published_at TIMESTAMP NULL COMMENT '发布时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    INDEX idx_kol_id (kol_id),
    INDEX idx_platform (platform),
    INDEX idx_published_at (published_at),
    INDEX idx_deleted_at (deleted_at),
    FOREIGN KEY (kol_id) REFERENCES orbia_kol(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='KOL视频表';