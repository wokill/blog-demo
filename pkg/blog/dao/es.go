package dao

import (
	"context"
	"fmt"
	"git.dustess.com/mk-base/es-driver/es"
	"git.dustess.com/mk-base/log"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/blog/model"
	statModel "git.dustess.com/mk-training/mk-blog-svc/pkg/blogstatistics/model"
	"github.com/olivere/elastic"
	"reflect"
	"sync"
	"time"
)

const (
	esIndex = "mk_blog"
	esType  = "blog"
)

var _once sync.Once

func createIndex() {
	_once.Do(func() {
		client := es.Client()
		exists, err := client.IndexExists(esIndex).Do(context.Background())
		if err != nil {
			panic(fmt.Errorf("es init fail! err: %s", err.Error()))
		}
		if !exists {
			mapping := `
{
	"mappings":{
		"blog":{
			"properties":{
				"id": {
					"type": "keyword"
				},
				"title":{
					"type":"text"
				},
				"content":{
					"type":"text"
				},
				"tags":{
					"type":"keyword"
				}
			}
		}
	}
}`
			createIndex, err := client.CreateIndex(esIndex).Body(mapping).Do(context.Background())
			if err != nil {
				panic(err)
			}
			if !createIndex.Acknowledged {
				log.Error("es 创建mapping error")
			}
		}
	})
}

type esDao struct {
	client *elastic.Client
}

// NewEsDao 初始化es
func NewEsDao() *esDao {
	createIndex()
	return &esDao{
		client: es.Client(),
	}
}

// CreateIndex 创建博客全文索引
func (c *esDao) CreateIndex(ctx context.Context, data model.Blog) error {
	data.CreatedAt = time.Now().Unix()
	result, err := c.client.Index().Index(esIndex).Type(esType).Id(data.ID).BodyJson(data).Do(ctx)
	if err != nil {
		return err
	}
	_ = result
	return nil
}

// UpdateIndex 更新博客全文索引
func (c *esDao) UpdateIndex(ctx context.Context, data model.Blog) error {
	updateData := map[string]interface{}{
		"title":   data.Title,
		"content": data.Content,
		"tags":    data.Tags,
	}
	result, err := c.client.Update().Index(esIndex).Type(esType).Id(data.ID).Doc(updateData).DocAsUpsert(true).Do(ctx)
	if err != nil {
		return err
	}
	_ = result
	return nil
}

//DeleteBlog 删除文章
func (c *esDao) DeleteBlog(ctx context.Context, id string) error {
	_, err := c.client.Delete().Index(esIndex).Type(esType).Id(id).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

// IncView 递增浏览数据
func (c *esDao) IncView(add []model.IncBlogStat) {
	if len(add) == 0 {
		return
	}
	var data []elastic.BulkableRequest
	for _, v := range add {
		var typ = statModel.PV
		if v.Type == statModel.UV {
			typ = statModel.UV
		}
		script := elastic.NewScript("ctx._source."+typ+" += params.num").Param("num", v.View)
		bulk := elastic.NewBulkUpdateRequest().Id(v.ID).Script(script)
		data = append(data, bulk)
	}
	_, err := c.client.Bulk().Index(esIndex).Type(esType).Add(data...).Do(context.Background())
	if err != nil {
		log.Errorf("update es statistics error: %s", err.Error())
		return
	}
}

// VisitorBlogSearch 搜索博客数据
func (c *esDao) VisitorBlogSearch(ctx context.Context, search model.ListSearch) ([]model.BlogList, int64, error) {
	bQuery := elastic.NewBoolQuery()
	var (
		q      []elastic.Query
		rq     *elastic.RangeQuery
		sort   = "createdAt"
		offset = 0
		size   = 20
		blog   model.BlogList
		resp   []model.BlogList
	)
	if len(search.Tags) > 0 {
		var temp = make([]interface{}, len(search.Tags))
		for k, v := range search.Tags {
			temp[k] = v
		}
		q = append(q, elastic.NewTermsQuery("tags", temp...))
	}
	if search.SearchKey != "" {
		mq := elastic.NewMultiMatchQuery(search.SearchKey, "title", "content")
		q = append(q, mq)
	}
	if search.CreateAt.Begin > 0 {
		rq = elastic.NewRangeQuery("createdAt").Gte(search.CreateAt.Begin)
	}
	if search.CreateAt.End > 0 {
		if rq == nil {
			rq = elastic.NewRangeQuery("createdAt")
		}
		rq = rq.Lte(search.CreateAt.End)
	}
	if rq != nil {
		q = append(q, rq)
	}

	// 排序值
	switch search.SortKey {
	case "pv":
		sort = "pv"
	case "uv":
		sort = "uv"
	}
	if len(q) > 0 {
		bQuery = bQuery.Must(q...)
	}
	if search.PageSize > 0 {
		size = search.PageSize
	}
	if search.Page > 0 {
		offset = (search.Page - 1) * size
	}
	do := c.client.Search().Index(esIndex)
	result, err := do.Query(bQuery).Sort(sort, false).From(offset).Size(size).Pretty(true).Do(ctx)
	if err != nil {
		return nil, 0, err
	}
	for _, item := range result.Each(reflect.TypeOf(blog)) {
		t := item.(model.BlogList)
		resp = append(resp, t)
	}

	return resp, result.Hits.TotalHits, nil
}
