package project_handler

import (
	project_service "github.com/Iyed-M/teamup-backend/service/project"
	"github.com/Iyed-M/teamup-backend/types"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func (ph *projectHandler) JoinProject(c *fiber.Ctx) error {
	log.Info("joinp")
	userId := c.Locals("userId").(uuid.UUID)
	var req types.JoinProjectRequest

	if err := c.BodyParser(&req); err != nil {
		log.Errorw("bad req", "err", err)
		return fiber.ErrBadRequest
	}

	return project_service.JoinProject(c.Context(), ph.conn, userId, req.ProjectId)
}
