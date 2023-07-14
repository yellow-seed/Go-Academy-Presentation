CREATE TABLE IF NOT EXISTS garbage_masters (
    id INT AUTO_INCREMENT PRIMARY KEY,
    public_code VARCHAR(255),
    garbage_id VARCHAR(255) UNIQUE,
    public_name VARCHAR(255),
    district VARCHAR(255),
    item VARCHAR(255),
    item_kana VARCHAR(255),
    item_eng VARCHAR(255),
    classify VARCHAR(255),
    note VARCHAR(255),
    remarks VARCHAR(255),
    large_fee VARCHAR(255),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);
