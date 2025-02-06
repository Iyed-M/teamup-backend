package project_handler

import (
	"context"

	"github.com/Iyed-M/teamup-backend/internal/repository"
	"github.com/Iyed-M/teamup-backend/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (ph *projectHandler) JoinProject(c *fiber.Ctx) error {
	userId := c.Locals("userId").(uuid.UUID)
	var req types.JoinProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	if err := ph.repo.JoinProject(c.Context(), repository.JoinProjectParams{Projectid: req.ProjectId, Userid: userId}); err != nil {
		return err
	}
	project, err := ph.repo.GetProjectByID(context.Background(), req.ProjectId)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"project": project})
}
