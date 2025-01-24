package team_handler

import (
	"github.com/Iyed-M/teamup-backend/internal/repository"
	team_service "github.com/Iyed-M/teamup-backend/service/team"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type createTeamRequest struct {
	Name  string  `json:"name"`
	Color *string `json:"color,omitempty"`
}
type createTeamResponse = repository.Team

func (h *teamHandler) CreateTeam(c *fiber.Ctx) error {
	var req createTeamRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	// TODO: parse userId from accessToken (make jwt middleware)
	var userId uuid.UUID
	team, err := team_service.CreateTeam(c.Context(), h.conn, userId, req.Name, req.Color)
	if err != nil {
		return err
	}
	return c.JSON(team)
}
