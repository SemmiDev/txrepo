package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func main() {
	conn, err := sql.Open("postgres", "postgresql://root:secret@localhost:5432/tx_repo?sslmode=disable")
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	if err := conn.Ping(); err != nil {
		log.Fatal("cannot ping db:", err)
	}

	store := NewStore(conn)
	app := fiber.New()

	server := &Server{app: app, store: store}
	server.setupRoutes()
	server.setupMiddlewares()

	server.app.Listen(":3000")
}
