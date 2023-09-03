package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var True = true
var False = false
var After = options.After
var UpsertOptions = options.UpdateOptions{Upsert: &True}
var AfterOptions = options.FindOneAndUpdateOptions{ReturnDocument: &After}
var CollSecPre = options.CollectionOptions{ReadPreference: readpref.SecondaryPreferred()}

// IndexModel indexModel 表示要创建的新索引。
type IndexModel struct {
	//描述应将哪些键用于索引的文档。
	//它不能为nil。 这必须是保序类型，例如 bson.D。
	//
	//bson.M 等地图类型无效。
	//
	//有关有效文档的示例，请参阅 https://www.mongodb.com/docs/manual/indexes/#indexes 。
	Keys any

	// 用于创建索引的设置。
	Options *options.IndexOptions
}

// ToIndexModel 转换为 mongo.IndexModel
func ToIndexModel(ms []IndexModel) (i []mongo.IndexModel) {
	for _, m := range ms {
		i = append(i, mongo.IndexModel{
			Keys:    m.Keys,
			Options: m.Options,
		})
	}
	return
}

type Cursor struct{ *mongo.Cursor }

/*return result*/

type UpdateResult struct {
	*mongo.UpdateResult
}

type SingleResult struct {
	*mongo.SingleResult
}

type InsertManyResult struct {
	*mongo.InsertManyResult
}

type InsertOneResult struct {
	*mongo.InsertOneResult
}

type DeleteResult struct {
	*mongo.DeleteResult
}

type ChangeStream struct {
	*mongo.ChangeStream
}

type Pipeline mongo.Pipeline
