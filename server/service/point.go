package service

import (
	"context"
	"encoding/json"
	"net/http"
	"tg-backend/db"
	"tg-backend/db/model"
	"tg-backend/server/types"
	"tg-backend/server/util"
)

type PointService struct {
	storage db.PointStorage
}

func NewPointService(storage db.PointStorage) *PointService {
	return &PointService{storage: storage}
}

func (p *PointService) GetPoint(ctx context.Context, tgUser *types.TelegramUser) (*model.Point, error) {
	point, err := p.storage.GetPoint(ctx, uint64(tgUser.ID))
	if err != nil {
		return nil, err
	}
	return point, nil
}

func (p *PointService) UpdatePoint(ctx context.Context, point *model.Point) error {
	return p.storage.UpdatePoint(ctx, point)
}

func GetPointHandler(pointService *PointService) util.HttpHandler {
	return func(r *http.Request) util.HandleResult {
		tgUser, _ := util.FromContext(r.Context())
		point, err := pointService.GetPoint(r.Context(), tgUser)
		if err != nil {
			return util.ErrorWith("Error Get Point", util.ErrorInternal, err)
		}
		if point == nil {
			return util.Success(&model.Point{
				Id: uint64(tgUser.ID),
			})
		}
		point.Id = uint64(tgUser.ID)
		return util.Success(types.GetPoint(point))
	}
}

func UpdatePointHandler(pointService *PointService) util.HttpHandler {
	return func(r *http.Request) util.HandleResult {
		tgUser, _ := util.FromContext(r.Context())
		request := types.Point{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return util.Error("payload error", util.ErrorBadData)
		}

		err := pointService.UpdatePoint(r.Context(), request.GetPoint(uint64(tgUser.ID)))
		if err != nil {
			return util.ErrorWith("Error update Point", util.ErrorInternal, err)
		}
		return util.HandleResult{}
	}
}
