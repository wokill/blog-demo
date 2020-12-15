package service

import (
	"context"
	"errors"
	"git.dustess.com/mk-base/util/crypto"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/common"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/user/dao"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/user/model"
	"time"
)

var (
	password = errors.New("密码或用户名不正确")
)

//  userService 用户服务
type userService struct {
	ctx   context.Context
	dao   *dao.UserDao
	cache *common.Cache
}

// NewUserService [...]
func NewUserService(ctx context.Context) *userService {
	return &userService{
		ctx:   ctx,
		dao:   dao.NewUserDao(ctx),
		cache: common.NewCache(common.SessionKey, int64(common.SessionExpired)),
	}
}

// Login 登录
func (s *userService) Login(account, pass string) (token string, user model.User, err error) {
	user, err = s.dao.FindByAccount(account)
	if err != nil {
		err = password
		return
	}
	if !common.Verify(pass, user.Password) {
		err = password
		return
	}
	token = crypto.Md5(crypto.RandID())
	err = s.cache.CreateJSON(token, model.UserSession{
		ID:      user.ID,
		Account: user.Account,
		Token:   token,
		LoginAt: time.Now().Unix(),
	})
	m := dao.NewUserDao(s.ctx)
	_ = m.Login(user.ID)
	return
}

// Register 注册
func (s *userService) Register(account, name, pass string) (id string, err error) {
	now := time.Now().Unix()
	us, err := s.dao.FindByAccount(account)
	if err == nil && us.ID != "" {
		return "", errors.New("账号已存在")
	}
	data := model.User{
		Account:   account,
		Name:      name,
		Password:  common.Encrypt(pass),
		Avatar:    "",
		CreatedAt: now,
	}
	id, err = s.dao.InsertUser(data)
	return
}

// Logout 退出登录
func (s *userService) Logout(d model.UserSession) error {
	err := s.dao.Logout(d.ID)
	if err != nil {
		return err
	}

	_ = s.cache.DelSession(d.Token)
	return nil
}
