-- 先删除旧表（如果存在）
DROP TABLE IF EXISTS users;

-- 创建新表
CREATE TABLE users (
    id            INT          AUTO_INCREMENT PRIMARY KEY,
    username      VARCHAR(50)  NOT NULL UNIQUE,
    email         VARCHAR(100) NOT NULL UNIQUE,
    password      VARCHAR(100) NOT NULL,
    province      VARCHAR(50)  DEFAULT NULL,
    exam_type     VARCHAR(20)  DEFAULT NULL,
    school_type   VARCHAR(20)  DEFAULT NULL,
    score         INT         DEFAULT NULL,
    province_rank INT         DEFAULT NULL,
    physics       BOOLEAN     DEFAULT FALSE,
    history       BOOLEAN     DEFAULT FALSE,
    chemistry     BOOLEAN     DEFAULT FALSE,
    biology       BOOLEAN     DEFAULT FALSE,
    geography     BOOLEAN     DEFAULT FALSE,
    politics      BOOLEAN     DEFAULT FALSE,
    holland       VARCHAR(20) DEFAULT NULL,
    interests     TEXT        DEFAULT NULL,
    created_at    TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP   DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email),
    INDEX idx_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建测试用户
INSERT INTO users (username, email, password) VALUES 
('admin', 'admin@example.com', SHA2('admin123', 256)); 