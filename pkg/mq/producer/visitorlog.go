package producer

import (
	"git.dustess.com/mk-base/kafka-driver/producer"
	"git.dustess.com/mk-base/log"
	"git.dustess.com/mk-training/mk-blog-svc/config"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/user/model"
	"sync"
	"time"
)

var _once sync.Once

var defaultProduce  *producer.Producer

func newProducer()  {
	_once.Do(func() {
		defaultProduce = producer.NewProducer(config.Get().Kafka.Addrs, config.Get().Kafka.Topic.MKWat.Name)
	})
}

// SendVisitorLog 生产浏览日志
func SendVisitorLog(uid, articleID, clientIP string)  {
	newProducer()
	var visitor = model.Visitor{
		UserID:    uid,
		Client:    clientIP,
		ArticleID: articleID,
		ViewAt:    time.Now().Unix(),
	}
	_, _, _err := defaultProduce.ProduceOneJSON(visitor)
	if _err != nil {
		log.Info("kafka error: ", _err)
	}
}