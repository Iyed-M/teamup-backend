package project_handler

import (
	project_service "github.com/Iyed-M/teamup-backend/service/project"
	"github.com/Iyed-M/teamup-backend/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *projectHandler) CreateProject(c *fiber.Ctx) error {
	var req types.CreateProjectRequest
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
