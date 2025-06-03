package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type User struct {
	Name     string
	Address  string
	Email    string
	Phone    string
	Username string
	Password string
	Role     string // "admin" atau "customer"
}

var users = []User{}
var products = []string{"Bola", "Decker", "Stick Baseball", "Stick Golf"}
var scanner = bufio.NewScanner(os.Stdin)

func main() {
	for {
		showMainMenu()
	}
}

func showMainMenu() {
	fmt.Println("=== Welcome to Bandit Sports ===")
	fmt.Println("1. Login")
	fmt.Println("2. Register")
	fmt.Println("3. Products List")
	fmt.Println("4. Exit")
	fmt.Print("Choose menu: ")
	choice := readInput()

	switch choice {
	case "1":
		handleLogin()
	case "2":
		handleRegister()
	case "3":
		showProducts()
	default:
		fmt.Println("Invalid option.")
	}
}

func handleLogin() {
	fmt.Println("=== Login ===")
	fmt.Print("Username: ")
	username := readInput()
	fmt.Print("Password: ")
	password := readInput()

	for _, user := range users {
		if user.Username == username && user.Password == password {
			fmt.Printf("Login successful. Welcome, %s!\n\n", user.Name)
			if user.Role == "admin" {
				showAdminMenu()
			} else {
				showCustomerMenu()
			}
			return
		}
	}
	fmt.Println("Invalid credentials.\n")
}

func handleRegister() {
	fmt.Println("=== Register ===")
	var user User

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
		fmt.Println("Passwords do not match.\n")
		return
	}
	user.Password = password

	users = append(users, user)
	fmt.Println("Registration successful.\n")
}

func showProducts() {
	fmt.Println("=== Product List ===")
	for i, product := range products {
		fmt.Printf("%d. %s\n", i+1, product)
	}
	fmt.Println()
}

func showAdminMenu() {
	fmt.Println("=== Admin Menu ===")
	fmt.Println("1. Create Category")
	fmt.Println("2. Create Product")
	fmt.Println("3. Report Most Sold Items")
	fmt.Println("4. Report Unpaid Bills")
	fmt.Println("5. Detail Revenue")
	fmt.Println("6. Exit")
	fmt.Print("Choose option: ")
	_ = readInput()
	fmt.Println("Returning to main menu...\n")
}

func showCustomerMenu() {
	fmt.Println("=== Customer Menu ===")
	fmt.Println("1. Add Order & Billing")
	fmt.Println("2. Update Order")
	fmt.Println("3. Add Payment")
	fmt.Println("4. Log Out")
	fmt.Print("Choose option: ")
	_ = readInput()
	fmt.Println("Returning to main menu...\n")
}

func readInput() string {
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}