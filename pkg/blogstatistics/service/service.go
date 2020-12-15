package service

import (
	"context"
	"fmt"
	blogDao "git.dustess.com/mk-training/mk-blog-svc/pkg/blog/dao"
	blogModel "git.dustess.com/mk-training/mk-blog-svc/pkg/blog/model"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/blogstatistics/dao"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/blogstatistics/model"
	userDao "git.dustess.com/mk-training/mk-blog-svc/pkg/user/dao"
	userModel "git.dustess.com/mk-training/mk-blog-svc/pkg/user/model"
	"go.mongodb.org/mongo-driver/bson"
)

const size = 20

type BlogStatService struct {
	dao *dao.BlogStatDao
	ctx context.Context
}

func NewBlogService(ctx context.Context) *BlogStatService {
	return &BlogStatService{
		dao: dao.NewBlogStatDao(ctx),
		ctx: ctx,
	}
}

// BlogHot 博客热度
func (b *BlogStatService) BlogHot(sort string, start, end int64) (result []model.StateOverViewResp) {
	result = make([]model.StateOverViewResp, 0)
	base := b.StatisticsOverview(sort, start, end)
	var (
		blogsID = make([]string, len(base))
		usersID []string
		blogMap = make(map[string]blogModel.Blog)
		userMap = make(map[string]string)
		ud *userDao.UserDao
		users []userModel.User
	)
	for k, v := range base {
		blogsID[k] = v.BlogID
		var (
			uv, pv int64
		)

		for _, v := range v.Stat {
			if v.Typ == model.PV {
				pv = v.Pv
				continue
			}
			uv = v.Uv
		}
		result = append(result, model.StateOverViewResp{
			ID: v.BlogID,
			Pv: pv,
			Uv: uv,
		})
	}
	if len(blogsID) == 0 {
		return
	}
	bd := blogDao.NewBlogDao(b.ctx)
	blogs, _ := bd.FindByIds(blogsID)
	for _, v := range blogs {
		blogMap[v.ID] = v
		usersID = append(usersID, v.Author)
	}
	if len(usersID) == 0 {
		goto end
	}
	ud = userDao.NewUserDao(b.ctx)
	users, _ = ud.FindManyByID(usersID)
	for _, v := range users {
		userMap[v.ID] = v.Name
	}
end:
	for k, v := range result {
		if _blog, ok := blogMap[v.ID]; ok {
			v.Tags = _blog.Tags
			v.Title = _blog.Title
			v.Author = _blog.Author
		}
		if userName, ok := userMap[v.Author]; ok {
			v.Author = userName
		}
		result[k] = v
	}
	return
}

// StatisticsOverview 统计概览数据
func (b *BlogStatService) StatisticsOverview(sort string, start, end int64) []model.StatOverview {
	filter := bson.M{"createdAt": bson.M{"$gte": start, "$lt": end}}
	overview := b.dao.StatisticsByBlog(filter, sort, size)
	return overview
}

// BlogDetail 博客详情
func (b *BlogStatService) BlogDetail(blogID string,start, end int64) []model.StatDetailResp {
	var result = make([]model.StatDetailResp, 0)
	detail := b.dao.StaticsDetail(blogID, start, end)
	fmt.Println(detail, "###")
	for _, v := range detail {
		var (
			uv, pv int64
		)
		for _, v := range v.Stat {
			if v.Typ == model.PV {
				pv = v.Pv
				continue
			}
			uv = v.Uv
		}
		result = append(result, model.StatDetailResp{
			Date: v.Date,
			Pv: pv,
			Uv: uv,
		})
	}
	return  result
}