package types

import (
	"time"

	"github.com/Iyed-M/teamup-backend/internal/repository"
	"github.com/google/uuid"
)

type CreateProjectRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}
type CreateProjectResponse = repository.Project

type CreateTaskRequest struct {
	Content       string     `json:"content"`
	Deadline      *time.Time `json:"deadline"`
	AttachmentUrl *string    `json:"attachmentUrl"`
	TaskOrder     int32      `json:"taskOrder"`
}
type CreateTaskResponse struct {
	AssignementId uuid.UUID `json:"assignementId"`
	TaskId        uuid.UUID `json:"taskId"`
}

type TaskWithUsers struct {
	repository.Task
	Users []repository.User `json:"users"`
}
type GetProjectTasksResponse struct {
	Tasks    []TaskWithUsers      `json:"tasks,omitempty"`
	SubTasks []repository.SubTask `json:"subTasks,omitempty"`
}

type InviteTeamMemberRequest struct {
	UserId    uuid.UUID `json:"userId"`
	ProejctId uuid.UUID `json:"projectId"`
}

type JoinProjectRequest struct {
	ProjectId uuid.UUID `json:"projectId"`
}
