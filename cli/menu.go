package cli

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"pairproject/entity"
	"pairproject/handler"
	"pairproject/utils"
	"strings"
)

var products = []string{"Bola", "Decker", "Stick Baseball", "Stick Golf"}
var scanner = bufio.NewScanner(os.Stdin)

func Menu(db *sql.DB, ctx *context.Context) {
	MainMenu: for{
		fmt.Println("=== Welcome to Bandit Sports ===")
		fmt.Println("1. Login")
		fmt.Println("2. Register")
		fmt.Println("3. Products List")
		fmt.Println("4. Exit")
		fmt.Print("Choose menu: ")
		choice := readInput()
	
		switch choice {
		case "1":
			ctx, err := handleLogin(db, ctx)
			if err != nil {
				log.Fatal(err)
				continue MainMenu
			}
	
			showAdminMenu(db, &ctx)
		case "2":
			handleRegister(db)
		case "3":
			showProducts()
		case "4":
			fmt.Println("Exiting program. Goodbye!")
			os.Exit(0)
		default:
			fmt.Println("Invalid option.")
		}
	}
}

func handleLogin(db *sql.DB, ctx *context.Context) (context.Context, error) {
	fmt.Println("=== Login ===")
	fmt.Print("Username: ")
	username := readInput()
	fmt.Print("Password: ")
	password := readInput()

	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	// Simulasi user lookup
	userHandler := handler.AuthHandler{DB: db}
	user, err := userHandler.LoginUser(username, password)

	if err != nil{
		return *ctx, err
	}

	return utils.WithUser(*ctx, user), nil
}

func handleRegister(db *sql.DB) {
	fmt.Println("=== Register ===")
	var user entity.CustomerRegister

	fmt.Print("Full Name: ")
	user.Name = readInput()
	fmt.Print("Address: ")
	user.Address = readInput()
	fmt.Print("Email: ")
	user.Email = readInput()
	fmt.Print("Phone Number: ")
	user.Phone = readInput()

	fmt.Println("\nCreate Login Credentials:")
	fmt.Print("Username: ")
	user.Username = readInput()
	fmt.Print("Password: ")
	password := readInput()
	fmt.Print("Confirm Password: ")
	confirm := readInput()

	if password != confirm {
		fmt.Println("Passwords do not match.")
		return
	}
	user.Password = password

	// Eksekusi handler.

	authHandler := handler.AuthHandler{DB: db}
	message, err := authHandler.Register(&user)

	if err != nil{
		log.Fatal(err)
	}

	fmt.Println(message)
}

func showProducts() {
	fmt.Println("=== Product List ===")
	for i, product := range products {
		fmt.Printf("%d. %s\n", i+1, product)
	}
	fmt.Println()
}

func showAdminMenu(db *sql.DB, ctx *context.Context) {
	// Ambil data user dari context (hasil login)
	user, ok := utils.GetUser(*ctx)
	if !ok {
		fmt.Println("Failed to get user from context.")
		return
	}

	role := user.Role

	// Jika user adalah admin, tampilkan menu admin
	if role == "admin" {

		for {
			// Tampilkan menu untuk admin
			fmt.Println("=== Admin Menu ===")
			fmt.Println("1. Create Category")
			fmt.Println("2. Create Product")
			fmt.Println("3. Report Most Sold Products")
			fmt.Println("4. Report Unpaid Bills")
			fmt.Println("5. Report Detail Revenue")
			fmt.Println("6. Log Out")
			fmt.Print("Choose option: ")
			choice := readInput() // baca input user

			switch choice {
			case "1":
				// Menu untuk membuat kategori baru
				h := handler.CategoryHandler{DB: db}
				var category entity.Category
				fmt.Print("Enter category name: ")
				category.Name = readInput()
				h.CreateCategory(ctx, &category)

			case "2":
				// Menu untuk membuat produk baru
				h := handler.ProductHandler{DB: db}
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

				h.CreateProduct(ctx, &product)

			case "3":
				// Menu untuk menampilkan laporan produk yang paling banyak terjual
				r := handler.ReportHandlerMostSold{DB: db}
				result, err := r.GetMostSoldProducts()
				if err != nil {
					fmt.Println("Error fetching report:", err)
					break
				}

				fmt.Println("=== Most Sold Products ===")
				for i, item := range result {
					fmt.Printf("%d. %s (Sold: %d)\n", i+1, item.Name, item.TotalSold)
				}

			case "4":
				// Menu untuk menampilkan tagihan yang belum dibayar
				r := handler.ReportHandlerUnpaidBill{DB: db}
				unpaidBills, err := r.GetUnpaidBills(ctx)
				if err != nil {
					fmt.Println("Error fetching unpaid bills:", err)
					break
				}

				if len(unpaidBills) == 0 {
					fmt.Println("No unpaid bills found.")
					break
				}

				// Format laporan tagihan yang belum dibayar
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
				// Menu untuk menampilkan detail pemasukan/revenue dari pembayaran
				r := handler.RevenueHandler{DB: db}
				revenueList, err := r.GetRevenueDetails(ctx)
				if err != nil {
					fmt.Println("Error fetching revenue details:", err)
					break
				}

				if len(revenueList) == 0 {
					fmt.Println("No revenue data found.")
					break
				}

				// Format laporan detail revenue
				fmt.Println("=== Revenue Detail Report ===")
				fmt.Printf("%-17s | %-20s | %-10s | %-12s | %-20s | %-15s\n",
					"Bill No", "Payment Date", "Amount", "Method", "Customer", "Order No")
				fmt.Println(strings.Repeat("-", 100))
				for _, r := range revenueList {
					fmt.Printf("%-15s | %-20s | %-10.2f | %-12s | %-20s | %-15s\n",
						r.BillNumber, r.PaymentDate, r.Amount, r.Method, r.CustomerName, r.OrderNumber)
				}

			case "6":
				// Keluar dari menu admin
				fmt.Println("Logging out...")
				return

			default:
				// Menangani input yang tidak valid
				fmt.Println("Invalid option.")
			}
		}

	} else if role == "customer" {
		// Menu sederhana jika user adalah customer
		fmt.Println("=== Customer Menu ===")
		fmt.Println("1. Add Order & Billing")
		fmt.Println("2. Update Order")
		fmt.Println("3. Add Payment")
		fmt.Println("4. Log Out")
		fmt.Print("Choose option: ")
		_ = readInput() // sementara belum di-handle detailnya
		fmt.Println("Returning to main menu...")
	}
}

func readInput() string {
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}
