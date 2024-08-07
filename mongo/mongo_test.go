package mongo_test

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"log"
	"reflect"
	"std-library/mongo"
	"testing"
	"time"
)

func TestInitNew(t *testing.T) {
	type args struct {
		opts []*mongo.Opt
	}
	tests := []struct {
		name string
		args args
	}{{
		"Init",
		args{opts: []*mongo.Opt{
			{
				AliasName:         "",
				Uri:               "mongodb://127.0.0.1:27017/",
				SkipTLSVerify:     false,
				MaxPoolSize:       10,
				MinPoolSize:       1,
				HeartbeatInterval: 0,
				MaxConnecting:     0,
				MaxConnIdleTime:   0,
				PoolMonitor: &event.PoolMonitor{Event: func(poolEvent *event.PoolEvent) {
					//log.Println("1->", poolEvent.ConnectionID)
				}},
				SocketTimeout: 0,
			},
			{
				AliasName:         "testA",
				Uri:               "mongodb://127.0.0.1:27017/",
				SkipTLSVerify:     false,
				MaxPoolSize:       20,
				MinPoolSize:       2,
				HeartbeatInterval: 0,
				MaxConnecting:     0,
				MaxConnIdleTime:   0,
				PoolMonitor: &event.PoolMonitor{Event: func(poolEvent *event.PoolEvent) {
					//log.Println("2->", poolEvent.ConnectionID)
				}},
				SocketTimeout: 0,
			},
		}},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mongo.Init(tt.args.opts...)
			doc := Doc{
				A: 11,
				B: "AAAA",
				C: 1.2356,
				D: true,
				E: time.Now(),
			}

			_, err := mongo.DB().InsertOne("testDB", "test", &doc)
			if err != nil {
				log.Fatal(err)
			}
			doc = Doc{
				A: 22,
				B: "BBBB",
				C: 1254.5887,
				D: false,
				E: time.Now().Add(time.Hour),
			}
			_, err = mongo.DB("testA").InsertOne("testDB", "test", &doc)
			if err != nil {
				log.Fatal(err)
			}
			doc = Doc{}
			err = mongo.DB().FindOne("testDB", "test", bson.M{"A": 11}).Decode(&doc)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(json.Marshal(doc))
			doc = Doc{}
			err = mongo.DB("testA").FindOne("testDB", "test", bson.M{"A": 22}).Decode(&doc)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(json.Marshal(doc))
			doc = Doc{
				A: 22,
				B: "BBBB",
				C: 1254.5887,
				D: false,
				E: time.Now().Add(time.Hour),
			}
			_, err = mongo.DB("testA").InsertOne("testDB", "test", &doc)
			if err != nil {
				log.Fatal(err)
			}
			err = mongo.DB("testA").UseSession(context.TODO(), nil, func(sCtx *mongo.SessionContext) error {
				err := sCtx.Begin(nil)
				if err != nil {
					return err
				}
				defer sCtx.End()
				_, err = mongo.DB("testA").WithCtx(sCtx.Context()).InsertOne("testDB", "test", doc)
				if err != nil {
					sCtx.Rollback()
					return err
				}
				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
			db := mongo.DB("testA")
			err = db.UseSession(context.TODO(), nil, func(sCtx *mongo.SessionContext) error {
				err := sCtx.Begin(nil)
				if err != nil {
					return err
				}
				defer sCtx.End()
				_, err = db.WithCtx(sCtx.Context()).InsertOne("testDB", "test", doc)
				if err != nil {
					sCtx.Rollback()
					return err
				}
				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
		})
	}
}

type Doc struct {
	A int       `bson:"A"`
	B string    `bson:"B"`
	C float32   `bson:"C"`
	D bool      `bson:"D"`
	E time.Time `bson:"E"`
}

func TestFilter(t *testing.T) {
	//mongo.Init(tt.args.opts...)
	cryptoOpt := `{"dbA":{"collA":["fieldA1","fieldA2"], "collB":["fieldB1", "fieldB2"]},"dbB":{"collC":["fieldC1","fieldC2"]}}`
	var m map[string]map[string][]string
	err := json.Unmarshal([]byte(cryptoOpt), &m)
	if err != nil {
		t.Fatal(err)
	}
	mongo.SetCryptoMap(m)
	mongo.SetPreExec(func(db, coll string, i ...any) error {
		log.Println(db, coll)
		log.Println(reflect.ValueOf(i).Type().String())
		log.Println(mongo.GetCrypto(db, coll))
		return nil
	})
	mongo.SetAfterExec(func(db, coll string, i ...any) error {
		log.Println(db, coll)
		log.Println(reflect.ValueOf(i).Type().String())
		log.Println(mongo.GetCrypto(db, coll))
		return nil
	})
	//Do Insert/Update/Find...
}
