package persistence

import (
	"context"
	"errors"
	ctime "github.com/am6737/headnexus/common/time"
	"github.com/am6737/headnexus/domain/host/entity"
	"github.com/am6737/headnexus/domain/host/repository"
	"github.com/am6737/headnexus/infra/persistence/converter"
	"github.com/am6737/headnexus/infra/persistence/po"
	"github.com/am6737/headnexus/pkg/code"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ruleCollectionName = "rules"
)

var _ repository.RuleRepository = &RuleMongodbRepository{}

func NewRuleMongodbRepository(client *mongo.Client, dbName string) *RuleMongodbRepository {
	db := client.Database(dbName)
	collection := db.Collection(ruleCollectionName)
	m := &RuleMongodbRepository{
		client:     client,
		db:         db,
		collection: collection,
	}
	return m
}

type RuleMongodbRepository struct {
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
}

func (m *RuleMongodbRepository) Create(ctx context.Context, userID string, rule *entity.Rule) (*entity.Rule, error) {
	model := converter.RuleEntityToPO(rule)

	currentTime := ctime.CurrentTimestampMillis()
	model.CreatedAt = currentTime
	model.UpdatedAt = currentTime
	model.ID = primitive.NewObjectID().Hex()
	model.UserID = userID // 添加 user_id

	_, err := m.collection.InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}

	e := converter.RulePOToEntity(model)
	return e, nil
}

func (m *RuleMongodbRepository) Gets(ctx context.Context, userID string, ids []string) ([]*entity.Rule, error) {
	if userID == "" || len(ids) == 0 {
		return nil, code.InvalidParameter
	}

	filter := bson.M{
		"user_id": userID,
		"_id":     bson.M{"$in": ids},
	}

	// 执行查询
	cursor, err := m.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 将结果解码为 PO 对象的切片
	var models []*po.Rule
	if err = cursor.All(ctx, &models); err != nil {
		return nil, err
	}

	// 将 PO 对象转换为 Entity 对象
	var rules []*entity.Rule
	for _, model := range models {
		rule := converter.RulePOToEntity(model)
		rules = append(rules, rule)
	}

	return rules, nil
}

func (m *RuleMongodbRepository) Get(ctx context.Context, userID, id string) (*entity.Rule, error) {
	filter := bson.M{"_id": id, "user_id": userID}
	model := &po.Rule{}
	err := m.collection.FindOne(ctx, filter).Decode(model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, code.NotFound
		}
		return nil, err
	}

	e := converter.RulePOToEntity(model)
	return e, nil
}

func (m *RuleMongodbRepository) Update(ctx context.Context, userID string, rule *entity.Rule) error {
	model := converter.RuleEntityToPO(rule)

	filter := bson.M{"_id": rule.ID, "user_id": userID}
	update := bson.M{"$set": model}

	_, err := m.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (m *RuleMongodbRepository) Delete(ctx context.Context, id string) error {

	// 使用 ObjectId 构建过滤条件
	filter := bson.M{"_id": id}
	_, err := m.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (m *RuleMongodbRepository) Find(ctx context.Context, userID string, options *entity.RuleFindOptions) ([]*entity.Rule, error) {
	if options == nil || userID == "" { // 检查 user_id 是否为空
		return nil, errors.New("user_id is empty")
	}

	filter := bson.M{
		"user_id": userID,
	}

	if options.HostID != "" {
		filter["host_id"] = options.HostID
	}
	if options.Name != "" {
		filter["name"] = options.Name
	}

	findOptions := optionsToMongoFindOptions(options)

	var rules []*entity.Rule

	cursor, err := m.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		model := &po.Rule{}
		if err := cursor.Decode(model); err != nil {
			return nil, err
		}

		e := converter.RulePOToEntity(model)
		rules = append(rules, e)
	}

	return rules, nil
}

func optionsToMongoFindOptions(efo *entity.RuleFindOptions) *options.FindOptions {
	findOptions := &options.FindOptions{}

	if efo.PageSize != 0 {
		findOptions.SetLimit(int64(efo.PageSize))
	}

	if efo.PageNum != 0 {
		findOptions.SetSkip(int64((efo.PageNum - 1) * efo.PageSize))
	}

	return findOptions
}
