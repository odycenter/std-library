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

type Cursor struct {
	*mongo.Cursor
	db   string
	coll string
}

func (c *Cursor) Decode(v any) error {
	err := c.Cursor.Decode(v)
	if err != nil {
		return err
	}
	return globalOpts.getAfterExec()(c.db, c.coll, v)
}

/*return result*/

type UpdateResult struct {
	*mongo.UpdateResult
}

type SingleResult struct {
	*mongo.SingleResult
	db   string
	coll string
}

func (s *SingleResult) Decode(v any) error {
	err := s.SingleResult.Decode(v)
	if err != nil {
		return err
	}
	return globalOpts.getAfterExec()(s.db, s.coll, v)
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
	db   string
	coll string
}

func (c *ChangeStream) Decode(v any) error {
	err := c.ChangeStream.Decode(v)
	if err != nil {
		return err
	}
	return globalOpts.getAfterExec()(c.db, c.coll, v)
}

type Pipeline mongo.Pipeline

type GlobalOptions struct {
	preExec   FilterFunc //执行前操作
	afterExec FilterFunc //执行后操作
}

func (g *GlobalOptions) setPreExec(f FilterFunc) {
	g.preExec = f
}

func (g *GlobalOptions) getPreExec() FilterFunc {
	if g.preExec == nil {
		return defaultFilterFunc
	}
	return g.preExec
}

func (g *GlobalOptions) setAfterExec(f FilterFunc) {
	g.afterExec = f
}

func (g *GlobalOptions) getAfterExec() FilterFunc {
	if g.afterExec == nil {
		return defaultFilterFunc
	}
	return g.afterExec
}

type FilterFunc func(db, coll string, i ...any) error

var defaultFilterFunc = func(db, coll string, i ...any) error {
	return nil
}
