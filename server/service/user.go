package service

import (
	"context"
	"encoding/json"
	"net/http"
	"tg-backend/db"
	"tg-backend/server/types"
	"tg-backend/server/util"
)

type UserService struct {
	storage db.UserStorage
}

func NewUserService(storage db.UserStorage) *UserService {
	return &UserService{storage: storage}
}

func (u *UserService) Login(ctx context.Context, tgUser *types.TelegramUser, info types.LoginInfo) error {
	user, err := u.storage.GetUser(ctx, uint64(tgUser.ID))
	if err != nil {
		return err
	}
	if user == nil {
		user = tgUser.GetUser()
		inviteId := uint64(0)
		if info.InviteCode != "" {
			inviteId = util.DecodeInvite(info.InviteCode)
		}
		return u.storage.InsertUser(ctx, user, inviteId)
	}
	return nil
}

func (u *UserService) Friends(ctx context.Context, tgUser *types.TelegramUser) ([]types.Friend, error) {
	myInvitee, err := u.storage.GetFriends(ctx, uint64(tgUser.ID))
	if err != nil {
		return nil, err
	}
	res := make([]types.Friend, 0, len(myInvitee))
	for _, mi := range myInvitee {
		res = append(res, types.GetFriend(&mi))
	}
	return res, nil
}

// LoginHandler
// @Tags 用户
// @Summary 登录
// @Accept json
// @Produce json
// @Failure 500 {string} string
// @Success 200 {object} types.Point
// @Router /api/user/login [post]
func LoginHandler(service *UserService, pointService *PointService) util.HttpHandler {
	return func(r *http.Request) util.HandleResult {
		tgUser, _ := util.FromContext(r.Context())

		request := types.LoginInfo{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return util.Error("payload error", util.ErrorBadData)
		}
		err := service.Login(r.Context(), tgUser, request)
		if err != nil {
			return util.ErrorWith("Error login", util.ErrorInternal, err)
		}
		point, err := pointService.GetPoint(r.Context(), tgUser)
		if err != nil {
			return util.ErrorWith("Error Get Point", util.ErrorInternal, err)
		}
		return util.Success(point)
	}
}

// InviteHandler
// @Tags 用户
// @Summary 获取邀请码
// @Accept json
// @Produce json
// @Failure 500 {string} string
// @Success 200 {string} types.Invite
// @Router /api/user/invite [get]
func InviteHandler(service *UserService) util.HttpHandler {
	return func(r *http.Request) util.HandleResult {
		tgUser, _ := util.FromContext(r.Context())
		code, _ := util.EncodeInvite(uint64(tgUser.ID))
		return util.Success(&types.Invite{Code: code})
	}
}

// FriendHandler
// @Tags 用户
// @Summary 获取邀请的好友
// @Accept json
// @Produce json
// @Failure 500 {string} string
// @Success 200 {array} types.Friend
// @Router /api/user/friends [get]
func FriendHandler(service *UserService) util.HttpHandler {
	return func(r *http.Request) util.HandleResult {
		tgUser, _ := util.FromContext(r.Context())
		friends, err := service.Friends(r.Context(), tgUser)
		if err != nil {
			return util.ErrorWith("Error get friends", util.ErrorInternal, err)
		}
		return util.Success(friends)
	}
}
