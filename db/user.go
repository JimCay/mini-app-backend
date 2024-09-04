package db

import (
	"context"
	"gorm.io/gorm"
	"tg-backend/config"
	"tg-backend/db/model"
	"time"
)

func (s *Storage) InsertUser(ctx context.Context, user *model.User, inviteId uint64) error {
	inviteUser, err := s.GetUser(ctx, inviteId)
	if err != nil {
		return err
	}
	if inviteUser == nil {
		inviteId = 0
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		err := s.db.Create(user).Error
		if err != nil {
			return err
		}
		point := &model.Point{
			Id:        user.Id,
			Value:     0,
			UpdatedAt: time.Now(),
		}
		if inviteId != 0 {
			friend := &model.Friends{
				Invitor:   inviteId,
				Invitee:   user.Id,
				Value:     config.INVITOR_POINT,
				CreatedAt: time.Now(),
			}
			err = s.db.Create(friend).Error
			if err != nil {
				return err
			}
			invitorPoint := &model.Point{}
			err = s.db.Where("id = ?", user.Id).First(invitorPoint).Error
			if err != nil {
				return err
			}
			invitorPoint.Value += config.INVITOR_POINT
			err = s.db.Model(invitorPoint).Update("point", invitorPoint.Value).Error
			if err != nil {
				return err
			}
			point.Value = config.INVITEE_POINT
		}
		err = s.db.Create(point).Error
		if err != nil {
			return err
		}
		return nil
	})

}

func (s *Storage) GetUser(ctx context.Context, id uint64) (*model.User, error) {
	user := &model.User{}
	err := s.db.Where("id = ?", id).First(user).Error
	if err != nil && err.Error() != RECORD_NOT_FOUND {
		return nil, err
	}
	return user, nil
}

func (s *Storage) GetFriends(ctx context.Context, id uint64) ([]model.Friends, error) {
	var friends []model.Friends
	err := s.db.Where("invitor", id).Find(&friends).Error
	if err != nil && err.Error() != RECORD_NOT_FOUND {
		return nil, err
	}
	return friends, nil
}

func (s *Storage) createUser(ctx context.Context, user *model.User) error {
	err := s.db.Create(user).Error
	return err
}
