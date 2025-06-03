package auth

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"pairproject/handler"
	"pairproject/utils"
	"strings"
)

func Login(db *sql.DB, ctx *context.Context) (context.Context, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')

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