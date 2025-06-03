package cli

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"release2/handler"
	"strings"
)

func ShowMenu(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("=== Video Game Report Menu ===")
		fmt.Println("1. Total Game Sales Report")
		fmt.Println("2. Most Popular Game Report")
		fmt.Println("3. Total Revenue Per Game Report")
		fmt.Println("4. Player Count Per Game Report")
		fmt.Println("5. Exit")
		fmt.Print("Choose menu [1-5]: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		fmt.Print("Mau berapa limit: ")

		inputLimit, _ := reader.ReadString('\n')
		inputLimit = strings.TrimSpace(inputLimit)

		switch input {
		case "1":

			handler.ShowTotalSales(db, inputLimit)
		case "2":
			handler.ShowMostPopular(db, inputLimit)
		case "3":
			handler.ShowTotalRevenue(db, inputLimit)
		case "4":
			handler.ShowPlayerEngagement(db, inputLimit)
		case "5":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please select from 1 to 5.")
		}

		fmt.Println("\nPress ENTER to continue...")
		reader.ReadString('\n')A
	}
}
