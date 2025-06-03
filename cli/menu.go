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
	
			showAdminMenu(&ctx)
		case "2":
			handleRegister(db)
		case "3":
			showProducts()
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

func showAdminMenu(ctx *context.Context) {
	user, ok := utils.GetUser(*ctx)
	
	if !ok {
		fmt.Println("Failed to get user from context.")
		return
	}

	role := user.Role

	if role == "admin" {
		fmt.Println("=== Admin Menu ===")
		fmt.Println("1. Create Category")
		fmt.Println("2. Create Product")
		fmt.Println("3. Report Most Sold Items")
		fmt.Println("4. Report Unpaid Bills")
		fmt.Println("5. Detail Revenue")
		fmt.Println("6. Exit")
		fmt.Print("Choose option: ")
		_ = readInput()
		fmt.Println("Returning to main menu...")
	} else if role == "customer" {
		fmt.Println("=== Customer Menu ===")
		fmt.Println("1. Add Order & Billing")
		fmt.Println("2. Update Order")
		fmt.Println("3. Add Payment")
		fmt.Println("4. Log Out")
		fmt.Print("Choose option: ")
		_ = readInput()
		fmt.Println("Returning to main menu...")
	}
}

func readInput() string {
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}