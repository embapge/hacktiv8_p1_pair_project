DROP DATABASE IF EXISTS pair_project; 
 
CREATE DATABASE pair_project; 
 
USE pair_project; 
 
CREATE TABLE users ( 
    id INT PRIMARY KEY AUTO_INCREMENT, 
    username VARCHAR(100) NOT NULL UNIQUE, 
    email VARCHAR(100) NOT NULL UNIQUE, 
    password VARCHAR(255) NOT NULL, 
    role ENUM('admin', 'staff', 'customer') NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_role (role), 
    INDEX idx_created_at (created_at) 
); 
 
CREATE TABLE customers ( 
    id INT PRIMARY KEY AUTO_INCREMENT, 
    name VARCHAR(100) NOT NULL, 
    address TEXT NOT NULL, 
    email VARCHAR(100) NOT NULL UNIQUE, 
    phone_number VARCHAR(20) NOT NULL UNIQUE, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by INT, 
    updated_by INT, 
    FOREIGN KEY (created_by) REFERENCES users(id), 
    FOREIGN KEY (updated_by) REFERENCES users(id), 
    INDEX idx_name (name), 
    INDEX idx_email (email) 
); 
 
CREATE TABLE user_customers ( 
    user_id INT NOT NULL, 
    customer_id INT NOT NULL PRIMARY KEY, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id), 
    FOREIGN KEY (customer_id) REFERENCES customers(id) 
); 
 
CREATE TABLE categories ( 
    id INT PRIMARY KEY AUTO_INCREMENT, 
    name VARCHAR(100) NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by INT NOT NULL, 
    updated_by INT,  
    FOREIGN KEY (created_by) REFERENCES users(id), 
    FOREIGN KEY (updated_by) REFERENCES users(id), 
    UNIQUE INDEX idx_name (name) 
); 
 
CREATE TABLE products ( 
    id INT PRIMARY KEY AUTO_INCREMENT, 
    name VARCHAR(100) NOT NULL, 
    stock INT DEFAULT 0 NOT NULL, 
    description TEXT, 
    category_id INT NOT NULL, 
    price DECIMAL(10,2) DEFAULT 0 NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, 
    created_by INT NOT NULL, 
    updated_by INT, 
    FOREIGN KEY (created_by) REFERENCES users(id), 
    FOREIGN KEY (updated_by) REFERENCES users(id), 
    FOREIGN KEY (category_id) REFERENCES categories(id), 
    INDEX idx_category_id (category_id), 
    INDEX idx_price (price), 
    INDEX idx_name (name) 
); 
 
CREATE TABLE orders ( 
    id INT PRIMARY KEY AUTO_INCREMENT, 
    customer_id INT NOT NULL, 
    number_display VARCHAR(50) UNIQUE, 
    date DATE NOT NULL, 
    status ENUM('processing', 'completed', 'cancel') NOT NULL DEFAULT 'processing',
    total DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (total >= 0), 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by INT NOT NULL,
    updated_by INT, 
    FOREIGN KEY (customer_id) REFERENCES customers(id),
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id),
    INDEX idx_date (date), 
    INDEX idx_status (status)
);
 
CREATE TABLE order_details ( 
    id INT PRIMARY KEY AUTO_INCREMENT, 
    order_id INT NOT NULL, 
    product_id INT NOT NULL,
    qty INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by INT NOT NULL, 
    updated_by INT, 
    FOREIGN KEY (created_by) REFERENCES users(id), 
    FOREIGN KEY (updated_by) REFERENCES users(id), 
    FOREIGN KEY (order_id) REFERENCES orders(id), 
    FOREIGN KEY (product_id) REFERENCES products(id) 
); 
 
CREATE TABLE billings ( 
    id INT PRIMARY KEY AUTO_INCREMENT, 
    order_id INT NOT NULL, 
    number_display VARCHAR(50) UNIQUE, 
    tax DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (tax >= 0), 
    total DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (total >= 0), 
    status ENUM('unpaid', 'paid', 'cancelled', 'refunded') NOT NULL DEFAULT 'unpaid', 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by INT NOT NULL,
    updated_by INT, 
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id),
    INDEX idx_status (status)
);
 
CREATE TABLE payments ( 
    id INT PRIMARY KEY AUTO_INCREMENT, 
    billing_id INT NOT NULL, 
    date DATETIME NOT NULL, 
    amount DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (amount >= 0), 
    method ENUM('credit_card', 'va', 'transfer') NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by INT NOT NULL,
    updated_by INT,
    FOREIGN KEY (billing_id) REFERENCES billings(id),
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id),
    INDEX idx_date (date), 
    INDEX idx_method (method)
);


INSERT INTO users (username, email, password, role)
VALUES 
/*('admin01', 'admin01@example.com', '$2b$12$rcVbOnmFPSu5S4sSccrYPuLZXHybabFFYCFIi9R4uEft1uTeq2rO2', 'admin'),
('staff01', 'staff01@example.com', '$2b$12$rcVbOnmFPSu5S4sSccrYPuLZXHybabFFYCFIi9R4uEft1uTeq2rO2', 'staff'),
('staff02', 'staff02@example.com', '$2b$12$rcVbOnmFPSu5S4sSccrYPuLZXHybabFFYCFIi9R4uEft1uTeq2rO2', 'staff'),
('custuser1', 'cust1@example.com', '$2b$12$rcVbOnmFPSu5S4sSccrYPuLZXHybabFFYCFIi9R4uEft1uTeq2rO2', 'customer'),
('custuser2', 'cust2@example.com', '$2b$12$rcVbOnmFPSu5S4sSccrYPuLZXHybabFFYCFIi9R4uEft1uTeq2rO2', 'customer');*/

('admin01', 'admin01@example.com', '123456', 'admin'),
('staff01', 'staff01@example.com', '123456', 'staff'),
('staff02', 'staff02@example.com', '123456', 'staff'),
('custuser1', 'cust1@example.com', '123456', 'customer'),
('custuser2', 'cust2@example.com', '123456', 'customer');


INSERT INTO customers (name, address, email, phone_number, created_by, updated_by)
VALUES 
('John Doe', 'Jl. Merdeka No. 123, Jakarta', 'john@example.com', '081234567890', 2, 2),
('Jane Smith', 'Jl. Sudirman No. 10, Bandung', 'jane@example.com', '089876543210', 3, 3);

INSERT INTO user_customers (user_id, customer_id)
VALUES 
(4, 1),
(5, 2);

INSERT INTO categories (name, created_by, updated_by)
VALUES 
('Fitness Equipment', 2, 2),
('Outdoor Gear', 2, 2),
('Nutrition & Supplements', 3, 3);

INSERT INTO products (name, stock, description, category_id, price, created_by, updated_by)
VALUES 
('Adjustable Dumbbell 20kg', 10, 'Customizable weight dumbbell for home workouts', 1, 1200000.00, 2, 2),
('Treadmill Compact X100', 4, 'Foldable treadmill with digital display', 1, 5500000.00, 2, 2),
('Tent 2-Person Waterproof', 8, 'Outdoor tent ideal for camping', 2, 875000.00, 3, 3),
('Camping Stove Mini', 15, 'Portable gas stove for outdoor use', 2, 320000.00, 3, 3),
('Whey Protein 1kg', 20, 'Chocolate flavor protein supplement', 3, 450000.00, 2, 2),
('Electrolyte Drink Pack (12x)', 30, 'Hydration booster during exercise', 3, 180000.00, 3, 3);

INSERT INTO orders (customer_id, number_display, date, status, total, created_by, updated_by)
VALUES 
(1, 'ORD-20250601-001', '2025-06-01', 'completed', 1700000.00, 2, 2),
(2, 'ORD-20250601-002', '2025-06-01', 'processing', 630000.00, 3, 3);

-- Order 1: Dumbbell + Whey Protein
INSERT INTO order_details (order_id, product_id, qty, created_by, updated_by)
VALUES 
(1, 1, 1, 2, 2),  -- Adjustable Dumbbell
(1, 5, 1, 2, 2);  -- Whey Protein

-- Order 2: Camping Stove + Electrolyte Drink
INSERT INTO order_details (order_id, product_id, qty, created_by, updated_by)
VALUES 
(2, 4, 1, 3, 3),  -- Camping Stove
(2, 6, 1, 3, 3);  -- Electrolyte Drink

-- Billing for Order 1
INSERT INTO billings (order_id, number_display, tax, total, status, created_by, updated_by)
VALUES 
(1, 'BILL-20250601-001', 50000.00, 1700000.00, 'paid', 2, 2);

-- Billing for Order 2
INSERT INTO billings (order_id, number_display, tax, total, status, created_by, updated_by)
VALUES 
(2, 'BILL-20250601-002', 30000.00, 630000.00, 'unpaid', 3, 3);

-- Payment for Billing 1
INSERT INTO payments (billing_id, date, amount, method, created_by, updated_by)
VALUES 
(1, '2025-06-01 10:00:00', 1700000.00, 'transfer', 2, 2);

SELECT id, username, email, role, password FROM users WHERE username = 'admin01';
