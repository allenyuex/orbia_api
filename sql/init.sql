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
    current_team_id BIGINT COMMENT '当前团队ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_wallet_address (wallet_address),
    INDEX idx_email (email),
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