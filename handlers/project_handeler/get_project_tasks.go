package project_handler

import (
	"github.com/Iyed-M/teamup-backend/internal/repository"
	"github.com/Iyed-M/teamup-backend/types"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func (ph *projectHandler) GetProjectTasks(c *fiber.Ctx) error {
	var response types.GetProjectTasksResponse
	projectId, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		log.Errorw("Invalid projectId in params", "err", err)
		return fiber.ErrBadRequest
	}

	tasks, err := ph.repo.ProjectTasks(c.Context(), projectId)
	if err != nil {
		return fiber.ErrNotFound
	}
	tasksMap := make(map[uuid.UUID]*types.TaskWithUsers)
	for _, t := range tasks {
		storedTask := tasksMap[t.Task.ID]

		if storedTask != nil {
			storedTask.Users = append(storedTask.Users, t.User)
		} else {
			tasksMap[t.Task.ID] = &types.TaskWithUsers{Task: t.Task, Users: []repository.User{t.User}}
		}
	}
	for _, task := range tasksMap {
		response.Tasks = append(response.Tasks, *task)
	}

	response.SubTasks, err = ph.repo.ProjectSubTasks(c.Context(), projectId)
	if err != nil {
		return fiber.ErrNotFound
	}
	return c.JSON(response)
}
