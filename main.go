package main

import (
	"context"

	"github.com/Iyed-M/teamup-backend/config"
	"github.com/Iyed-M/teamup-backend/handlers/auth_handler"
	project_handler "github.com/Iyed-M/teamup-backend/handlers/project_handeler"
	"github.com/Iyed-M/teamup-backend/internal/repository"
	jwt_service "github.com/Iyed-M/teamup-backend/service/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
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
	p, _ := bcrypt.GenerateFromPassword([]byte("testtest"), bcrypt.DefaultCost)
	log.Infow("HASH", "pas", string(p))

	app.Use(logger.New())
	app.Use(recover.New())
	app.Group("/api")

	jwtService := jwt_service.NewJwtService([]byte(envVars.JWTSecret), envVars.JWTAccessDuration, envVars.JWTRefreshDuration)
	addRestEndpoints(app, conn, repo, jwtService)
	log.Fatal(app.Listen(":" + envVars.Port))
}

func addRestEndpoints(app *fiber.App, conn *pgx.Conn, repo *repository.Queries, jwtService jwt_service.JwtService) {
	authHandler := auth_handler.NewAuthHandler(jwtService, repo)
	app.Post("/signup", authHandler.Signup)
	app.Post("/login", authHandler.Login)
	app.Post("/refresh", authHandler.Refresh)
	app.Post("/logout", authHandler.Logout)
	app.Use(jwtService.Middleware)
	projectHandler := project_handler.NewProjectHandler(repo, conn)
	app.Post("/projects", projectHandler.CreateProject)
	app.Post("/projects/invite", projectHandler.InviteProjectMember)
	app.Post("/projects/join", projectHandler.JoinProject)
	app.Get("/projects", projectHandler.ListProjects)
	app.Get("/projects/:projectId", projectHandler.GetProjectByID)
	app.Get("/projects/:projectId/tasks", projectHandler.GetProjectTasks)
	app.Get("/projects/invitations", projectHandler.ListInvitations)
	// app.Post("/projects/:projectId/task", projectHandler.)
}

func restErrorHandler(c *fiber.Ctx, err error) error {
	if e, ok := err.(*fiber.Error); ok {
		return c.Status(e.Code).JSON(fiber.Map{
			"error": e.Message,
		})
	}
	log.Tracew("Unhandled Error", "err", err)
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "Internal Server Error",
	})
}
