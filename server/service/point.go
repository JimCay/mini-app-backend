package service

import (
	"context"
	"encoding/json"
	"net/http"
	"tg-backend/db"
	"tg-backend/server/types"
	"tg-backend/server/util"
)

type PointService struct {
	storage db.PointStorage
}

func NewPointService(storage db.PointStorage) *PointService {
	return &PointService{storage: storage}
}

func (p *PointService) GetPoint(ctx context.Context, tgUser *types.TelegramUser) (*types.Point, error) {
	point, err := p.storage.GetPoint(ctx, tgUser.ID)
	if err != nil {
		return nil, err
	}
	if point == nil {
		return &types.Point{
			Id: tgUser.ID,
		}, nil
	}
	point.Id = tgUser.ID
	return types.GetPoint(point), nil
}

func (p *PointService) UpdatePoint(ctx context.Context, point *types.Point) error {
	return p.storage.UpdatePoint(ctx, point.GetPoint())
}

func (p *PointService) Rank(ctx context.Context, tgUser *types.TelegramUser) (*types.Rank, error) {
	point, err := p.storage.GetPoint(ctx, uint64(tgUser.ID))
	if err != nil {
		return nil, err
	}

	pos, err := p.storage.MyRank(ctx, uint64(tgUser.ID))
	if err != nil {
		return nil, err
	}

	self := tgUser.UserPoint(point.Value, int32(pos)+1)
	ranks, err := p.storage.Ranks(ctx)
	if err != nil {
		return nil, err
	}
	return &types.Rank{
		Self:  self,
		Ranks: types.GetRanks(ranks),
	}, nil
}

// GetPointHandler
// @Summary 获取积分信息
// @Tags 积分
// @Accept json
// @Produce json
// @Failure 500 {string} string
// @Success 200 {object} types.Point
// @Router /api/point/query [get]
func GetPointHandler(pointService *PointService) util.HttpHandler {
	return func(r *http.Request) util.HandleResult {
		tgUser, _ := util.FromContext(r.Context())
		point, err := pointService.GetPoint(r.Context(), tgUser)
		if err != nil {
			return util.ErrorWith("Error Get Point", util.ErrorInternal, err)
		}
		return util.Success(point)
	}
}

// GetRankHandler
// @Summary 获取积分排名
// @Tags 积分
// @Accept json
// @Produce json
// @Failure 500 {string} string
// @Success 200 {array} types.Rank
// @Router /api/point/rank [get]
func GetRankHandler(pointService *PointService) util.HttpHandler {
	return func(r *http.Request) util.HandleResult {
		tgUser, _ := util.FromContext(r.Context())
		rank, err := pointService.Rank(r.Context(), tgUser)
		if err != nil {
			return util.ErrorWith("Error Get Point", util.ErrorInternal, err)
		}
		return util.Success(rank)
	}
}

// UpdatePointHandler
// @Summary 获取积分排名
// @Tags 积分
// @Accept json
// @Produce json
// @Param  param body types.Point true "积分信息"
// @Failure 500 {string} string
// @Success 200 {bool} string "OK"
// @Router /api/point/update [post]
func UpdatePointHandler(pointService *PointService) util.HttpHandler {
	return func(r *http.Request) util.HandleResult {
		tgUser, _ := util.FromContext(r.Context())
		request := &types.Point{}
		if err := json.NewDecoder(r.Body).Decode(request); err != nil {
			return util.Error("payload error", util.ErrorBadData)
		}
		p, err := pointService.GetPoint(r.Context(), tgUser)
		if p != nil {
			if p.Value > request.Value {
				return util.Error("point shrink", util.ErrorBadData)
			}
		}
		request.Id = tgUser.ID
		err = pointService.UpdatePoint(r.Context(), request)
		if err != nil {
			return util.ErrorWith("Error update Point", util.ErrorInternal, err)
		}
		return util.Success("OK")
	}
}
