package log

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/qingw1230/studyim/pkg/common/config"

	elasticV7 "github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

// esHook CUSTOMIZED ES hook
type esHook struct {
	moduleName string
	client     *elasticV7.Client
}

// newEsHook Initialization
func newEsHook(moduleName string) *esHook {
	es, err := elasticV7.NewClient(
		elasticV7.SetURL(config.Config.Log.ElasticSearchAddr...),
		elasticV7.SetBasicAuth(config.Config.Log.ElasticSearchUser, config.Config.Log.ElasticSearchPassword),
		elasticV7.SetSniff(false),
		elasticV7.SetHealthcheckInterval(60*time.Second),
		elasticV7.SetErrorLog(log.New(os.Stderr, "ES:", log.LstdFlags)),
	)

	if err != nil {
		log.Fatal("failed to create Elastic V7 Client: ", err)
	}
	return &esHook{client: es, moduleName: moduleName}
}

// Fire log hook interface
func (hook *esHook) Fire(entry *logrus.Entry) error {
	doc := newEsLog(entry)
	go hook.sendEs(doc)
	return nil
}

func (hook *esHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// sendEs
func (hook *esHook) sendEs(doc appLogDocModel) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("send entry to es failed: ", r)
		}
	}()
	_, err := hook.client.Index().Index(hook.moduleName).Type(doc.indexName()).BodyJson(doc).Do(context.Background())
	if err != nil {
		log.Println(err)
	}

}

// appLogDocModel es model
type appLogDocModel map[string]interface{}

func newEsLog(e *logrus.Entry) appLogDocModel {
	ins := make(map[string]interface{})
	ins["level"] = strings.ToUpper(e.Level.String())
	ins["time"] = e.Time.Format("2006-01-02 15:04:05")
	for kk, vv := range e.Data {
		ins[kk] = vv
	}
	ins["tipInfo"] = e.Message

	return ins
}

// indexName es index name
func (m *appLogDocModel) indexName() string {
	return time.Now().Format("2006-01-02")
}
