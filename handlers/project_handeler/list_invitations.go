package project_handler

import (
	"github.com/Iyed-M/teamup-backend/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (ph *projectHandler) ListInvitations(c *fiber.Ctx) error {
	userId := c.Locals("userId").(uuid.UUID)
	invitations, err := ph.repo.GetProjectInvitations(c.Context(), userId)
	if err != nil {
		return err
	}
	users := make(map[uuid.UUID]*repository.User)
	for _, inv := range invitations {
		if users[inv.User.ID] == nil {
			users[inv.User.ID] = &inv.User
		}
	}
	return c.JSON(fiber.Map{"invitations": invitations})
}
