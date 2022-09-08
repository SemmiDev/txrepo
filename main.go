package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Server struct {
	app   *fiber.App
	store Store
}

type CreateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *Server) createUserHandler(c *fiber.Ctx) error {
	var input CreateUserInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user := NewUser(input.Name, input.Email)

	err := s.store.CreateUser(c.Context(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

type UpdateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *Server) updatUserHandler(c *fiber.Ctx) error {
	var input UpdateUserInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userId, err := uuid.Parse(c.Query("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user := NewUserWithID(userId, input.Name, input.Email)
	err = s.store.UpdateUserTx(c.Context(), userId, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

func (s *Server) getUsersHandler(c *fiber.Ctx) error {
	user, err := s.store.FindUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

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
	app.Use(logger.New())
	app.Use(recover.New())

	server := &Server{
		app:   app,
		store: store,
	}

	server.app.Get("/users", server.getUsersHandler)
	server.app.Post("/users", server.createUserHandler)
	server.app.Put("/users", server.updatUserHandler)

	server.app.Listen(":3000")
}
