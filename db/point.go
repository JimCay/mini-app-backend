package db

import (
	"context"
	"tg-backend/db/model"
)

func (s *Storage) GetPoint(ctx context.Context, id uint64) (*model.Point, error) {
	point := &model.Point{}
	err := s.db.Where("id = ?", id).First(point).Error
	if err != nil {
		if err.Error() == RECORD_NOT_FOUND {
			return nil, nil
		}
		return nil, err
	}
	return point, nil
}

func (s *Storage) UpdatePoint(ctx context.Context, point *model.Point) error {
	return s.db.Model(point).Updates(*point).Error
}
