package mysql

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"tg-backend/db/model"
)

func (s *MysqlStorage) GetPoint(ctx context.Context, id uint64) (*model.Point, error) {
	point := &model.Point{}
	err := s.db.Where("id = ?", id).First(point).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return point, nil
}

func (s *MysqlStorage) UpdatePoint(ctx context.Context, point *model.Point) error {
	return s.db.Model(point).Select("value", "energy", "updated_at").Updates(*point).Error
}

func (s *MysqlStorage) Ranks(ctx context.Context) ([]model.Rank, error) {
	var rating []model.Rank
	err := s.db.Raw("SELECT u.id,u.user_name ,u.first_name,u.last_name,p.value " +
		"FROM point p left join user u on u.id = p.id order by value desc limit 100 ").Scan(&rating).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return rating, nil
}

func (s *MysqlStorage) MyRank(ctx context.Context, id uint64) (int64, error) {
	var count int64
	err := s.db.Model(&model.Point{}).Where("value > (select value from point where id = ?)", id).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
