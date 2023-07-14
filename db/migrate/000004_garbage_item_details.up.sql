CREATE TABLE IF NOT EXISTS garbage_item_details (
  id INT AUTO_INCREMENT PRIMARY KEY,
  garbage_id VARCHAR(255),
  garbage_item_id INT,
  language_code VARCHAR(255),
  translated_name VARCHAR(255),
  translated_category VARCHAR(255),
  translated_description TEXT,
  FOREIGN KEY (garbage_item_id) REFERENCES garbage_items(id)
);