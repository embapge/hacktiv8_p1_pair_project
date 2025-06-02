DROP DATABASE IF EXISTS pair_project;

CREATE DATABASE pair_project;

USE pair_project;

CREATE TABLE Users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role ENUM('admin', 'staff', 'customer') NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    INDEX idx_role (role),
    INDEX idx_created_at (created_at)
);

CREATE TABLE Customers (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    address TEXT NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    phone_number VARCHAR(20) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    INDEX idx_name (name),
    INDEX idx_email (email)
);

CREATE TABLE UserCustomer (
    user_id INT NOT NULL,
    customer_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, customer_id),
    FOREIGN KEY (user_id) REFERENCES Users(id),
    FOREIGN KEY (customer_id) REFERENCES Customers(id)
);

CREATE TABLE Categories (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE INDEX idx_name (name)
);

CREATE TABLE Products (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    stock INT DEFAULT 0,
    description TEXT,
    category_id INT NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (category_id) REFERENCES Categories(id),
    INDEX idx_category_id (category_id),
    INDEX idx_price (price),
    INDEX idx_name (name)
);

CREATE TABLE Orders (
    id INT PRIMARY KEY AUTO_INCREMENT,
    customer_id INT NOT NULL,
    number_display VARCHAR(50),
    date DATE NOT NULL,
    status ENUM('pending', 'done', 'process', 'void', 'delivered') NOT NULL,
    subtotal DECIMAL(10,2) DEFAULT 0 CHECK (subtotal >= 0),
    tax DECIMAL(10,2) DEFAULT 0 CHECK (tax >= 0),
    total DECIMAL(10,2) DEFAULT 0 CHECK (total >= 0),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES Customers(id),
    INDEX idx_order_number_display (number_display),
    INDEX idx_date (DATE),
    INDEX idx_status (STATUS)
);

CREATE TABLE OrderDetails (
    id INT PRIMARY KEY AUTO_INCREMENT,
    order_id INT NOT NULL,
    product_id INT NOT NULL,
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    qty INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (order_id) REFERENCES Orders(id),
    FOREIGN KEY (product_id) REFERENCES Products(id)
);

CREATE TABLE Billings (
    id INT PRIMARY KEY AUTO_INCREMENT,
    order_id INT NOT NULL,
    number_display VARCHAR(50),
    amount DECIMAL(10,2) DEFAULT 0 CHECK (amount >= 0),
    status ENUM('unpaid', 'paid', 'cancelled', 'refunded') DEFAULT 'unpaid',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (order_id) REFERENCES Orders(id),
    INDEX idx_billing_number_display (number_display),
    INDEX idx_status (STATUS)
);

CREATE TABLE Payments (
    id INT PRIMARY KEY AUTO_INCREMENT,
    billing_id INT NOT NULL,
    date DATETIME NOT NULL,
    amount DECIMAL(10,2) DEFAULT 0 CHECK (amount >= 0),
    method ENUM('credit_card', 'va', 'transfer') NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (billing_id) REFERENCES billings(id),
    INDEX idx_date (date),
    INDEX idx_method (method)
);
