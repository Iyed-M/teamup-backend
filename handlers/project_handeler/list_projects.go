package project_handler

import (
	"github.com/Iyed-M/teamup-backend/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

type TaskData struct {
	repository.Task
	SubTasks []struct {
		TaskData
	} `json:"subTasks,omitempty"`
	Assignments []struct {
		repository.TaskAssignment
	} `json:"assignments,omitempty"`
}
type ProjectData struct {
	repository.Project
	Tasks []TaskData `json:"tasks,omitempty"`
}

type listProjectResponse []ProjectData

func (h *projectHandler) ListProjects(c *fiber.Ctx) error {
	id := c.Locals("userId").(uuid.UUID)
	rows, err := h.repo.ListProjects(c.Context(), id)
	if err != nil {
		log.Errorw("db", "err", err)
		return err
	}
	return c.JSON(rows)
}
