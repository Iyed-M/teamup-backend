package main

import (
	"log"

	"github.com/Iyed-M/teamup-backend/config"
	"github.com/Iyed-M/teamup-backend/service/auth"
	"github.com/Iyed-M/teamup-backend/types"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := initDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	auth := auth.NewAuthSerive([]byte(cfg.JWTSecret), 1000, 1000, db)

	e.POST("/signup", auth.Signup)
	e.POST("/login", auth.Login)
	e.POST("/refresh", auth.Refresh)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

func initDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DbURL), &gorm.Config{})
	db.AutoMigrate(&types.User{}, &types.Team{}, &types.TeamPermission{}, &types.ProjectPermission{}, &types.Project{}, &types.DirectMessage{}, &types.TeamMessages{}, &types.Task{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
