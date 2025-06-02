# Deskripsi  
- Online Store yang menjual peralatan olahraga. Mulai dari sepatu, celana, baju dsb. Pelanggan dapat menjadi member agar seluruh transaksi dapat ter-record.
- Batasan aplikasi. Untuk berbelanja customer wajib memiliki user. Customer bisa berbelanja tanpa akses login namun harus melalui admin.

# Entity & Attribute
- Users
    - id INT PK AUTO_INCREMENT
    - username Varchar UNIQUE NOT NULL
    - email Varchar(100) UNIQUE NOT NULL
    - password Varchar(255) NOT NULL
    - role Enum (admin, staff, customer)
    - created_at TIMESTAMP NOT NULL
    - updated_at TIMESTAMP NOT NULL
- Customer
    - id INT PK AUTO_INCREMENT
    - name Varchar(100) NOT NULL
    - address TEXT NOT NULL
    - email Varchar UNIQUE NOT NULL
    - phone_number Varchar UNIQUE NOT NULL
    - created_at TIMESTAMP NOT NULL
    - updated_at TIMESTAMP NOT NULL
- UserCustomer
    - customer_id INT PK FK NN
    - user_id INT PK FK NN
    - created_at TIMESTAMP NOT NULL
    - updated_at TIMESTAMP NOT NULL
- Categories
    - id INT PK AUTO_INCREMENT
    - name Varchar(100) NOT NULL
    - created_at TIMESTAMP NOT NULL
    - updated_at TIMESTAMP NOT NULL
- Products
    - id INT PK AUTO_INCREMENT
    - name Varchar(100) NOT NULL
    - stock INT DEFAULT 0
    - description TEXT
    - category_id INT FK NOT NULL
    - price DECIMAL(10,2) NOT NULL
    - created_at TIMESTAMP NOT NULL
    - updated_at TIMESTAMP NOT NULL
- Orders
    - id INT PK AUTO_INCREMENT
    - customer_id INT FK NOT NULL
    - number_display Varchar(50)
    - date DATE NOT NULL
    - status ENUM(pending, done, process, void, delivered)
    - subtotal Decimal default 0 check (subtotal >=0)
    - tax Decimal default 0 check (tax >=0)
    - total Decimal default 0 check (total >=0)
    - created_at TIMESTAMP NOT NULL
    - updated_at TIMESTAMP NOT NULL
- OrderDetails
    - id INT PK AUTO_INCREMENT
    - order_id INT FK NOT NULL
    - product_id INT FK NOT NULL
    - price DECIMAL(10,2) NOT NULL (price >=0)
    - qty INT NOT NULL
    - created_at TIMESTAMP NOT NULL
    - updated_at TIMESTAMP NOT NULL
- Billings
    - id INT PK AUTO_INCREMENT
    - order_id INT FK NOT NULL
    - number_display Varchar(50)
    - amount DECIMAL(10,2) DEFAULT 0
    - status ENUM('unpaid', 'paid', 'cancelled', 'refunded') DEFAULT 'unpaid',
    - created_at TIMESTAMP NOT NULL
    - updated_at TIMESTAMP NOT NULL
- Payments
    - id INT PK AUTO_INCREMENT
    - billing_id INT FK NOT NULL,
    - date DATETIME NOT NULL
    - amount DECIMAL(10, 2) DEFAULT 0 CHECK (amount >= 0)
    - method ENUM (credit_card, va, transfer)
    - created_at TIMESTAMP NOT NULL
    - updated_at TIMESTAMP NOT NULL

# Relationship

# Integrity Constraint

- One to One
- One to Many
- Many to Many

- Berhasil membuat sebuah entity dari input (CREATE UPDATE DELETE)

## Main Feature:
1. Bikin orders, hapus, edit
2. Bayar orders
<!-- 3. Bikin User -->
3. Bikin Customer
4. Bikin Product/Category