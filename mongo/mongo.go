// Package mongo mongo数据库操作
package mongo

import (
	"context"
	"crypto/tls"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	app "std-library/app/conf"
	"sync"
	"time"
)

var pool sync.Map                         //map[string]*mongo.Client
var CryptoMap *sync.Map                   //*sync.Map[string]*sync.Map[string]*[]string -> map[db]map[coll][fields...]
var globalOpts = new(GlobalOptions)       //全局操作
var ErrNoDocuments = mongo.ErrNoDocuments //mongo无记录错误，方便做error判断
var defaultTimeout = 180 * time.Second

// Opt 配置结构
type Opt struct {
	AliasName         string
	Uri               string
	SkipTLSVerify     bool
	ReadPreference    *readpref.ReadPref `json:"-"`
	MaxPoolSize       uint64
	MinPoolSize       uint64
	HeartbeatInterval time.Duration         `json:"-"`
	MaxConnecting     uint64                //default is 2. Values greater than 100 are not recommended
	MaxConnIdleTime   time.Duration         `json:"-"` //default is 0, meaning a connection can remain unused indefinitely.
	PoolMonitor       *event.PoolMonitor    `json:"-"` //mongo线程池监听
	CommandMonitor    *event.CommandMonitor `json:"-"` //mongo命令事件监听
	SocketTimeout     time.Duration         `json:"-"` //default is 0, meaning no timeout is used and socket operations can block indefinitely.
	Timeout           time.Duration         `json:"-"`
}

// Cli mongo客户端封装
type Cli struct {
	AliasName string
	db        *mongo.Client
	dbOpt     *options.DatabaseOptions
	collOpt   *options.CollectionOptions
	ctx       context.Context
}

// WithReadPreference 设置读取优先配置
func (opt *Opt) WithReadPreference(readPref *readpref.ReadPref) *Opt {
	opt.ReadPreference = readPref
	return opt
}

// WithPoolMonitor 设置连接池监控
func (opt *Opt) WithPoolMonitor(monitor *event.PoolMonitor) *Opt {
	opt.PoolMonitor = monitor
	return opt
}

// SetPreExec 设置前置拦截器
func SetPreExec(f FilterFunc) {
	globalOpts.setPreExec(f)
}

// SetAfterExec 设置后置拦截器
func SetAfterExec(f FilterFunc) {
	globalOpts.setAfterExec(f)
}

// SetCryptoMap 设置加密映射表
func SetCryptoMap(m map[string]map[string][]string) {
	if m == nil {
		return
	}
	dbM := sync.Map{}
	for dbK, dbV := range m {
		collM := sync.Map{}
		for collK, fields := range dbV {
			var tmp []string
			copy(tmp, fields)
			collM.Store(collK, tmp)
			//collM.Store(collK, slices.Clone(fields))
		}
		dbM.Store(dbK, &collM)
	}
	CryptoMap = &dbM
	return
}

func (opt *Opt) getAliasName() string {
	if opt.AliasName == "" {
		return "default"
	}
	return opt.AliasName
}

func (opt *Opt) getUri() string {
	if opt.Uri == "" {
		log.Panicln("[mongo]Mongo Uri is Empty")
	}
	return opt.Uri
}

func (opt *Opt) getMaxConnecting() uint64 {
	if opt.MaxConnecting == 0 {
		return 2
	}
	return opt.MaxConnecting
}

func GetCrypto(db, coll string) []string {
	if CryptoMap == nil {
		return nil
	}
	collM, ok := CryptoMap.Load(db)
	if !ok {
		return nil
	}
	collSM, ok := collM.(*sync.Map)
	if !ok {
		return nil
	}
	fields, ok := collSM.Load(coll)
	if !ok {
		return nil
	}
	fieldsS, ok := fields.([]string)
	if !ok {
		return nil
	}
	return fieldsS
}

func IsCrypto(db, coll, field string) bool {
	if CryptoMap == nil {
		return false
	}
	collM, ok := CryptoMap.Load(db)
	if !ok {
		return false
	}
	collSM, ok := collM.(*sync.Map)
	if !ok {
		return false
	}
	fields, ok := collSM.Load(coll)
	if !ok {
		return false
	}
	fieldsS, ok := fields.([]string)
	if !ok {
		return false
	}
	for _, f := range fieldsS {
		if f == field {
			return true
		}
	}
	return false
}

// Init 初始化连接池
func Init(opts ...*Opt) {
	for _, opt := range opts {
		if _, ok := pool.Load(opt.getAliasName()); ok {
			log.Panicf("[mongo]Mongo <%s> already registered\n", opt.getAliasName())
		}
		pool.Store(opt.getAliasName(), newCli(opt))
	}
}

func InitMigration(name string, client *mongo.Client) {
	pool.Store(name, client)
}

// 创建连接
func newCli(opt *Opt) *mongo.Client {
	cliOp := options.Client().ApplyURI(opt.getUri())
	cliOp.AppName = &app.Name
	if opt.SkipTLSVerify {
		cliOp.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}
	if opt.ReadPreference != nil {
		cliOp.SetReadPreference(opt.ReadPreference)
	}
	cliOp.SetMaxPoolSize(opt.MaxPoolSize)
	cliOp.SetMinPoolSize(opt.MinPoolSize)
	if opt.HeartbeatInterval != 0 {
		cliOp.SetHeartbeatInterval(opt.HeartbeatInterval)
	}
	cliOp.SetMaxConnecting(opt.getMaxConnecting())
	if opt.MaxConnIdleTime == 0 {
		opt.MaxConnIdleTime = 30 * time.Minute
	}
	cliOp.SetMaxConnIdleTime(opt.MaxConnIdleTime)
	cliOp.SetPoolMonitor(opt.PoolMonitor)
	cliOp.SetMonitor(opt.CommandMonitor)
	cliOp.SetSocketTimeout(opt.SocketTimeout)
	cliOp.SetRetryReads(true)
	cliOp.SetRetryWrites(true)
	if opt.Timeout <= 0 {
		opt.Timeout = defaultTimeout
	}
	cliOp.SetTimeout(opt.Timeout)
	cliOp.SetMonitor(Monitor)
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	c, err := mongo.Connect(ctx, cliOp)
	if err != nil {
		log.Panicf("[mongo]Connect to <%s> Failed %v\n", opt.Uri, err)
	}
	return c
}

// DB 获取一个DB连接对象
func DB(aliasName ...string) *Cli {
	name := "default"
	if len(aliasName) != 0 {
		name = aliasName[0]
	}
	v, ok := pool.Load(name)
	if !ok {
		log.Panicf("no %s cli in mongoDB pool\n", name)
	}
	db, ok := v.(*mongo.Client)
	if ok {
		return &Cli{
			AliasName: name,
			db:        db,
		}
	}
	return nil
}

// GenObjectID 生成ObjectId
func GenObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

// GenObjectIDFromTimestamp 使用指定时间戳生成objectId
func GenObjectIDFromTimestamp(ts time.Time) primitive.ObjectID {
	return primitive.NewObjectIDFromTimestamp(ts)
}

// WithDBOpt 为当前选到的连接对象，使用DatabaseOptions配置参数
func (c *Cli) WithDBOpt(dbOpt *options.DatabaseOptions) *Cli {
	c.dbOpt = dbOpt
	return c
}

// WithCollOpt 为当前选到的连接对象，使用CollectionOptions配置参数
func (c *Cli) WithCollOpt(collOpt *options.CollectionOptions) *Cli {
	c.collOpt = collOpt
	return c
}

// WithCtx 传入自定义context
func (c *Cli) WithCtx(ctx context.Context) *Cli {
	c.ctx = ctx
	return c
}

func (c *Cli) getDBOpt() *options.DatabaseOptions {
	return c.dbOpt
}

func (c *Cli) getCollOpt() *options.CollectionOptions {
	return c.collOpt
}

func (c *Cli) getCtx() context.Context {
	if c.ctx == nil {
		return context.TODO()
	}
	return c.ctx
}

// Database 获取名为db的Database对象
func (c *Cli) Database(db string) *mongo.Database {
	return c.db.Database(db, c.getDBOpt())
}

// Collection 获取Database中的名为coll的集合对象
func (c *Cli) Collection(db, coll string) *mongo.Collection {
	return c.db.Database(db, c.getDBOpt()).Collection(coll, c.getCollOpt())
}

// Indexes 执行 createIndexes 命令以在集合上创建多个索引并返回新索引的名称。
//
// 对于模型参数中的每个 IndexModel，可以通过 Options 字段指定索引名称。
//
// 如果没有给出名称，它将从 Keys 文档中生成。
//
// opts 参数可用于指定此操作的选项（请参阅 options.CreateIndexesOptions 文档）。
//
// 有关该命令的更多信息，请参阅 https://www.mongodb.com/docs/manual/reference/command/createIndexes/。
func (c *Cli) Indexes(db, coll string, idxes []IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {
	return c.db.Database(db, c.getDBOpt()).Collection(coll, c.getCollOpt()).Indexes().CreateMany(c.getCtx(), ToIndexModel(idxes), opts...)
}

// Find 执行查找命令并返回集合中匹配文档的Cursor
//
// filter 参数必须是包含查询运算符的文档，可用于选择结果中包含哪些文档。不能为nil。应该使用空文档（例如 bson.D{}）来表示包含所有文档。
//
// opts 参数可用于指定操作的选项（请参阅 options.FindOptions 文档）。
//
// 有关该命令的更多信息，请参阅 https://www.mongodb.com/docs/manual/reference/command/find/。
func (c *Cli) Find(db, coll string, filter any, opts ...*options.FindOptions) (*Cursor, error) {
	cursor, err := c.db.Database(db, c.getDBOpt()).Collection(coll, c.getCollOpt()).Find(c.getCtx(), filter, opts...)
	return &Cursor{cursor, db, coll}, err
}

// FindOne 执行查找命令并为集合中的一个文档返回SingleResult
//
// filter 参数必须是包含查询运算符的文档，可用于选择要返回的文档。
// 不能为nil。如果过滤器不匹配任何文档，将返回错误设置为ErrNoDocuments的SingleResult 。
// 如果过滤器匹配多个文档，将从匹配的集合中选择一个。
//
// opts 参数可用于指定此操作的选项（请参阅选项 options.FindOneOptions 文档）。
//
// 有关该命令的更多信息，请参阅 https://www.mongodb.com/docs/manual/reference/command/find/。
func (c *Cli) FindOne(db, coll string, filter any, opts ...*options.FindOneOptions) *SingleResult {
	return &SingleResult{c.db.Database(db, c.getDBOpt()).Collection(coll, c.getCollOpt()).FindOne(c.getCtx(), filter, opts...), db, coll}
}

// FindOneAndUpdate 执行 findAndModify 命令以更新集合中至多一个文档，并返回更新前出现的文档。
//
// filter 参数必须是包含查询运算符的文档，可用于选择要更新的文档。它不能为nil。如果过滤器不匹配任何文档，将返回一个错误设置为 ErrNoDocuments 的 SingleResult。
// 如果过滤器匹配多个文档，将从匹配的集合中选择一个。
//
// 更新参数必须是包含更新操作符的文档（https://www.mongodb.com/docs/manual/reference/operator/update/），可用于指定要对所选文档进行的修改。
// 它不能为nil或为空。
//
// opts 参数可用于指定操作的选项（请参阅 options.FindOneAndUpdateOptions 文档）。
//
// 有关该命令的更多信息，请参阅 https://www.mongodb.com/docs/manual/reference/command/findAndModify/。
func (c *Cli) FindOneAndUpdate(db, coll string, filter, update any, opts ...*options.FindOneAndUpdateOptions) *SingleResult {
	if err := globalOpts.getPreExec()(db, coll, update); err != nil {
		return nil
	}
	return &SingleResult{c.db.Database(db, c.getDBOpt()).Collection(coll, c.getCollOpt()).FindOneAndUpdate(c.getCtx(), filter, update, opts...), db, coll}
}

// FindOneAndReplace 执行 findAndModify 命令以替换集合中至多一个文档，并返回替换前出现的文档。
//
// filter 参数必须是包含查询运算符的文档，可用于选择要替换的文档。它不能为nil。
// 如果过滤器不匹配任何文档，将返回一个错误设置为ErrNoDocuments的SingleResult 。
// 如果过滤器匹配多个文档，将从匹配的集合中选择一个。
//
// 替换参数必须是将用于替换所选文档的文档。
// 不能为 nil 并且不能包含任何更新运算符 (https://www.mongodb.com/docs/manual/reference/operator/update/)。
// opts 参数可用于指定操作的选项（请参阅选项 options.FindOneAndReplaceOptions 文档）。
//
// 有关该命令的更多信息，请参阅 https://www.mongodb.com/docs/manual/reference/command/findAndModify/。
func (c *Cli) FindOneAndReplace(db, coll string, filter, replacement any, opts ...*options.FindOneAndReplaceOptions) *SingleResult {
	if err := globalOpts.getPreExec()(db, coll, replacement); err != nil {
		return nil
	}
	return &SingleResult{c.db.Database(db, c.getDBOpt()).Collection(coll, c.getCollOpt()).FindOneAndReplace(c.getCtx(), filter, replacement, opts...), db, coll}
}

// FindOneAndDelete 执行 findAndModify 命令以删除集合中至多一个文档。并返回删除前出现的文档。
//
// filter 参数必须是包含查询运算符的文档，可用于选择要删除的文档。不能为nil。
// 如果过滤器不匹配任何文档，将返回一个错误设置为ErrNoDocuments的SingleResult 。
// 如果过滤器匹配多个文档，将从匹配的集合中选择一个。
//
// opts 参数可用于指定操作选项（请参阅选项 options.FindOneAndDeleteOptions 文档）。
//
// 有关该命令的更多信息，请参阅 https://www.mongodb.com/docs/manual/reference/command/findAndModify/。
func (c *Cli) FindOneAndDelete(db, coll string, filter any, opts ...*options.FindOneAndDeleteOptions) *SingleResult {
	return &SingleResult{c.db.Database(db, c.getDBOpt()).Collection(coll, c.getCollOpt()).FindOneAndDelete(c.getCtx(), filter, opts...), db, coll}
}

// Exists 判断符合条件的文档是否存在，只返回是否存在
//
// filter 参数必须是包含查询运算符的文档，可用于选择要返回的文档。
// 不能为nil。如果过滤器不匹配任何文档，将返回false 。
// 如果过滤器匹配多个文档，将返回true。
//
// opts 参数可用于指定此操作的选项（请参阅选项 options.FindOneOptions 文档）。
// 有关该命令的更多信息，请参阅 https://www.mongodb.com/docs/manual/reference/command/find/.
func (c *Cli) Exists(db, coll string, filter any, opts ...*options.FindOneOptions) bool {
	opts = append(opts, &options.FindOneOptions{
		Projection: bson.M{"_id": 1},
	})
	return c.db.Database(db, c.getDBOpt()).Collection(coll, c.getCollOpt()).FindOne(c.getCtx(), filter, opts...).Err() == nil
}

// EExists 判断符合条件的文档是否存在，只返回是否存在，已Error为nil作为依据条件
//
// filter 参数必须是包含查询运算符的文档，可用于选择要返回的文档。
// 不能为nil。如果过滤器不匹配任何文档，将返回 ErrNoDocuments 。
// 如果过滤器遇到其他错误，将返回对应的 Error。
// 如果过滤器匹配单个或多个文档，将返回nil。
//
// opts 参数可用于指定此操作的选项（请参阅选项 options.FindOneOptions 文档）。
// 有关该命令的更多信息，请参阅 https://www.mongodb.com/docs/manual/reference/command/find/.
func (c *Cli) EExists(db, coll string, filter any, opts ...*options.FindOneOptions) error {
	opts = append(opts, &options.FindOneOptions{
		Projection: bson.M{"_id": 1},
	})
	return c.db.Database(db, c.getDBOpt()).Collection(coll, c.getCollOpt()).FindOne(c.getCtx(), filter, opts...).Err()
}

// Update 执行更新命令来更新集合中的文档。
//
// filter 参数必须是包含查询运算符的文档，可用于选择要更新的文档。不能为nil。如果过滤器不匹配任何文档，操作将成功并返回 MatchedCount 为 0 的UpdateResult 。
//
// 更新参数必须是一个包含更新操作符的文档（ https://www.mongodb.com/docs/manual/reference/operator/update/ ），可用于指定对所选文档进行的修改。 不能为nil或为空。
//
// opts 参数可用于指定操作的选项（请参阅 options.UpdateOptions 文档）。
//
// 有关该命令的更多信息，请参阅 https://www.mongodb.com/docs/manual/reference/command/update/。
func (c *Cli) Update(db, coll string, filter, update any, opts ...*options.UpdateOptions) (*UpdateResult, error) {
	if err := globalOpts.getPreExec()(db, coll, update); err != nil {
		return nil, err
	}
	result, err := c.db.Database(db, c.getDBOpt()).Collection(coll, c.getCollOpt()).UpdateMany(c.getCtx(), filter, update, opts...)
	return &UpdateResult{result}, err
}

// UpdateOne 执行更新命令以更新集合中至多一个文档。
//
// filter 参数必须是包含查询运算符的文档，可用于选择要更新的文档。它不能为nil。如果过滤器不匹配任何文档，操作将成功并返回 MatchedCount 为 0 的 UpdateResult。
// 如果过滤器匹配多个文档，将从匹配的集合中选择一个，MatchedCount 将等于 1。
// 更新参数必须是包含更新操作符的文档（https://www.mongodb.com/docs/manual/reference/operator/update/），可用于指定要修改的内容对所选文档进行。 它不能为nil或为空。
//
// opts 参数可用于指定操作选项（请参阅 options.UpdateOptions 文档）。
// 有关该命令的更多信息，请参阅 https://www.mongodb.com/docs/manual/reference/command/update/.。
func (c *Cli) UpdateOne(db, coll string, filter, update any, opts ...*options.UpdateOptions) (*UpdateResult, error) {
	if err := globalOpts.getPreExec()(db, coll, update); err != nil {
		return nil, err
	}
	result, err := c.db.Database(db, c.getDBOpt()).Collection(coll, c.getCollOpt()).UpdateOne(c.getCtx(), filter, update, opts...)
	return &UpdateResult{result}, err
}

// UpdateByID 执行更新命令以更新其 _id 值与集合中提供的 ID 匹配的文档。
// 相当于运行UpdateOne (ctx, bson.D{{"_id", id}}, update, opts...)。
//
// id 参数是要更新的文档的_id。不能为nil。如果 ID 不匹配任何文档，操作将成功并返回 MatchedCount 为 0 的UpdateResult 。
//
// update 参数必须是一个包含更新操作符的文档（ https://www.mongodb.com/docs/manual/reference/operator/update/ ），可用于指定要对所选文档进行的修改。它不能为nil或为空。
//
// opts 参数可用于指定操作的选项（请参阅 options.UpdateOptions 文档）。
//
// 有关该命令的更多信息，请参阅 https://www.mongodb.com/docs/manual/reference/command/update/。
func (c *Cli) UpdateByID(db, coll string, id, update any, opts ...*options.UpdateOptions) (*UpdateResult, error) {
	if err := globalOpts.getPreExec()(db, coll, update); err != nil {
		return nil, err
	}
	result, err := c.db.Database(db, c.getDBOpt()).Collection(coll, c.getCollOpt()).UpdateByID(c.getCtx(), id, update, opts...)
	return &UpdateResult{result}, err
}

// InsertMany 执行插入命令以将多个文档插入到集合中。
// 如果在操作期间发生写入错误（例如重复键错误），此方法将返回BulkWriteException错误。
//
// documents 参数必须是要插入的文档切片。切片不能为nil或为空。元素必须全部为非nil。对于转换为 BSON 时没有 _id 字段的任何文档，一个将自动添加到编组文档中。原始文档不会被修改。可以从返回的InsertManyResult的 InsertedIDs 字段中检索插入文档的 _id 值。
//
// opts 参数可用于指定操作的选项（请参阅 options.InsertManyOptions 文档。）
//
// 有关该命令的更多信息，请参阅 https://www.mongodb.com/docs/manual/reference/command/insert/
func (c *Cli) InsertMany(db, coll string, documents []any, opts ...*options.InsertManyOptions) (*InsertManyResult, error) {
	if err := globalOpts.getPreExec()(db, coll, documents...); err != nil {
		return nil, err
	}
	result, err := c.db.Database(db, c.getDBOpt()).Collection(coll, c.getCollOpt()).InsertMany(c.getCtx(), documents, opts...)
	return &InsertManyResult{result}, err
}

// InsertOne 执行插入命令以将单个文档插入到集合中。
//
// 文档参数必须是要插入的文档。不能为nil。
// 如果文档在转换为 BSON 时没有 _id 字段，则会自动将一个字段添加到编组文档中，原始文档不会被修改。
// _id 可以从返回的InsertOneResult的 InsertedID 字段中检索。
//
// opts 参数可用于指定操作的选项（请参阅 options.InsertOneOptions 文档。）
//
// 有关该命令的更多信息，请参阅 https://www.mongodb.com/docs/manual/reference/command/insert/
func (c *Cli) InsertOne(db, coll string, document any, opts ...*options.InsertOneOptions) (*InsertOneResult, error) {
	if err := globalOpts.getPreExec()(db, coll, document); err != nil {
		return nil, err
	}
	result, err := c.db.Database(db, c.getDBOpt()).Collection(coll, c.getCollOpt()).InsertOne(c.getCtx(), document, opts...)
	return &InsertOneResult{result}, err
}

// Delete 执行删除命令以从集合中删除文档。
//
// filter 参数必须是包含查询运算符的文档，可用于选择要删除的文档。不能为nil。应该使用一个空文档（例如 bson.D{}）来删除集合中的所有文档。
// 如果过滤器不匹配任何文档，操作将成功并返回 DeletedCount 为 0 的DeleteResult 。
//
// opts 参数可用于指定操作的选项（请参阅 options.DeleteOptions 文档）。
//
// 有关该命令的更多信息，请参阅 https://www.mongodb.com/docs/manual/reference/command/delete/ 。
func (c *Cli) Delete(db, coll string, filter any, opts ...*options.DeleteOptions) (*DeleteResult, error) {
	result, err := c.db.Database(db, c.getDBOpt()).Collection(coll, c.getCollOpt()).DeleteMany(c.getCtx(), filter, opts...)
	return &DeleteResult{result}, err
}

// DeleteOne 执行删除命令以从集合中删除最多一个文档。
//
// filter 参数必须是包含查询运算符的文档，可用于选择要删除的文档。不能为nil。
// 如果过滤器不匹配任何文档，操作将成功并返回 DeletedCount 为 0 的DeleteResult。
// 如果过滤器匹配多个文档，将从匹配的集合中选择一个。
//
// opts 参数可用于指定操作的选项（请参阅 options.DeleteOptions 文档）。
//
// 有关该命令的更多信息，请参阅 https://www.mongodb.com/docs/manual/reference/command/delete/ 。
func (c *Cli) DeleteOne(db, coll string, filter any, opts ...*options.DeleteOptions) (*DeleteResult, error) {
	result, err := c.db.Database(db, c.getDBOpt()).Collection(coll, c.getCollOpt()).DeleteOne(c.getCtx(), filter, opts...)
	return &DeleteResult{result}, err
}

// Count 返回集合中的文档数。
//
// 有关集合中文档的快速计数，请参阅 EstimatedDocumentCount 方法。
//
// filter 参数必须是文档，可用于选择哪些文档有助于计数。不能为nil。应该使用空文档（例如 bson.D{}）来计算集合中的所有文档。这将导致完整的集合扫描。
//
// opts 参数可用于指定操作的选项（请参阅 options.CountOptions 文档）。
func (c *Cli) Count(db, coll string, filter any, opts ...*options.CountOptions) (int64, error) {
	return c.db.Database(db, c.getDBOpt()).Collection(coll, c.getCollOpt()).CountDocuments(c.getCtx(), filter, opts...)
}

// EstimatedCount 执行计数命令并使用集合元数据返回集合中文档数量的估计值。
//
// opts 参数可用于指定操作的选项（请参阅 options.EstimatedDocumentCountOptions 文档）。
//
// 有关该命令的更多信息，请参阅 https://www.mongodb.com/docs/manual/reference/command/count/ 。
func (c *Cli) EstimatedCount(db, coll string, opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	return c.db.Database(db, c.getDBOpt()).Collection(coll, c.getCollOpt()).EstimatedDocumentCount(c.getCtx(), opts...)
}

// Aggregate 对集合执行聚合命令并返回结果文档上的游标。
//
// pipeline 参数必须是一个文档数组，每个文档代表一个聚合阶段。管道不能为零，但可以为空。阶段文件必须全部为非零。对于 bson.D 文件的管道， mongo 。可以使用管道类型。有关聚合中有效阶段的列表，请参阅https://www.mongodb.com/docs/manual/reference/operator/aggregation-pipeline/#db-collection-aggregate-stages 。
//
// opts 参数可用于指定操作选项（请参阅 options.AggregateOptions 文档。）
//
// 有关该命令的更多信息，请参阅 https://www.mongodb.com/docs/manual/reference/command/aggregate/ 。
func (c *Cli) Aggregate(db, coll string, pip any, opts ...*options.AggregateOptions) (*Cursor, error) {
	cursor, err := c.db.Database(db, c.getDBOpt()).Collection(coll, c.getCollOpt()).Aggregate(c.getCtx(), pip, opts...)
	return &Cursor{cursor, db, coll}, err
}

// Distinct 执行不同的命令以查找集合中指定字段的唯一值。
//
// fieldName 参数指定应为其返回不同值的字段名称。
//
// filter 参数必须是包含查询运算符的文档，可用于选择考虑哪些文档。它不能为零。应该使用空文档（例如 bson.D{}）来选择所有文档。
//
// opts 参数可用于指定操作的选项（请参阅 options.DistinctOptions 文档）。
//
// 有关该命令的更多信息，请参阅 https://www.mongodb.com/docs/manual/reference/command/distinct/ 。
func (c *Cli) Distinct(db, coll string, fieldName string, filter any, opts ...*options.DistinctOptions) ([]any, error) {
	return c.db.Database(db, c.getDBOpt()).Collection(coll, c.getCollOpt()).Distinct(c.getCtx(), fieldName, filter, opts...)
}

// BulkWrite 执行批量写入操作 https://www.mongodb.com/docs/manual/core/bulk-write-operations/ 。
// models 参数必须是要在此批量写入中执行的 slice 。它不能为零或为 nil。所有模型都必须非 nil。
// mongo.WriteModel 文档提供有效模型类型的列表以及如何使用它们的示例。已在 func 中提供
// opts 参数可用于指定操作的选项（请参阅 options.BulkWriteOptions 文档。）
func (c *Cli) BulkWrite(db, coll string, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	return c.db.Database(db, c.getDBOpt()).Collection(coll, c.getCollOpt()).BulkWrite(c.getCtx(), models, opts...)
}

// Watch 为相应集合上的所有更改返回一个更改流。
//
// 有关更改流的更多信息，请参阅 https://www.mongodb.com/docs/manual/changeStreams/ 。
//
// 集合必须配置为read concern majority 或no read concern才能成功创建更改流。
//
// 管道参数必须是文档数组，每个文档代表一个管道阶段。管道不能为零，但可以为空。阶段文件必须全部为非零。
// 有关可与更改流一起使用的管道阶段列表，请参阅 https://www.mongodb.com/docs/manual/changeStreams/ 。对于 bson.D 文件的pipeline，可以使用 mongo.Pipeline{} 类型。
//
// opts 参数可用于指定更改流创建的选项（请参阅选项 options.ChangeStreamOptions 文档）。
func (c *Cli) Watch(db, coll string, pip any, opts ...*options.ChangeStreamOptions) (*ChangeStream, error) {
	result, err := c.db.Database(db, c.getDBOpt()).Collection(coll, c.getCollOpt()).Watch(c.getCtx(), pip, opts...)
	return &ChangeStream{result, db, coll}, err
}

// RunCommand 对数据库执行给定的命令。此函数不服从数据库的读取首选项。要指定读取首选项，必须使用 RunCmdOptions.ReadPreference 选项。
//
// runCommand 参数必须是要执行的命令的文档。不能为nil。这必须是保序类型，例如 bson.D。
//
// bson.M 等map类型无效。
//
// opts 参数可用于指定此操作的选项（请参阅 options.RunCmdOptions 文档）。
// 如果命令文档包含以下任何内容，则 RunCommand 的行为未定义：
//   - 会话 ID 或任何特定于事务的字段
//   - 当已在客户端上声明 API 版本时的 API 版本控制选项
//   - 当在客户端上设置超时时的 maxTimeMS
func (c *Cli) RunCommand(db string, runCommand any, opts ...*options.RunCmdOptions) *SingleResult {
	return &SingleResult{c.db.Database(db, c.getDBOpt()).RunCommand(c.getCtx(), runCommand, opts...), "", ""}
}

// SessionContext 事务结构
type SessionContext struct {
	cli      *Cli
	ctx      mongo.SessionContext
	rollback bool
}

// Begin 在此会话上启动一个新事务，配置了给定的选项。
// 如果此会话中已有正在进行的事务，则此方法将返回错误。
func (ctx *SessionContext) Begin(opt *options.TransactionOptions) error {
	return ctx.ctx.StartTransaction(opt)
}

// Rollback 设置回滚标记
func (ctx *SessionContext) Rollback() {
	ctx.rollback = true
}

// Context 获取当前Session的Context
func (ctx *SessionContext) Context() context.Context {
	return ctx.ctx
}

// End 结束当前会话上的事务
// 按照是否设置过 Rollback() 分一下两种情景：
// - CommitTransaction 提交此会话的活动事务。如果此会话没有活动事务或事务已中止，则此方法将返回错误。（此处忽略错误）
// - AbortTransaction 中止此会话的活动事务。如果此会话没有活动事务或事务已提交或中止，则此方法将返回错误。（此处忽略错误）
func (ctx *SessionContext) End() {
	if ctx.rollback {
		_ = ctx.ctx.AbortTransaction(ctx.ctx)
	}
	_ = ctx.ctx.CommitTransaction(ctx.ctx)
}

// UseSession 使用 SessionOptions 创建一个新的Session，并用它创建一个新的SessionContext，用于调用fn回调。
//
// SessionContext 参数必须用作应在会话下执行的 fn 回调中的任何操作的上下文参数。
// 回调返回后，创建的 Session 结束，这意味着即使 fn 返回错误，任何由 fn 启动的正在进行的事务也将被中止。
//
// UseSession 可以安全地从多个 goroutines 同时调用。但是，传递给 UseSession 回调函数的 SessionContext 对于多个 goroutines 的并发使用是不安全的。
// 如果 ctx 参数已经包含一个 Session，则该 Session 将被新创建的 Session 替换。
//
// fn 回调返回的任何错误都将在不做任何修改的情况下返回。但是，传递给 UseSessionWithOptions 回调函数的 SessionContext 对于多个 goroutines 的并发使用是不安全的。
func (c *Cli) UseSession(ctx context.Context, opt *options.SessionOptions, fn func(sCtx *SessionContext) error) error {
	return c.db.UseSessionWithOptions(ctx, opt, func(s mongo.SessionContext) error {
		return fn(&SessionContext{
			cli: c,
			ctx: s,
		})
	})
}
