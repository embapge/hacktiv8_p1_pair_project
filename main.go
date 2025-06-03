package main

import (
	"fmt"
	"pairproject/config"
)

func main() {
	db := config.InitDB()
	fmt.Println(db)
}
