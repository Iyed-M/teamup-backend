package project_handler

import (
	"github.com/Iyed-M/teamup-backend/internal/repository"
	project_service "github.com/Iyed-M/teamup-backend/service/project"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CreateProjectRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}
type createProjectResponse = repository.Project

func (h *projectHandler) CreateProject(c *fiber.Ctx) error {
	var req CreateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	userId := c.Locals("userId").(uuid.UUID)
	project, err := project_service.CreateProject(c.Context(), h.conn, userId, req.Name, req.Color)
	if err != nil {
		return err
	}
	return c.JSON(project)
}
