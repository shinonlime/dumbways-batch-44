package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func DbConnect() {
	var err error
	dbUrl := "postgres://postgres:shinonlime818@localhost:5432/personal-web"

	Conn, err = pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable connect to database: %v", err)
		os.Exit(1)
	}

	fmt.Println("Success connect to database")
}
