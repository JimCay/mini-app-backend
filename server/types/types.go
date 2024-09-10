package types

import (
	"tg-backend/db/model"
	"time"
)

type LoginInfo struct {
	InviteCode string `json:"invite_code"`
}

type Invite struct {
	Code string `json:"code"`
}

type Point struct {
	Id     uint64 `json:"id"`
	Value  uint64 `json:"value"`
	Rate   int32  `json:"rate"`
	Limit  int32  `json:"limit"`
	Energy int32  `json:"energy"`
}

func GetPoint(point *model.Point) *Point {
	return &Point{
		Id:     point.Id,
		Value:  point.Value,
		Rate:   point.Rate,
		Energy: point.GetEnergy(),
		Limit:  point.Limit,
	}
}

func (p *Point) GetPoint() *model.Point {
	return &model.Point{
		Id:        p.Id,
		Value:     p.Value,
		Rate:      p.Rate,
		Energy:    p.Energy,
		UpdatedAt: time.Now(),
	}
}

type Friend struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	UserName string `json:"userName"`
	Reward   int32  `json:"reward"`
}

func GetFriend(mi *model.MyInvitee) Friend {
	return Friend{
		Id:       mi.Id,
		Name:     mi.FirstName + " " + mi.LastName,
		UserName: mi.UserName,
		Reward:   mi.Reward,
	}
}

type Rank struct {
	Self  UserPoint   `json:"self"`
	Ranks []UserPoint `json:"ranks"`
}

type UserPoint struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	UserName string `json:"userName"`
	Pos      int32  `json:"pos"`
	Value    uint64 `json:"value"`
}

func GetRanks(ranks []model.Rank) []UserPoint {
	res := make([]UserPoint, 0, len(ranks))
	for index, rank := range ranks {
		res = append(res, UserPoint{
			Id:       rank.Id,
			UserName: rank.UserName,
			Name:     rank.FirstName + " " + rank.LastName,
			Pos:      int32(index) + 1,
			Value:    rank.Value,
		})
	}
	return res
}

type Task struct {
	Id          uint64 `json:"id"`
	Name        string `json:"name"`        // 任务名称
	Description string `json:"description"` // 任务细节
	Type        int    `json:"type"`        // 任务类型 1: 邀请任务 2：X任务
	Reward      int32  `json:"reward"`      // 奖励
	Status      bool   `json:"status"`      // 任务完成状态 false 未完成 true 已完成
}

func GetTasks(task []model.Task, userTask []model.UserTask) []Task {
	userDone := make(map[uint64]bool, len(userTask))
	if userTask != nil && len(userTask) > 0 {
		for _, ut := range userTask {
			userDone[ut.TaskId] = true
		}
	}
	res := make([]Task, 0, len(task))
	for _, item := range task {
		res = append(res, Task{
			Id:          item.Id,
			Name:        item.Name,
			Description: item.Description,
			Reward:      item.Reward,
			Type:        item.TaskType,
			Status:      userDone[item.Id],
		})
	}
	return res
}

type TaskCheck struct {
	Id     uint64 `json:"id"`
	Status bool   `json:"status"`
}
