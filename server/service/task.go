package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"tg-backend/db"
	"tg-backend/db/model"
	"tg-backend/server/types"
	"tg-backend/server/util"
)

type TaskService struct {
	storage db.Storage
}

func NewTaskService(storage db.Storage) *TaskService {
	return &TaskService{storage: storage}
}

func (p *TaskService) GetTask(ctx context.Context, tgUser *types.TelegramUser) ([]types.Task, error) {
	tasks, err := p.storage.GetTasks(ctx)
	if err != nil {
		return nil, err
	}
	userTask, err := p.storage.UserTask(ctx, tgUser.ID)
	if err != nil {
		return nil, err
	}
	return types.GetTasks(tasks, userTask), nil
}

func (p *TaskService) TaskCheck(ctx context.Context, tgUser *types.TelegramUser, taskId uint64) (bool, error) {
	task, err := p.storage.GetTask(ctx, taskId)
	if err != nil {
		return false, err
	}
	if task == nil || task.Status == model.DELETE {
		return false, errors.New("task not found")
	}
	exist, err := p.storage.GetUserStatus(ctx, taskId, tgUser.ID)
	if err != nil {
		return false, err
	}
	if exist {
		return false, nil
	}
	switch task.TaskType {
	case model.TaskTypeFriend:
		fs, err := p.storage.GetFriends(ctx, tgUser.ID)
		if err != nil {
			return false, err
		}
		if len(fs) > task.Condition {
			err := p.storage.TaskDone(ctx, task, tgUser.ID)
			if err != nil {
				return false, err
			}
			return true, nil
		}
	case model.TaskTypeX:
		err := p.storage.TaskDone(ctx, task, tgUser.ID)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

// GetTasksHandler
// @Tags 任务
// @Summary 获取任务信息
// @Accept json
// @Produce json
// @Failure 500 {string} string
// @Success 200 {array} types.Task
// @Router /api/task/get [get]
func GetTasksHandler(taskService *TaskService) util.HttpHandler {
	return func(r *http.Request) util.HandleResult {
		tgUser, _ := util.FromContext(r.Context())
		tasks, err := taskService.GetTask(r.Context(), tgUser)
		if err != nil {
			return util.ErrorWith("Error Get Tasks", util.ErrorInternal, err)
		}
		return util.Success(tasks)
	}
}

// TaskCheckHandler
// @Tags 任务
// @Summary 检查任务是否完成
// @Accept json
// @Produce json
// @Param  param body types.TaskCheck true "任务ID"
// @Failure 500 {string} string
// @Success 200 {string} string "true/false"
// @Router /api/task/check [post]
func TaskCheckHandler(taskService *TaskService) util.HttpHandler {
	return func(r *http.Request) util.HandleResult {
		tgUser, _ := util.FromContext(r.Context())
		request := &types.TaskCheck{}
		if err := json.NewDecoder(r.Body).Decode(request); err != nil {
			return util.Error("payload error", util.ErrorBadData)
		}

		res, err := taskService.TaskCheck(r.Context(), tgUser, request.Id)
		if err != nil {
			return util.ErrorWith("Error check task", util.ErrorInternal, err)
		}
		return util.Success(res)
	}
}
