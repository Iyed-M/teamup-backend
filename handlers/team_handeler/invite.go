package team_handler

import (
	"github.com/Iyed-M/teamup-backend/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type inviteTeamMemberRequest struct {
	UserId uuid.UUID `json:"userId"`
	TeamId uuid.UUID `json:"teamId"`
}

func (h *teamHandler) InviteTeamMember(c *fiber.Ctx) error {
	senderId := c.Locals("userId").(uuid.UUID)
	var req inviteTeamMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if err := h.repo.InviteToTeam(c.Context(), repository.InviteToTeamParams{Teamid: req.TeamId, Senderid: senderId, Receiverid: req.UserId}); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error inviting user to team")
	}
	return nil
}
