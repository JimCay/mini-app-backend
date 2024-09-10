package mysql

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"tg-backend/db/model"
)

func (s *MysqlStorage) GetTasks(ctx context.Context) ([]model.Task, error) {
	var task []model.Task
	err := s.db.Where("status = ?", model.NORMAL).Find(&task).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return task, nil
}

func (s *MysqlStorage) GetTask(ctx context.Context, id uint64) (*model.Task, error) {
	task := &model.Task{}
	err := s.db.Where("id = ?", id).First(task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return task, nil
}

func (s *MysqlStorage) UserTask(ctx context.Context, userId uint64) ([]model.UserTask, error) {
	var task []model.UserTask
	err := s.db.Where("user_id = ?", userId).Find(&task).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return task, nil
}

func (s *MysqlStorage) TaskDone(ctx context.Context, task *model.Task, userId uint64) error {
	ut := &model.UserTask{
		TaskId: task.Id,
		UserId: userId,
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(ut).Error
		if err != nil {
			return err
		}
		userPoint := &model.Point{}
		err = tx.Model(userPoint).Where("id = ?", userId).First(userPoint).Error
		if err != nil {
			return err
		}
		userPoint.Value += uint64(task.Reward)
		userPoint.Energy = userPoint.GetEnergy()
		return tx.Model(userPoint).Select("value", "energy", "updated_at").Updates(userPoint).Error
	})
}

func (s *MysqlStorage) GetUserStatus(ctx context.Context, taskId, userId uint64) (bool, error) {
	userTask := &model.UserTask{}
	err := s.db.Where("task_id = ? AND user_id = ?", taskId, userId).Limit(1).Find(userTask).Error
	if err != nil {
		return false, err
	}
	return userTask.UserId == userId, nil
}
