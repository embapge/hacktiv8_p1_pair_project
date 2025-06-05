package cli

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"pairproject/entity"
	"pairproject/handler"
	"pairproject/utils"
	"strconv"
	"strings"
)

// Variabel global products (contoh daftar produk statis, bisa dihapus jika menggunakan database)
var products = []string{"Bola", "Decker", "Stick Baseball", "Stick Golf"}

// Scanner untuk membaca input dari terminal
var scanner = bufio.NewScanner(os.Stdin)

// Struct cliHandler berisi koneksi database dan context aplikasi
type cliHandler struct {
	db  *sql.DB
	ctx context.Context
}

// Konstruktor untuk membuat instance baru cliHandler
func NewCLIHandler(db *sql.DB, ctx context.Context) *cliHandler {
	return &cliHandler{
		db:  db,
		ctx: ctx,
	}
}

// Fungsi utama untuk menampilkan menu utama CLI
func (c *cliHandler) Menu() {
MainMenu: for {
		fmt.Println("\n\n=== Welcome to Bandit Sports ===")
		fmt.Println("1. Login")
		fmt.Println("2. Register")
		fmt.Println("3. Products List")
		fmt.Println("4. Exit")
		fmt.Print("Choose menu: ")
		choice := readInput() // Fungsi pembantu untuk baca input

		switch choice {
		case "1":
			// Proses login user
			fmt.Println("=== Login ===")
			fmt.Print("Username: ")
			username := readInput()
			fmt.Print("Password: ")
			password := readInput()

			// Hapus spasi berlebih agar validasi lebih akurat
			username = strings.TrimSpace(username)
			password = strings.TrimSpace(password)

			// Panggil handler untuk login, validasi username dan password di database
			userHandler := handler.AuthHandler{DB: c.db}
			user, err := userHandler.LoginUser(username, password)

			if err != nil {
				// Jika login gagal, beri pesan dan kembali ke menu utama
				fmt.Println("Username or password incorrect please loggin again.")
				continue MainMenu
			}

			// Jika login berhasil, simpan user di context agar bisa dipakai di menu lain
			c.ctx = utils.WithUser(c.ctx, user)
			c.showLoggedInMenu() // Tampilkan menu sesuai role user
		case "2":
			// Proses registrasi user baru
			fmt.Println("=== Register ===")
			var userRegis entity.CustomerRegister

			// Ambil data registrasi dari user via input
			fmt.Print("Full Name: ")
			userRegis.Name = readInput()
			fmt.Print("Address: ")
			userRegis.Address = readInput()
			fmt.Print("Email: ")
			userRegis.Email = readInput()
			fmt.Print("Phone Number: ")
			userRegis.Phone = readInput()

			fmt.Println("\nCreate Login Credentials:")
			fmt.Print("Username: ")
			userRegis.Username = readInput()
			fmt.Print("Password: ")
			password := readInput()
			fmt.Print("Confirm Password: ")
			confirm := readInput()

			// Validasi password dan konfirmasi password harus sama
			if password != confirm {
				fmt.Println("Passwords do not match.")
				return
			}
			userRegis.Password = password

			// Panggil handler untuk menyimpan user baru ke database
			authHandler := handler.AuthHandler{DB: c.db}
			err := authHandler.Register(&userRegis)

			if err != nil {
				fmt.Printf("%s\n\n", err)
				continue
			}

			fmt.Println("Register successfully")
		case "3":
			// Menampilkan daftar produk dari database
			productHandler := handler.ProductHandler{DB: c.db}
			products, _ := productHandler.GetProducts()
			PrintProducts(products) // Fungsi bantu untuk print produk ke terminal
		case "4":
			// Exit program
			fmt.Println("Exiting program. Goodbye!")
			os.Exit(0)
		default:
			// Jika input pilihan tidak valid
			fmt.Println("Invalid option.")
		}
	}
}

// Menampilkan menu khusus setelah user berhasil login
func (c *cliHandler) showLoggedInMenu() {
	// Ambil user dari context
	user, ok := utils.GetUser(c.ctx)

	if !ok {
		fmt.Println("Failed to get user from context.")
		return
	}

	role := user.Role

	// Tampilkan menu berbeda berdasarkan role user (admin/customer)
	if role == "admin" {
		c.adminMenu()
	} else if role == "customer" {
		c.customerMenu()
	}
}

// Menu khusus untuk user dengan role admin
func (c *cliHandler) adminMenu() {
	// Handler untuk kategori, produk, dan laporan, siap digunakan di menu admin
	categoryHandler := handler.CategoryHandler{DB: c.db, Ctx: &c.ctx}
	productHandler := handler.ProductHandler{DB: c.db, Ctx: &c.ctx}
	reportHandler := handler.ReportHandler{DB: c.db}

AdminMenuLabel: for {
		fmt.Println("=== Admin Menu ===")
		fmt.Println("1. Create Category")
		fmt.Println("2. Create Product")
		fmt.Println("3. Report Most Sold Items")
		fmt.Println("4. Report Unpaid Bills")
		fmt.Println("5. Detail Revenue")
		fmt.Println("6. Logout")
		fmt.Print("Choose option: ")
		choice := readInput()

		switch choice {
		case "1":
			// Membuat kategori baru
			fmt.Print("Enter category name: ")
			categoryName := readInput()
			err := categoryHandler.CreateCategory(categoryName)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Kategori berhasil dibuat")
			}
		case "2":
			// Membuat produk baru
			var product entity.Product

			fmt.Print("Enter product name: ")
			product.Name = readInput()

			fmt.Print("Enter product stock: ")
			stockInput := readInput()
			_, err := fmt.Sscanf(stockInput, "%d", &product.Stock)
			if err != nil {
				fmt.Println("Invalid stock value.")
				break
			}

			fmt.Print("Enter category ID: ")
			categoryInput := readInput()
			_, err = fmt.Sscanf(categoryInput, "%d", &product.CategoryID)
			if err != nil {
				fmt.Println("Invalid category ID.")
				break
			}

			fmt.Print("Enter product description: ")
			product.Description = readInput()

			fmt.Print("Enter product price: ")
			priceInput := readInput()
			_, err = fmt.Sscanf(priceInput, "%f", &product.Price)
			if err != nil {
				fmt.Println("Invalid price.")
				break
			}

			// Simpan produk baru via handler
			err = productHandler.CreateProduct(product)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Product berhasil dibuat")
			}
		case "3":
			// Tampilkan laporan produk paling banyak terjual
			result, err := reportHandler.GetMostSoldProducts()
			if err != nil {
				fmt.Println("Error fetching report:", err)
				break
			}

			fmt.Println("=== Most Sold Products ===")
			for i, item := range result {
				fmt.Printf("%d. %s (Sold: %d)\n", i+1, item.Name, item.TotalSold)
			}
		case "4":
			// Tampilkan laporan tagihan yang belum dibayar
			unpaidBills, err := reportHandler.GetUnpaidBills()
			if err != nil {
				fmt.Println("Error fetching unpaid bills:", err)
				break
			}

			if len(unpaidBills) == 0 {
				fmt.Println("No unpaid bills found.")
				break
			}

			fmt.Println("=== Unpaid Bills Report ===")
			fmt.Printf("%-5s | %-17s | %-16s | %-20s | %-10s | %-10s | %-10s | %-20s\n",
				"ID", "Bill No", "Order No", "Customer", "Tax", "Total", "Status", "Created At")
			fmt.Println(strings.Repeat("-", 110))
			for _, bill := range unpaidBills {
				fmt.Printf("%-5d | %-15s | %-15s | %-20s | %-10.2f | %-10.2f | %-10s | %-20s\n",
					bill.ID, bill.BillNumber, bill.OrderNumber, bill.CustomerName,
					bill.Tax, bill.Total, bill.Status, bill.CreatedAt)
			}
		case "5":
			// Tampilkan laporan detail revenue/pemasukan
			revenueList, err := reportHandler.GetRevenueDetails()
			if err != nil {
				fmt.Println("Error fetching revenue details:", err)
				break
			}

			if len(revenueList) == 0 {
				fmt.Println("No revenue data found.")
				break
			}

			fmt.Println("=== Revenue Detail Report ===")
			fmt.Printf("%-17s | %-20s | %-10s | %-12s | %-20s | %-15s\n",
				"Bill No", "Payment Date", "Amount", "Method", "Customer", "Order No")
			fmt.Println(strings.Repeat("-", 100))
			for _, r := range revenueList {
				fmt.Printf("%-15s | %-20s | %-10.2f | %-12s | %-20s | %-15s\n",
					r.BillNumber, r.PaymentDate, r.Amount, r.Method, r.CustomerName, r.OrderNumber)
			}
		case "6":
			// Logout user dan kembali ke menu utama
			fmt.Println("User Logout...")
			c.ctx = utils.ClearUser(c.ctx)
			break AdminMenuLabel
		default:
			// Pilihan tidak dikenali, kembali ke menu admin
			fmt.Println("Returning to main menu...")
		}
	}
}

// customerMenu menampilkan menu khusus customer dan menangani aksi-aksi terkait order, billing, dan pembayaran
func (c *cliHandler) customerMenu(){
	orderHandler := handler.OrderHandler{DB: c.db, Ctx: &c.ctx}
	billingHandler := handler.BillingHandler{DB: c.db, Ctx: &c.ctx}

	CustomerMenuLabel: for {
		fmt.Println("\n\n=== Customer Menu ===")
		fmt.Println("1. Add Order")       // Tambah order baru
		fmt.Println("2. Update Order")    // Update detail order (kuantitas)
		fmt.Println("3. Create Billing")  // Generate tagihan dari order
		fmt.Println("4. Add Payment")     // Tambah pembayaran untuk tagihan
		fmt.Println("5. Log Out")         // Logout user
		fmt.Print("Choose option: ")
		option := readInput()

		switch option {
		case "1":
			// Mendapatkan daftar produk dari database
			productHandler := handler.ProductHandler{DB: c.db}
			products, err := productHandler.GetProducts()
			if err != nil {
				// Jika error, set produk kosong agar tidak crash
				var p []entity.Product
				products = p
			}

			// Tampilkan daftar produk
			fmt.Println("List Product:")
			for _, p := range products {
				fmt.Printf("Id: %d, Name: %s, Category: %s, Stock: %d, Price: %.2f, Description: %s\n",
					p.ID, p.Name, p.Category.Name, p.Stock, p.Price, p.Description)
			}
			fmt.Println("")

			// Input order secara berulang sampai user tidak ingin menambah lagi
			var orders []entity.OrderProduct
			for {
				fmt.Print("Masukkan ProductId: ")
				productIdStr := readInput()
				productId, err := strconv.Atoi(productIdStr)
				if err != nil {
					fmt.Println("Invalid ProductId. Please enter a number.")
					continue
				}

				fmt.Print("Masukkan Qty: ")
				qtyStr := readInput()
				qty, err := strconv.Atoi(qtyStr)
				if err != nil {
					fmt.Println("Invalid Qty. Please enter a number.")
					continue
				}

				orders = append(orders, entity.OrderProduct{ProductId: productId, Qty: qty})

				fmt.Print("Masih ingin memesan (y/n): ")
				isStillOrder := readInput()
				if isStillOrder != "y" {
					break
				}
			}

			// Kirim data order ke handler untuk dibuatkan order baru
			_, err = orderHandler.CreateOrder(orders)
			if err != nil {
				fmt.Printf("%v\n\n", err)
				continue CustomerMenuLabel
			}

			fmt.Println("Order berhasil dibuat.")

		case "3":
			// Ambil daftar order customer dari DB
			orders, err := orderHandler.GetOrders()
			if err != nil {
				fmt.Println("Failed to get orders:", err)
				return
			}
			printOrders(orders) // Tampilkan order yang ada

			if len(orders) == 0{
				continue
			}

			fmt.Print("Silahkan masuk nomor order: ")
			orderNumber := readInput()

			// Cari order berdasarkan nomor yang dimasukkan
			var filteredOrder entity.Order
			var findOrder bool
			for _, order := range orders {
				if order.NumberDisplay == orderNumber {
					filteredOrder = order
					findOrder = true
					break
				}
			}

			if !findOrder {
				fmt.Printf("Order id tidak ditemukan.\n\n")
				break
			}

			// Generate tagihan dari order terpilih
			billing, err := billingHandler.GenerateBill(filteredOrder)
			if err != nil {
				fmt.Printf("%v.\n\n", err)
				break
			}

			// Tampilkan info tagihan dan instruksi pembayaran
			fmt.Printf("Silahkan melakukan pembayaran atas tagihan: %s dengan nominal: %.2f maksimal di pukul: %s\n\n",
				billing.NumberDisplay, billing.Total, billing.DueDate)

		case "2":
			// Ambil daftar order customer untuk update detail order
			orders, err := orderHandler.GetOrders()
			if err != nil {
				fmt.Println("Failed to get orders:", err)
				return
			}
			printOrders(orders) // Tampilkan order

			if len(orders) == 0{
				continue
			}

			// Input order detail yang ingin diupdate
			fmt.Print("Silahkan masuk id order detail: ")
			orderDetailIdStr := readInput()
			orderDetailId, err := strconv.Atoi(orderDetailIdStr)
			if err != nil {
				fmt.Println("Invalid orderDetailId.")
				continue
			}

			// Input kuantitas baru
			fmt.Print("Kuantitas baru: ")
			qtyStr := readInput()
			qty, err := strconv.Atoi(qtyStr)
			if err != nil {
				fmt.Println("Invalid Qty. Please enter a number.")
				continue
			}

			// Update order detail via handler
			orderDetailHandler := handler.OrderDetailHandler{DB: c.db, Ctx: &c.ctx}
			_, err = orderDetailHandler.UpdateDetail(orderDetailId, qty)
			if err != nil {
				fmt.Printf("%v\n", err)
				continue
			}

			fmt.Println("Order detail berhasil diupdate")

		case "4":
			// Input dan proses pembayaran tagihan
			var paymentMethod entity.Method
			var isOkPay bool
			paymentHandler := handler.PaymentHandler{DB: c.db, Ctx: &c.ctx}

			for {
				// Input nomor tagihan yang ingin dibayar
				fmt.Print("Silahkan masukan nomor bill: ")
				billNumberDisplay := readInput()
				billing, err := billingHandler.GetBillByNumberDisplay(billNumberDisplay)
				if err != nil {
					fmt.Printf("%v\n", err)
					continue CustomerMenuLabel
				}

				// Tampilkan opsi metode pembayaran yang tersedia
				fmt.Println("===== List Jenis Pembayaran =====")
				fmt.Printf("- %s\n", entity.MethodCredit)
				fmt.Printf("- %s\n", entity.MethodVA)
				fmt.Printf("- %s\n", entity.MethodTransfer)
				fmt.Print("Silahkan masukan jenis pembayaran: ")
				paymentMethodInput := readInput()

				// Validasi input metode pembayaran
				if paymentMethodInput == string(entity.MethodCredit) {
					paymentMethod = entity.MethodCredit
					isOkPay = true
				} else if paymentMethodInput == string(entity.MethodVA) {
					paymentMethod = entity.MethodVA
					isOkPay = true
				} else if paymentMethodInput == string(entity.MethodTransfer) {
					paymentMethod = entity.MethodTransfer
					isOkPay = true
				}

				if !isOkPay {
					fmt.Println("Input tidak valid. Silahkan input ulang")
					continue
				}

				// Input nominal pembayaran
				fmt.Print("Nominal Pembayaran: ")
				amountInput := readInput()

				amount, err := strconv.ParseFloat(amountInput, 64)
				if err != nil {
					fmt.Println("Error converting amount:", err)
					return
				}

				// Buat pembayaran menggunakan handler
				err = paymentHandler.CreatePayment(&billingHandler, billing, amount, paymentMethod)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Pembayaran berhasil dibuat.")
				}

				break
			}

		case "5":
			// Logout dan hapus user dari context
			c.ctx = utils.ClearUser(c.ctx)
			break CustomerMenuLabel

		default:
			// Default jika input tidak valid
			fmt.Println("Returning to main menu...")
			break
		}
	}
}

// readInput membaca input baris dari stdin dan menghapus spasi kosong di depan dan belakang
func readInput() string {
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

// printOrders menampilkan daftar order dengan detailnya dalam format tabel
func printOrders(orders []entity.Order) {
	fmt.Printf("\n===== List Order =====\n")
	if len(orders) == 0 {
		fmt.Println("No orders found.")
		return
	}

	// Header tabel order
	fmt.Printf("%-10s %-15s %-12s %-10s %-10s\n", "Order ID", "NumberDisplay", "Order Date", "Status", "Total")
	fmt.Println(strings.Repeat("-", 60))

	// Loop tiap order dan print detailnya
	for _, order := range orders {
		fmt.Printf("%-10d %-15s %-12s %-10s %-10.2f\n",
			order.ID, order.NumberDisplay, order.Date, order.Status, order.Total)
		fmt.Printf("%-20s %-10s %-20s %-8s %-10s\n",
			"OrderDetailID", "ProductID", "Name", "Qty", "Subtotal")

		// Loop tiap detail produk dalam order
		for _, detail := range order.Details {
			fmt.Printf("%-20d %-10d %-20s %-8d %-10.2f\n",
				detail.ID, detail.ProductID, detail.Product.Name, detail.Qty, detail.Total)
		}
		fmt.Println(strings.Repeat("-", 60))
	}
}

// PrintProducts menampilkan daftar produk dalam format tabel
func PrintProducts(products []entity.Product) {
	fmt.Printf("%-5s %-20s %-6s %-10s %-15s %-10s\n",
		"ID", "Name", "Stock", "Price", "Category", "Description")
	fmt.Println(strings.Repeat("-", 70))

	// Loop tiap produk dan tampilkan info dengan truncate description jika panjang
	for _, p := range products {
		fmt.Printf("%-5d %-20s %-6d %-10.2f %-15s %-10s\n",
			p.ID,
			p.Name,
			p.Stock,
			p.Price,
			p.Category.Name,
			truncateString(p.Description, 10),
		)
	}
}

// truncateString memotong string s jika lebih panjang dari maxLen dan menambahkan "..." di akhir
func truncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}
