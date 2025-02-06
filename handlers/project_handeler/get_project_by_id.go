package project_handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func (ph *projectHandler) GetProjectByID(c *fiber.Ctx) error {
	projectId, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		log.Errorw("Invalid projectId in params", "err", err)
		return fiber.ErrBadRequest
	}
	project, err := ph.repo.GetProjectByID(c.Context(), projectId)
	if err != nil {
		return err
	}
	return c.JSON(project)
}
