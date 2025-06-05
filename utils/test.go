package utils

import (
	"context"
	"database/sql"
	"pairproject/entity"
	"testing"
)

const userTestKey ContextKey = "user"

func NewTestContextWithUser() context.Context {
	user := &entity.User{
		ID: 1,
		Customer: entity.Customer{
			ID: 1,
		},
	}
	ctx := context.Background()
	return context.WithValue(ctx, userTestKey, user)
}

func SetupTestMainDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	// Skema database lengkap untuk semua tabel yang relevan
	// Termasuk orders, billings, payments, dan tabel lain yang mungkin ada
	schema := `
-- Tabel users
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('admin', 'staff', 'customer')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Trigger untuk updated_at di tabel users
CREATE TRIGGER update_users_updated_at
AFTER UPDATE ON users
FOR EACH ROW
BEGIN
    UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

CREATE INDEX idx_users_role ON users (role);
CREATE INDEX idx_users_created_at ON users (created_at);

---

-- Tabel customers
CREATE TABLE customers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    phone_number TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER,
    updated_by INTEGER,
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id)
);

-- Trigger untuk updated_at di tabel customers
CREATE TRIGGER update_customers_updated_at
AFTER UPDATE ON customers
FOR EACH ROW
BEGIN
    UPDATE customers SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

CREATE INDEX idx_customers_name ON customers (name);
CREATE INDEX idx_customers_email ON customers (email);

---

-- Tabel user_customers
CREATE TABLE user_customers (
    user_id INTEGER NOT NULL,
    customer_id INTEGER NOT NULL PRIMARY KEY,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id)
);

-- Trigger untuk updated_at di tabel user_customers
CREATE TRIGGER update_user_customers_updated_at
AFTER UPDATE ON user_customers
FOR EACH ROW
BEGIN
    UPDATE user_customers SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

---

-- Tabel categories
CREATE TABLE categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER NOT NULL,
    updated_by INTEGER,
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id)
);

-- Trigger untuk updated_at di tabel categories
CREATE TRIGGER update_categories_updated_at
AFTER UPDATE ON categories
FOR EACH ROW
BEGIN
    UPDATE categories SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

---

-- Tabel products
CREATE TABLE products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    stock INTEGER DEFAULT 0 NOT NULL,
    description TEXT,
    category_id INTEGER NOT NULL,
    price NUMERIC DEFAULT 0 NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER NOT NULL,
    updated_by INTEGER,
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id),
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

-- Trigger untuk updated_at di tabel products
CREATE TRIGGER update_products_updated_at
AFTER UPDATE ON products
FOR EACH ROW
BEGIN
    UPDATE products SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

CREATE INDEX idx_products_category_id ON products (category_id);
CREATE INDEX idx_products_price ON products (price);
CREATE INDEX idx_products_name ON products (name);

---

-- Tabel orders
CREATE TABLE orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    customer_id INTEGER NOT NULL,
    number_display TEXT UNIQUE,
    date TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('processing', 'completed', 'cancel')) DEFAULT 'processing',
    total NUMERIC NOT NULL DEFAULT 0 CHECK (total >= 0),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER NOT NULL,
    updated_by INTEGER,
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id)
);

-- Trigger untuk updated_at di tabel orders
CREATE TRIGGER update_orders_updated_at
AFTER UPDATE ON orders
FOR EACH ROW
BEGIN
    UPDATE orders SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

CREATE INDEX idx_orders_number_display ON orders (number_display);
CREATE INDEX idx_orders_date ON orders (date);
CREATE INDEX idx_orders_status ON orders (status);

---

-- Tabel order_details
CREATE TABLE order_details (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    qty INTEGER NOT NULL DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER NOT NULL,
    updated_by INTEGER,
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id),
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

-- Trigger untuk updated_at di tabel order_details
CREATE TRIGGER update_order_details_updated_at
AFTER UPDATE ON order_details
FOR EACH ROW
BEGIN
    UPDATE order_details SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

---

-- Tabel billings
CREATE TABLE billings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_id INTEGER NOT NULL,
    number_display TEXT,
    issue_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    due_date DATETIME,
    tax NUMERIC NOT NULL DEFAULT 0 CHECK (tax >= 0),
    total NUMERIC NOT NULL DEFAULT 0 CHECK (total >= 0),
    status TEXT NOT NULL CHECK (status IN ('unpaid', 'lesspaid', 'paid', 'cancelled', 'refunded')) DEFAULT 'unpaid',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER NOT NULL,
    updated_by INTEGER,
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id),
    FOREIGN KEY (order_id) REFERENCES orders(id)
);

-- Trigger untuk updated_at di tabel billings
CREATE TRIGGER update_billings_updated_at
AFTER UPDATE ON billings
FOR EACH ROW
BEGIN
    UPDATE billings SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

CREATE INDEX idx_billings_number_display ON billings (number_display);
CREATE INDEX idx_billings_status ON billings (status);

---

-- Tabel payments
CREATE TABLE payments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    billing_id INTEGER NOT NULL,
    date DATETIME DEFAULT CURRENT_TIMESTAMP,
    amount NUMERIC NOT NULL DEFAULT 0 CHECK (amount >= 0),
    method TEXT NOT NULL CHECK (method IN ('credit_card', 'va', 'transfer')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER NOT NULL,
    updated_by INTEGER,
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id),
    FOREIGN KEY (billing_id) REFERENCES billings(id)
);

-- Trigger untuk updated_at di tabel payments
CREATE TRIGGER update_payments_updated_at
AFTER UPDATE ON payments
FOR EACH ROW
BEGIN
    UPDATE payments SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

CREATE INDEX idx_payments_date ON payments (date);
CREATE INDEX idx_payments_method ON payments (method);

---

-- Triggers Pengganti Stored Procedures

-- Trigger pengganti trg_order_details_after_insert
CREATE TRIGGER trg_order_details_after_insert
AFTER INSERT ON order_details
FOR EACH ROW
BEGIN
    UPDATE orders
    SET total = (SELECT IFNULL(SUM(od.qty * p.price), 0)
                 FROM order_details od
                 JOIN products p ON od.product_id = p.id
                 WHERE od.order_id = NEW.order_id)
    WHERE id = NEW.order_id;
END;

-- Trigger pengganti trg_order_details_after_update
CREATE TRIGGER trg_order_details_after_update
AFTER UPDATE ON order_details
FOR EACH ROW
BEGIN
    -- Recalculate for old order ID if it was changed
    -- Only execute if OLD.order_id is different from NEW.order_id
    UPDATE orders
    SET total = (SELECT IFNULL(SUM(od.qty * p.price), 0)
                 FROM order_details od
                 JOIN products p ON od.product_id = p.id
                 WHERE od.order_id = OLD.order_id)
    WHERE id = OLD.order_id AND OLD.order_id != NEW.order_id; -- Hanya update jika order_id berubah

    -- Recalculate for new order ID
    -- This update always happens as it applies to the current new state of the order_detail
    UPDATE orders
    SET total = (SELECT IFNULL(SUM(od.qty * p.price), 0)
                 FROM order_details od
                 JOIN products p ON od.product_id = p.id
                 WHERE od.order_id = NEW.order_id)
    WHERE id = NEW.order_id;
END;

-- Trigger pengganti trg_order_details_after_delete
CREATE TRIGGER trg_order_details_after_delete
AFTER DELETE ON order_details
FOR EACH ROW
BEGIN
    UPDATE orders
    SET total = (SELECT IFNULL(SUM(od.qty * p.price), 0)
                 FROM order_details od
                 JOIN products p ON od.product_id = p.id
                 WHERE od.order_id = OLD.order_id)
    WHERE id = OLD.order_id;
END;

-- Trigger pengganti ValidatePaymentAmount dan trg_payment_before_insert
CREATE TRIGGER trg_payment_before_insert
BEFORE INSERT ON payments
FOR EACH ROW
BEGIN
    -- Hitung total pembayaran saat ini + jumlah pembayaran baru
    -- Bandingkan dengan total tagihan dari tabel billings
    SELECT
        CASE
            WHEN (
                (SELECT IFNULL(SUM(amount), 0) FROM payments WHERE billing_id = NEW.billing_id)
                + NEW.amount
            ) > (SELECT total FROM billings WHERE id = NEW.billing_id)
            THEN RAISE(ABORT, 'Total payment exceeds billing total')
        END;
END;

---

-- Data Initial

INSERT INTO users (username, email, password, role)
VALUES
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
(1, 'ORD-202506-001', '2025-06-01', 'completed', 1700000.00, 2, 2),
(2, 'ORD-202506-002', '2025-06-01', 'processing', 630000.00, 3, 3);

-- Order 1: Dumbbell + Whey Protein
INSERT INTO order_details (order_id, product_id, qty, created_by, updated_by)
VALUES
(1, 1, 1, 2, 2),   -- Adjustable Dumbbell
(1, 5, 1, 2, 2);   -- Whey Protein

-- Order 2: Camping Stove + Electrolyte Drink
INSERT INTO order_details (order_id, product_id, qty, created_by, updated_by)
VALUES
(2, 4, 1, 3, 3),   -- Camping Stove
(2, 6, 1, 3, 3);   -- Electrolyte Drink

-- Billing for Order 1
INSERT INTO billings (order_id, number_display, tax, total, status, created_by, updated_by)
VALUES
(1, 'BIL-202506-001', 50000.00, 1700000.00, 'paid', 2, 2);

-- Billing for Order 2
INSERT INTO billings (order_id, number_display, tax, total, status, created_by, updated_by)
VALUES
(2, 'BIL-202506-002', 30000.00, 630000.00, 'unpaid', 3, 3);

-- Payment for Billing 1
INSERT INTO payments (billing_id, date, amount, method, created_by, updated_by)
VALUES
(1, '2025-06-01 10:00:00', 1700000.00, 'transfer', 2, 2);

SELECT id, username, email, role, password FROM users WHERE username = 'admin01';
	`

	_, err = db.Exec(schema)
	if err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}
	return db
}