# 📊 Database Design Analysis - `pair_project`

## Anggota
### Mohammad Barata Putra Gusti
### Rudito Nugroho

## 🧾 Deskripsi Sistem
Sistem ini adalah *Online Store* yang menjual perlengkapan olahraga seperti alat fitness, perlengkapan outdoor, dan suplemen. Transaksi harus dilakukan oleh customer yang sudah memiliki akun (*user*) dan terhubung dengan data customer. Pembuatan akun admin hanya bisa dilakukan via SQL injection (secara manual) menggunakan prosedur tersimpan (*stored procedure*).

---

## 🗂️ Entity & Attributes

### 1. Users
- **PK**: `id`
- **Unique**: `username`, `email`
- **Other Constraints**: `NOT NULL`, Enum `role`
- **Timestamps**: `created_at`, `updated_at`

### 2. Customers
- **PK**: `id`
- **Unique**: `email`, `phone_number`
- **FK**: `created_by`, `updated_by → users(id)`
- **Timestamps**: `created_at`, `updated_at`

### 3. UserCustomers
- **PK**: `customer_id`
- **FKs**: `user_id → users(id)`, `customer_id → customers(id)`
- **Cardinality**: One-to-One mapping user to customer

### 4. Categories
- **PK**: `id`
- **Unique**: `name`
- **FKs**: `created_by`, `updated_by → users(id)`

### 5. Products
- **PK**: `id`
- **FKs**: `category_id → categories(id)`, `created_by`, `updated_by → users(id)`
- **Other Constraints**: `price`, `stock` default values
- **Indexes**: `name`, `price`, `category_id`

### 6. Orders
- **PK**: `id`
- **Unique**: `number_display`
- **FKs**: `customer_id → customers(id)`, `created_by`, `updated_by → users(id)`
- **Enum**: `status` (`processing`, `completed`, `cancel`)
- **Other Constraints**: `total >= 0`
- **Indexes**: `date`, `status`, `number_display`

### 7. OrderDetails
- **PK**: `id`
- **FKs**: `order_id → orders(id)`, `product_id → products(id)`, `created_by`, `updated_by → users(id)`
- **Constraints**: `qty >= 1`

### 8. Billings
- **PK**: `id`
- **FKs**: `order_id → orders(id)`, `created_by`, `updated_by → users(id)`
- **Enum**: `status` (`unpaid`, `lesspaid`, `paid`, `cancelled`, `refunded`)
- **Other Constraints**: `tax`, `total >= 0`
- **Unique**: `number_display`

### 9. Payments
- **PK**: `id`
- **FKs**: `billing_id → billings(id)`, `created_by`, `updated_by → users(id)`
- **Enum**: `method` (`credit_card`, `va`, `transfer`)
- **Constraints**: `amount >= 0`

---

## 🔗 Modality & Cardinality

| Relationship | Cardinality | Modality | Notes |
|--------------|-------------|----------|-------|
| Users → All Table | 1:N | Optional/Mandatory (updated_by/created_by) |  |
| Customers → Orders | 1:N | Optional | Order harus punya customer |
| Customer → UserCustomer | 1:1 | Mandatory | Kemunculan Customer harus beriringan dengan User |
| Orders → OrderDetails | 1:N | Mandatory | Order harus memiliki minimal satu detail |
| Products → OrderDetails | 1:N | Optional | Product terlibat dalam order |
| Orders → Billings | 1:1 | Optional | Billing opsional per order |
| Billings → Payments | 1:N | Optional | Bisa tidak dibayar, atau dibayar sebagian |

---

## 🔒 Integrity Constraints

### 1. Unique Constraints
- Users: `username`, `email`
- Customers: `email`, `phone_number`
- Orders & Billings: `number_display`

### 2. Foreign Keys & Referential Integrity
- Semua relasi antar tabel menggunakan `FOREIGN KEY` dengan cascading default.
- `created_by` dan `updated_by` mengacu pada `users(id)`.

### 3. Business Logic Validations
- Stored Procedure `sp_update_order_total` menjaga konsistensi `orders.total`.
- Trigger otomatis pada `order_details` memanggil SP ini saat INSERT/UPDATE/DELETE.
- Validasi jumlah pembayaran (via `ValidatePaymentAmount`) sebelum insert payment dilakukan dengan trigger.

---

## ⚙️ Triggers & Stored Procedures

### Stored Procedures:
- `sp_update_order_total(p_order_id)` → Hitung ulang total order.
- `ValidatePaymentAmount(p_billing_id, adjustment, OUT is_valid)` → Validasi total pembayaran terhadap tagihan.

### Triggers:
- `AFTER INSERT/UPDATE/DELETE` on `order_details` → Update total order.
- `BEFORE INSERT` on `payments` → Cegah kelebihan bayar via SIGNAL ERROR.

---

## ✅ Main Features Supported

- Create Orders
- Update Order Detail Qty
- Payment Validation & Processing
- Create Product
- Create Category
- Customer Registration (linked with User)

---

## 💡 Catatan Tambahan
- Struktur `user_customers` menghindari data duplikasi dan mempermudah traceability.

---

