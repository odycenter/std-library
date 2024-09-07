package elastic_test

import (
	"context"
	"fmt"
	"github.com/odycenter/std-library/elastic"
	"testing"
)

func TestElastic(t *testing.T) {
	elastic.NewClient(&elastic.Option{
		URLs:     []string{"https://127.0.0.1:9200"},
		Scheme:   "https",
		Username: "elastic",
		Password: "BsdegDSXpAYXah+j2SZk",
		Sniff:    false,
		SkipTLS:  true,
	})
	// 更新索引
	mapping := `{
		"properties":{
			"CreateTime":{
				"type":"integer"
			}
		}
	}`
	fmt.Println(elastic.Cli().Version("https://127.0.0.1:9200"))
	fmt.Println(elastic.Cli().CreateIndex(context.Background(), "testIndex1", mapping))
	fmt.Println(elastic.Cli().IndexExists(context.Background(), "testIndex1"))
	fmt.Println(elastic.Cli().SearchSimple(context.Background(), "testIndex1", "5f7f7ce161d14b5d985eb1e6b832488e"))
	fmt.Println(elastic.Cli().SearchSimpleWithBody(context.Background(), &elastic.SimpleSearchQuery{
		Index:  []string{"testIndex1"},
		IDs:    []string{"5f7f7ce161d14b5d985eb1e6b832488e"},
		From:   elastic.Int(0),
		Size:   elastic.Int(500),
		Pretty: elastic.Bool(true),
		Sort: []elastic.Sort{{
			Field:     "Time",
			Ascending: false,
		}},
	}))
	fmt.Println(elastic.Cli().SearchAdvance(context.Background(), &elastic.SearchQuery{
		Index:  []string{"testIndex1"},
		From:   elastic.Int(0),
		Size:   elastic.Int(500),
		Pretty: elastic.Bool(true),
		Sort: []elastic.Sort{{
			Field:     "Time",
			Ascending: false,
		}},
		QueryEntity: []elastic.QueryEntity{{
			Key: elastic.Range,
			Values: []elastic.QueryBody{
				{
					Key: "FieldB",
					Gte: 1,
					Lte: 87,
				},
				{
					Key: "FieldC",
					Gt:  1,
					Lt:  5,
				},
			},
		}},
	}))
	fmt.Println(elastic.Cli().SearchAdvance(context.Background(), &elastic.SearchQuery{
		Index:  []string{"testIndex1"},
		From:   elastic.Int(0),
		Size:   elastic.Int(500),
		Pretty: elastic.Bool(true),
		Sort: []elastic.Sort{{
			Field:     "Time",
			Ascending: false,
		}},
		QueryEntity: []elastic.QueryEntity{{
			Key: elastic.Term,
			Values: []elastic.QueryBody{
				{
					Key:   "FieldA",
					Value: "TestValue1",
				},
			},
		}},
	}))
}
