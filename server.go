package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
)

type Server struct {
	app   *fiber.App
	store Store
}

func (s *Server) setupRoutes() {
	s.app.Get("/users", s.getUsersHandler)
	s.app.Post("/users", s.createUserHandler)
	s.app.Put("/users", s.updatUserHandler)
}

func (s *Server) setupMiddlewares() {
	s.app.Use(logger.New())
	s.app.Use(recover.New())
}

func (s *Server) createUserHandler(c *fiber.Ctx) error {
	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

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

func (s *Server) updatUserHandler(c *fiber.Ctx) error {
	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

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

	user := NewUserWithoutID(input.Name, input.Email)
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
