-- Migration: 001_create_speed_tests_table.up.sql
-- Create speed tests table to store test results

CREATE TABLE IF NOT EXISTS speed_tests (
    id VARCHAR(36) PRIMARY KEY,
    client_ip VARCHAR(45) NOT NULL,
    user_agent TEXT,
    test_type ENUM('download', 'upload', 'ping', 'full') NOT NULL,
    
    -- Test results
    download_speed_mbps DECIMAL(10,2) DEFAULT NULL,
    upload_speed_mbps DECIMAL(10,2) DEFAULT NULL,
    ping_latency_ms DECIMAL(8,2) DEFAULT NULL,
    jitter_ms DECIMAL(8,2) DEFAULT NULL,
    
    -- Test parameters
    download_size_bytes BIGINT DEFAULT NULL,
    upload_size_bytes BIGINT DEFAULT NULL,
    test_duration_seconds DECIMAL(8,2) DEFAULT NULL,
    
    -- Client information
    isp VARCHAR(255) DEFAULT NULL,
    country VARCHAR(2) DEFAULT NULL,
    region VARCHAR(255) DEFAULT NULL,
    city VARCHAR(255) DEFAULT NULL,
    
    -- Ookla-compatible fields
    server_name VARCHAR(255) DEFAULT 'Krea Speed Test Server',
    server_country VARCHAR(2) DEFAULT 'IN',
    server_city VARCHAR(255) DEFAULT 'Sri City',
    sponsor VARCHAR(255) DEFAULT 'Krea University',
    
    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_client_ip (client_ip),
    INDEX idx_test_type (test_type),
    INDEX idx_created_at (created_at)
);

-- Create API keys table for authentication
CREATE TABLE IF NOT EXISTS api_keys (
    id VARCHAR(36) PRIMARY KEY,
    key_hash VARCHAR(64) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    rate_limit_per_minute INT DEFAULT 60,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP NULL,
    
    INDEX idx_key_hash (key_hash),
    INDEX idx_is_active (is_active)
);

-- Create rate limiting table
CREATE TABLE IF NOT EXISTS rate_limits (
    id VARCHAR(255) PRIMARY KEY,
    identifier VARCHAR(255) NOT NULL, -- IP address or API key
    endpoint VARCHAR(255) NOT NULL,
    request_count INT DEFAULT 1,
    window_start TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE KEY unique_identifier_endpoint_window (identifier, endpoint, window_start),
    INDEX idx_identifier (identifier),
    INDEX idx_window_start (window_start)
);

-- Create whitelist table for rate limit exemptions
CREATE TABLE IF NOT EXISTS rate_limit_whitelist (
    id VARCHAR(36) PRIMARY KEY,
    ip_address VARCHAR(45),
    ip_range VARCHAR(50), -- CIDR notation
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    
    INDEX idx_ip_address (ip_address),
    INDEX idx_is_active (is_active)
);

-- Insert default API key (for demo purposes)
INSERT INTO api_keys (id, key_hash, name, description, rate_limit_per_minute) 
VALUES (
    UUID(), 
    SHA2('demo-api-key-2025', 256), 
    'Demo API Key', 
    'Default API key for demonstration', 
    1000
) ON DUPLICATE KEY UPDATE id=id;

-- Insert default whitelist entries
INSERT INTO rate_limit_whitelist (id, ip_address, description) VALUES 
(UUID(), '127.0.0.1', 'Localhost'),
(UUID(), '::1', 'IPv6 localhost')
ON DUPLICATE KEY UPDATE id=id;
