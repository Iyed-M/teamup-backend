package project_handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (ph *projectHandler) ListInvitations(c *fiber.Ctx) error {
	userId := c.Locals("userId").(uuid.UUID)
	invitations, err := ph.repo.GetProjectInvitations(c.Context(), userId)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"invitations": invitations})
}
