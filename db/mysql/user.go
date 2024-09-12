package mysql

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"tg-backend/config"
	"tg-backend/db/model"
	"time"
)

func (s *MysqlStorage) InsertUser(ctx context.Context, user *model.User, inviteId uint64) error {
	inviteUser, err := s.GetUser(ctx, inviteId)
	if err != nil {
		return err
	}
	if inviteUser == nil {
		inviteId = 0
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(user).Error
		if err != nil {
			return err
		}
		point := &model.Point{
			Id:        user.Id,
			Value:     0,
			Limit:     config.DefaultDayLimit,
			Energy:    config.DefaultDayLimit,
			UpdatedAt: time.Now(),
		}
		if inviteId != 0 {
			friend := &model.Friend{
				Invitor:      inviteId,
				Invitee:      user.Id,
				Value:        int32(config.InvitorPoint),
				InviteeValue: int32(config.InviteePoint),
				CreatedAt:    time.Now(),
			}
			err = tx.Create(friend).Error
			if err != nil {
				return err
			}
			invitorPoint := &model.Point{}
			err = tx.Model(invitorPoint).Where("id = ?", inviteId).First(invitorPoint).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					invitorPoint.Id = inviteId
					invitorPoint.Energy = config.DefaultDayLimit
					invitorPoint.Limit = config.DefaultDayLimit
					err = tx.Create(invitorPoint).Error
				} else {
					return err
				}
			}
			invitorPoint.Value += config.InvitorPoint
			err = tx.Model(invitorPoint).Update("value", invitorPoint.Value).Error
			if err != nil {
				return err
			}
			point.Value = config.InviteePoint
		}

		err = tx.Clauses(clause.OnConflict{DoNothing: true}).Create(point).Error

		if err != nil {
			return err
		}
		return nil
	})

}

func (s *MysqlStorage) GetUser(ctx context.Context, id uint64) (*model.User, error) {
	user := &model.User{}
	err := s.db.Where("id = ?", id).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (s *MysqlStorage) UpdateUserDays(ctx context.Context, user *model.User) error {
	return s.db.Model(user).Select("login_days", "updated_at").Updates(*user).Error
}

func (s *MysqlStorage) GetFriends(ctx context.Context, id uint64) ([]model.MyInvitee, error) {
	var friends []model.MyInvitee
	err := s.db.Raw("SELECT u.id,u.user_name ,u.first_name,u.last_name,f.value as reward "+
		"FROM friend f left join user u on u.id = f.invitee  WHERE f.invitor = ?", id).Scan(&friends).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return friends, nil
}
