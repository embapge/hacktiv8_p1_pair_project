package main

import (
	"context"
	"pairproject/cli"
	"pairproject/config"
)

func main() {
	db := config.InitDB()
	ctx := context.Background()

	cli.Menu(db, &ctx)
	// for {
	// 	ctx, err := auth.Login(db, &ctx)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	user, ok := utils.GetUser(ctx)
	// 	fmt.Println(db, user, ok)

	// 	cli.Menu(db, &ctx)
	// 	// Panggil 1 handler menu

	// 	// // Tampilkan menu sesuai role user
	// 	// if !ok {
	// 	// 	fmt.Println("Guest menu")
	// 	// 	cli.GuestMenu()
	// 	// } else if user.Role == "admin" {
	// 	// 	fmt.Println("Admin menu")
	// 	// 	menu.AdminMenu(user)
	// 	// } else if user.Role == "customer" {
	// 	// 	fmt.Println("Customer menu")
	// 	// 	customer.CustomerMenu(user)
	// 	// }

	// 	// // Contoh logout
	// 	// var input string
	// 	// fmt.Print("Logout? (y/n): ")
	// 	// fmt.Scanln(&input)
	// 	// if input == "y" || input == "Y" {
	// 	// 	ctx = context.Background() // reset context (logout)
	// 	// 	continue // kembali ke login
	// 	// } else {
	// 	// 	break // keluar dari aplikasi
	// 	// }
	// }
}
