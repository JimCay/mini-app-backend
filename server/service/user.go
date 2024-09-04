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
	friends, err := u.storage.GetFriends(ctx, uint64(tgUser.ID))
	if err != nil {
		return nil, err
	}
	res := make([]types.Friend, 0, len(friends))
	for _, friend := range friends {
		res = append(res, types.GetFriend(friend))
	}
	return res, nil
}

func LoginHandler(service *UserService) util.HttpHandler {
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
		return util.HandleResult{}
	}
}

func InviteHandler(service *UserService) util.HttpHandler {
	return func(r *http.Request) util.HandleResult {
		tgUser, _ := util.FromContext(r.Context())
		code, _ := util.EncodeInvite(uint64(tgUser.ID))
		return util.Success(&types.Invite{Code: code})
	}
}

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
