package main

import (
	"context"

	"github.com/Iyed-M/teamup-backend/config"
	"github.com/Iyed-M/teamup-backend/internal/repository"
	"github.com/Iyed-M/teamup-backend/service/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5"
)

func main() {
	log.SetLevel(log.LevelTrace)
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, cfg.DbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	queries := repository.New(conn)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if e, ok := err.(*fiber.Error); ok {
				return c.Status(e.Code).JSON(fiber.Map{
					"error": e.Message,
				})
			}
			log.Errorw("Unhandled Error", "err", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		},
	})

	app.Use(logger.New())
	app.Use(recover.New())

	auth := auth.NewAuthService([]byte(cfg.JWTSecret), cfg.JWTAccessDuration, cfg.JWTRefreshDuration, queries)

	app.Post("/signup", auth.Signup)
	app.Post("/login", auth.Login)
	app.Post("/refresh", auth.Refresh)
	app.Post("/logout", auth.Logout)

	log.Fatal(app.Listen(":8080"))
}
