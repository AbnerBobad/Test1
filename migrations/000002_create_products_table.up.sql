CREATE TABLE IF NOT EXISTS products (
  product_id BigSerial PRIMARY KEY,
  product_name varchar(100) NOT NULL,
  product_quantity int NOT NULL DEFAULT 0,
  product_price decimal(10,2) NOT NULL,
  product_description text NOT NULL,
  category_id int,
  added_by int,
  updated_by int,
  created_at timestamp DEFAULT now(),
  updated_at timestamp DEFAULT now()
);

ALTER TABLE products ADD FOREIGN KEY (category_id) REFERENCES categories (category_id);
ALTER TABLE products ADD FOREIGN KEY (added_by) REFERENCES users (user_id);
ALTER TABLE products ADD FOREIGN KEY (updated_by) REFERENCES users (user_id);
