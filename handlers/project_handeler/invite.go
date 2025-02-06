package project_handler

import (
	"github.com/Iyed-M/teamup-backend/internal/repository"
	"github.com/Iyed-M/teamup-backend/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *projectHandler) InviteProjectMember(c *fiber.Ctx) error {
	senderId := c.Locals("userId").(uuid.UUID)
	var req types.InviteTeamMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if err := h.repo.InviteToProject(c.Context(), repository.InviteToProjectParams{Projectid: req.ProejctId, Senderid: senderId, Receiverid: req.UserId}); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error inviting user to team")
	}
	return nil
}
