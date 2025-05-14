package elasticsearchapi

import "github.com/elastic/go-elasticsearch/v8"

type Settings struct {
	NameRegionalObject string
	Host               string
	User               string
	Passwd             string
	IndexDB            string
	Port               int
}

type ElasticsearchDB struct {
	client   *elasticsearch.Client
	settings Settings
}
