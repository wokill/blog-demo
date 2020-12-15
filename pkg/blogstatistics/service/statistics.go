package service

import (
	"context"
	"errors"
	"git.dustess.com/mk-base/util/crypto"
	blogDao "git.dustess.com/mk-training/mk-blog-svc/pkg/blog/dao"
	blogModel "git.dustess.com/mk-training/mk-blog-svc/pkg/blog/model"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/blogstatistics/dao"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/blogstatistics/model"
	"time"
)

type stat struct {
	dao    *dao.BlogStatDao
	logDao *dao.BlogLogDao
}

// NewStat 生成
func NewStat() *stat {
	ctx := context.Background()
	return &stat{
		dao:    dao.NewBlogStatDao(ctx),
		logDao: dao.NewBlogLogDao(ctx),
	}
}

// StatisticsView 统计访问情况
func (s *stat) StatisticsView() error {

	blog := s.FindSourceData()
	if len(blog) == 0 {
		return errors.New("no statistics data")
	}
	now := time.Now().Unix()
	var (
		insertData []model.Statistics
		esView []blogModel.IncBlogStat
	)
	for _, v := range blog {
		temp := []model.Statistics{
			{
				ID:        crypto.UUID(),
				BlogID:    v.BlogID,
				Typ:       model.PV,
				Author:    v.Author,
				View:      v.Pv,
				CreatedAt: now,
			},
			{
				ID:        crypto.UUID(),
				BlogID:    v.BlogID,
				Typ:       model.UV,
				Author:    v.Author,
				View:      v.Uv,
				CreatedAt: now,
			},
		}
		insertData = append(insertData, temp...)
		// ES 浏览数据递增
		if v.Uv > 0 {
			esView = append(esView, blogModel.IncBlogStat{
				ID:   v.BlogID,
				Type: model.UV,
				View: v.Uv,
			})
		}
		if v.Pv > 0 {
			esView = append(esView, blogModel.IncBlogStat{
				ID:   v.BlogID,
				Type: model.PV,
				View: v.Pv,
			})
		}
	}
	err := s.dao.InsertMany(insertData)
	esModel := blogDao.NewEsDao()
	esModel.IncView(esView)
	return err
}

type blogSt struct {
	BlogID string
	Author string
	Pv     int
	Uv     int
}

// FindSourceData 通过原始数据得到 pv 与uv
func (s *stat) FindSourceData() map[string]blogSt {
	at := s.findLastStatAt(model.PV)
	logs, err := s.logDao.FindLogsByTime(at)
	if err != nil {
		return nil
	}
	var (
		userMap = make(map[string]map[string]struct{}, 0)
		client  = make(map[string]map[string]struct{}, 0)
		blog    = make(map[string]blogSt)
		author  []string
	)
	for _, v := range logs {
		var _st = blogSt{
			BlogID: v.BlogID,
			Pv:     1,
		}
		author = append(author, v.BlogID)
		if _b, ok := blog[v.BlogID]; ok {
			_st = _b
			_st.Pv++
		}
		blog[v.BlogID] = _st
		if v.UserID != "" {
			if _, ok := userMap[v.BlogID]; !ok {
				userMap[v.BlogID] = make(map[string]struct{})
			}
			userMap[v.BlogID][v.UserID] = struct{}{}
			continue
		}
		if _, ok := client[v.BlogID]; !ok {
			client[v.BlogID] = make(map[string]struct{})
		}
		client[v.BlogID][v.ClientIP] = struct{}{}
	}
	_blogAuthor := s.blogAuthor(author)
	for k, v := range blog {
		if _b, ok := _blogAuthor[k]; ok {
			v.Author = _b
			blog[k] = v
		}
	}
	for k, v := range userMap {
		var users []string
		for vv := range v {
			users = append(users, vv)
		}
		_uv := s.userUv(users, k, at)
		if _blog, ok := blog[k]; ok {
			_blog.Uv += _uv
			blog[k] = _blog
		}
	}
	for k, v := range client {
		var ips []string
		for vv := range v {
			ips = append(ips, vv)
		}
		_uv := s.ipUv(ips, k, at)
		if _blog, ok := blog[k]; ok {
			_blog.Uv += _uv
			blog[k] = _blog
		}
	}
	return blog
}

func (s *stat) userUv(users []string, blogID string, at int64) int {
	result, err := s.logDao.FindLogsByUserIDs(users, blogID, at)
	if err != nil {
		return 0
	}
	var _reUser = make(map[string]struct{})
	for _, v := range result {
		_reUser[v.UserID] = struct{}{}
	}
	return len(users) - len(_reUser)
}

func (s *stat) ipUv(ips []string, blogID string, at int64) int {
	result, err := s.logDao.FindLogsByIps(ips, blogID, at)
	if err != nil {
		return 0
	}
	var _reIp = make(map[string]struct{})
	for _, v := range result {
		_reIp[v.ClientIP] = struct{}{}
	}
	return len(ips) - len(_reIp)
}

// findLastStatAt 访问最后的统计时间
func (s *stat) findLastStatAt(typ string) int64 {
	last, err := s.dao.FindLast(typ)
	if err != nil {
		return 0
	}
	return last.CreatedAt
}

// blogAuthor 博客作者
func (s *stat) blogAuthor(blogID []string) map[string]string {
	var res = make(map[string]string)
	_dao := blogDao.NewBlogDao(context.Background())
	result, err := _dao.FindByIds(blogID)
	if err != nil {
		return res
	}
	for _, v := range result {
		res[v.ID] = v.Author
	}
	return res
}
