CREATE TABLE IF NOT EXISTS garbage_masters (
    id INT AUTO_INCREMENT PRIMARY KEY,
    public_code VARCHAR(255),
    garbage_id VARCHAR(255) UNIQUE,
    public_name VARCHAR(255),
    district VARCHAR(255),
    item VARCHAR(255),
    item_kana VARCHAR(255),
    item_eng VARCHAR(255),
    classify VARCHAR(255), -- TODO: categoryのほうが自然 可燃ごみ　不燃ごみ　資源　粗大ごみ　その他　で大別できる
    note VARCHAR(255), -- TODO: cautionのほうがいいかも　注意点: 現状使われていない
    remarks VARCHAR(255), -- 備考欄　注意点: 頻出
    large_fee VARCHAR(255), -- 粗大ごみ回収料金　注意点: 現状使われていない
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
