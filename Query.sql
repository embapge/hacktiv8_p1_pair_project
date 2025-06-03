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
    created_by INT NOT NULL, 
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
 
CREATE TABLE Categories ( 
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
    FOREIGN KEY (category_id) REFERENCES Categories(id), 
    INDEX idx_category_id (category_id), 
    INDEX idx_price (price), 
    INDEX idx_name (name) 
); 
 
CREATE TABLE orders ( 
    id INT PRIMARY KEY AUTO_INCREMENT, 
    customer_id INT NOT NULL, 
    number_display VARCHAR(50), 
    date DATE NOT NULL, 
    status ENUM('processing', 'completed', 'cancel') DEFAULT('processing') NOT NULL,
    total DECIMAL(10,2) DEFAULT 0 CHECK (total >= 0) NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by INT NOT NULL, 
    updated_by INT, 
    FOREIGN KEY (created_by) REFERENCES users(id), 
    FOREIGN KEY (updated_by) REFERENCES users(id), 
    FOREIGN KEY (customer_id) REFERENCES customers(id), 
    INDEX idx_order_number_display (number_display), 
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
    number_display VARCHAR(50), 
    subtotal DECIMAL(10,2) DEFAULT 0 CHECK (subtotal >= 0) NOT NULL, 
    tax DECIMAL(10,2) DEFAULT 0 CHECK (tax >= 0) NOT NULL, 
    total DECIMAL(10,2) DEFAULT 0 CHECK (total >= 0) NOT NULL, 
    status ENUM('unpaid', 'paid', 'cancelled', 'refunded') DEFAULT 'unpaid', 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by INT NOT NULL, 
    updated_by INT, 
    FOREIGN KEY (created_by) REFERENCES users(id), 
    FOREIGN KEY (updated_by) REFERENCES users(id), 
    FOREIGN KEY (order_id) REFERENCES orders(id), 
    INDEX idx_billing_number_display (number_display), 
    INDEX idx_status (status)
); 
 
CREATE TABLE payments ( 
    id INT PRIMARY KEY AUTO_INCREMENT, 
    billing_id INT NOT NULL, 
    date DATETIME NOT NULL, 
    amount DECIMAL(10,2) DEFAULT 0 CHECK (amount >= 0) NOT NULL, 
    method ENUM('credit_card', 'va', 'transfer') NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by INT NOT NULL, 
    updated_by INT, 
    FOREIGN KEY (created_by) REFERENCES users(id), 
    FOREIGN KEY (updated_by) REFERENCES users(id), 
    FOREIGN KEY (billing_id) REFERENCES billings(id), 
    INDEX idx_date (date), 
    INDEX idx_method (method) 
); 
