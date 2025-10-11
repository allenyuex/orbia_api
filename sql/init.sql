-- 创建数据库
CREATE DATABASE IF NOT EXISTS orbia_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE orbia_db;

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
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_wallet_address (wallet_address),
    INDEX idx_email (email),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';