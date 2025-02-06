package project_handler

import (
	"github.com/Iyed-M/teamup-backend/internal/repository"
	project_service "github.com/Iyed-M/teamup-backend/service/project"
	"github.com/Iyed-M/teamup-backend/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (ph *projectHandler) CreateTask(c *fiber.Ctx) error {
	projectId, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		return fiber.ErrBadRequest
	}
	userId := c.Locals("userId").(uuid.UUID)
	var req types.CreateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	project_service.NewTask(c.Context(), ph.conn, repository.CreateTaskParams{
		Content:       req.Content,
		ProjectID:     projectId,
		Deadline:      req.Deadline,
		AttachmentUrl: req.AttachmentUrl,
		TaskOrder:     req.TaskOrder,
	}, userId)
	return c.JSON(types.CreateTaskResponse{AssignementId: uuid.Nil, TaskId: uuid.Nil})
}
