package mongo

import "go.mongodb.org/mongo-driver/mongo"

// IsDuplicateKeyError 是否为主键重复Error
func IsDuplicateKeyError(err error) bool {
	return mongo.IsDuplicateKeyError(err)
}

// IsTimeout 是否为超时Error
func IsTimeout(err error) bool {
	return mongo.IsTimeout(err)
}

// IsNetworkError 是否为网络Error
func IsNetworkError(err error) bool {
	return mongo.IsNetworkError(err)
}

// NewInsertOneModel creates a new InsertOneModel.
func NewInsertOneModel() *mongo.InsertOneModel {
	return mongo.NewInsertOneModel()
}

// NewDeleteOneModel creates a new DeleteOneModel.
func NewDeleteOneModel() *mongo.DeleteOneModel {
	return mongo.NewDeleteOneModel()
}

// NewDeleteManyModel creates a new DeleteManyModel.
func NewDeleteManyModel() *mongo.DeleteManyModel {
	return mongo.NewDeleteManyModel()
}

// NewReplaceOneModel creates a new ReplaceOneModel.
func NewReplaceOneModel() *mongo.ReplaceOneModel {
	return mongo.NewReplaceOneModel()
}

// NewUpdateOneModel creates a new UpdateOneModel.
func NewUpdateOneModel() *mongo.UpdateOneModel {
	return mongo.NewUpdateOneModel()
}

// NewUpdateManyModel creates a new UpdateManyModel.
func NewUpdateManyModel() *mongo.UpdateManyModel {
	return mongo.NewUpdateManyModel()
}

// NewWriteModels 创建[]WriteModels
func NewWriteModels() []mongo.WriteModel {
	return []mongo.WriteModel{}
}
