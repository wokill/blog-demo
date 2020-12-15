package service

import (
	"context"
	"errors"
	"git.dustess.com/mk-base/log"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/blog/dao"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/blog/model"
	statDao "git.dustess.com/mk-training/mk-blog-svc/pkg/blogstatistics/dao"
	statModel "git.dustess.com/mk-training/mk-blog-svc/pkg/blogstatistics/model"
	tagDao "git.dustess.com/mk-training/mk-blog-svc/pkg/tags/dao"
	userDao "git.dustess.com/mk-training/mk-blog-svc/pkg/user/dao"
	userModel "git.dustess.com/mk-training/mk-blog-svc/pkg/user/model"
	"go.mongodb.org/mongo-driver/bson"
)

// BlogService 博客服务
type BlogService struct {
	dao *dao.BlogDao
	ctx context.Context
}

// NewBlogService [...]
func NewBlogService(ctx context.Context) *BlogService {
	return &BlogService{
		dao: dao.NewBlogDao(ctx),
		ctx: ctx,
	}
}

// CreateBlog 创建博客
func (b *BlogService) CreateBlog(data model.Blog, tags []string) (err error) {
	if len(tags) > 0 {
		td := tagDao.NewTagDao(b.ctx)
		go func() {
			_, _ = td.InsertTags(tags)
		}()
	}
	data.Tags = tags
	var id string
	id, err = b.dao.InsertOne(data)
	es := dao.NewEsDao()
	data.ID = id
	_ = es.CreateIndex(b.ctx, data)
	return
}

// FindBlog 查找单条数据
func (b *BlogService) FindBlog(id string) (model.Blog, error) {
	return b.dao.FindByID(id)
}

// UpdateBlog 修改博客
func (b *BlogService) UpdateBlog(session userModel.UserSession, id string, blog model.Blog) error {
	article, err := b.FindBlog(id)
	if err != nil {
		return err
	}
	if article.Author != session.ID {
		return errors.New("无权限修改当前用户的文章")
	}
	if len(blog.Tags) > 0 {
		td := tagDao.NewTagDao(b.ctx)
		go func() {
			_, _ = td.InsertTags(blog.Tags)
		}()
	}
	article.Title = blog.Title
	article.Content = blog.Content
	article.Tags = blog.Tags
	err = b.dao.UpdateBlog(id, article)
	es := dao.NewEsDao()
	err = es.UpdateIndex(b.ctx, article)
	return err
}

// Delete 删除
func (b *BlogService) Delete(session userModel.UserSession, id string) error {
	article, err := b.FindBlog(id)
	if err != nil {
		return err
	}
	if article.Author != session.ID {
		return errors.New("无权限修改当前用户的文章")
	}
	err = b.dao.DeleteByID(id)
	es := dao.NewEsDao()
	_ = es.DeleteBlog(b.ctx, id)
	return err
}

// Lists 列表
func (b *BlogService) Lists(session userModel.UserSession, offset, limit int64) (resp []model.BlogList, count int64) {
	resp = make([]model.BlogList, 0)
	filter := bson.M{"author": session.ID}
	result, err := b.dao.FindAll(offset, limit, filter)
	if err != nil {
		return
	}
	resp = b.HandleLists(result)
	count = b.dao.FindCount(filter)
	return
}

// HandleLists 处理list
func (b *BlogService) HandleLists(result []model.Blog) (resp []model.BlogList) {
	if len(result) < 1 {
		return
	}
	type _st struct {
		Pv int64
		Uv int64
	}
	var (
		bid     = make([]string, len(result))
		uids    = make(map[string]struct{}, 0)
		users   []string
		stMap   = make(map[string]_st)
		userMap = make(map[string]userModel.User)
	)
	for k, v := range result {
		resp = append(resp, model.BlogList{
			CreatedAt: v.CreatedAt,
			Title:     v.Title,
			ID:        v.ID,
			Author:    v.Author,
			Tags:      v.Tags,
		})
		bid[k] = v.ID
		uids[v.Author] = struct{}{}
	}
	stDao := statDao.NewBlogStatDao(b.ctx)
	query := stDao.BlogStatQuery(bid, 0, 0, "")
	stResult := stDao.AsynStaticsByBlog(query)
	for _, v := range stResult {
		var temp  _st
		if _st, ok := stMap[v.BlogID]; ok {
			temp  = _st
		}
		if v.Typ == statModel.UV {
			temp.Uv =v.View
		} else {
			temp.Pv = v.View
		}
		stMap[v.BlogID] = temp
	}

	for k := range uids {
		users = append(users, k)
	}
	ud := userDao.NewUserDao(b.ctx)
	userRe, _ := ud.FindManyByID(users)
	for _, v := range userRe {
		userMap[v.ID] = v
	}
	for k, v := range resp {
		if st, ok := stMap[v.ID]; ok {
			resp[k].Uv = st.Uv
			resp[k].Pv = st.Pv
		}
		if us, ok := userMap[v.Author]; ok {
			resp[k].Author = us.Name
		}
	}
	return
}

// AttrUser 添加用户数据
func (b *BlogService) AttrUser(result []model.BlogList)(resp []model.BlogList) {
	var (
		users   = make([]string, len(result))
		userMap = make(map[string]userModel.User)
	)
	for k, v := range result {
		users[k] = v.Author
	}
	ud := userDao.NewUserDao(b.ctx)
	userRe, _ := ud.FindManyByID(users)
	for _, v := range userRe {
		userMap[v.ID] = v
	}
	for k, v := range result {
		if us, ok := userMap[v.Author]; ok {
			result[k].Author = us.Name
		}
	}
	return result
}

// Detail 博客详情
func (b *BlogService) Detail(id string) ( model.BlogDetail ,error) {
	blog, err := b.FindBlog(id)
	if err != nil {
		return model.BlogDetail{}, err
	}
	bl := []model.Blog{ blog }
	lists := b.HandleLists(bl)
	resp := model.BlogDetail{
		BlogList: model.BlogList{
			Uv:        0,
			Pv:        0,
			CreatedAt: blog.CreatedAt,
			Title:     blog.Title,
			ID:        blog.ID,
			Tags:      blog.Tags,
		},
		Content:  blog.Content,
	}
	if len(lists) > 0 {
		resp.BlogList = lists[0]
	}
	return resp, nil
}

//VisitorLists 访客博客
func (b *BlogService) VisitorLists(search model.ListSearch) ([]model.BlogList, int64 ){
	var (
		result = make([]model.BlogList, 0)
		err error
		count int64
	)
	es := dao.NewEsDao()
	result,count, err = es.VisitorBlogSearch(b.ctx, search)
	if err != nil {
		log.Error("elastic search error: ", err)
		return result, 0
	}
	result = b.AttrUser(result)
	return result, count
}