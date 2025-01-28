package main

import (
	"context"

	"github.com/Iyed-M/teamup-backend/config"
	"github.com/Iyed-M/teamup-backend/handlers/auth_handler"
	team_handler "github.com/Iyed-M/teamup-backend/handlers/team_handeler"
	"github.com/Iyed-M/teamup-backend/internal/repository"
	"github.com/Iyed-M/teamup-backend/service/jwt"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5"
)

func main() {
	log.SetLevel(log.LevelTrace)
	envVars, err := config.ParseEnv()
	ctx := context.Background()
	if err != nil {
		log.Fatal(err)
	}

	repo, conn, err := config.InitDB(envVars.DbURL, ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	app := fiber.New(fiber.Config{
		ErrorHandler: restErrorHandler,
		Prefork:      true,
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Group("/api")

	jwtService := jwt.NewJwtService([]byte(envVars.JWTSecret), envVars.JWTAccessDuration, envVars.JWTRefreshDuration)
	addRestEndpoints(app, conn, repo, jwtService)
	log.Fatal(app.Listen(":" + envVars.Port))
}

func addRestEndpoints(app *fiber.App, conn *pgx.Conn, repo *repository.Queries, jwtService jwt.JwtService) {
	authHandler := auth_handler.NewAuthHandler(jwtService, repo)
	app.Post("/signup", authHandler.Signup)
	app.Post("/login", authHandler.Login)
	app.Post("/refresh", authHandler.Refresh)
	app.Post("/logout", authHandler.Logout)
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("secret")},
	}))
	teamHandler := team_handler.NewTeamHandler(repo, conn)
	app.Post("/teams", teamHandler.CreateTeam)
	app.Post("/teams/invite", teamHandler.InviteTeamMember)
}

func restErrorHandler(c *fiber.Ctx, err error) error {
	if e, ok := err.(*fiber.Error); ok {
		return c.Status(e.Code).JSON(fiber.Map{
			"error": e.Message,
		})
	}
	log.Errorw("Unhandled Error", "err", err)
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "Internal Server Error",
	})
}
