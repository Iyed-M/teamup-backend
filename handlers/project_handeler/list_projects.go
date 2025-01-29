package project_handler

import (
	"fmt"

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
	fmt.Println("LOGGED", "user_id", id)
	fmt.Println("jwt context id", "user_id", id)
	log.Errorw("LOGGED", "user_id", id)
	log.Tracew("jwt context id", "user_id", id)
	// if err != nil {
	// 	return fmt.Errorf("no access to user_id in handler %v", err)
	// }
	rows, err := h.repo.ListProjects(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(rows)
}
