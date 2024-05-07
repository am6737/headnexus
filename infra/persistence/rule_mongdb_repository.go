package persistence

import (
	"context"
	"errors"
	"fmt"
	ctime "github.com/am6737/headnexus/common/time"
	"github.com/am6737/headnexus/domain/rule/entity"
	"github.com/am6737/headnexus/domain/rule/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ruleCollectionName = "rules"
)

type RuleModel struct {
	Type        uint8              `bson:"type"`
	ID          primitive.ObjectID `bson:"_id"`
	CreatedAt   int64              `bson:"created_at"`
	UpdatedAt   int64              `bson:"updated_at"`
	DeletedAt   int64              `bson:"deleted_at"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	HostID      string             `bson:"host_id"`
	Port        string             `bson:"port"`
	Proto       string             `bson:"proto"`
	Action      uint8              `bson:"action"`
	Host        []string           `bson:"host,omitempty"`
}

func (m *RuleModel) FromEntity(e *entity.Rule) error {
	var id primitive.ObjectID
	if e.ID != "" {
		oid, err := primitive.ObjectIDFromHex(e.ID)
		if err != nil {
			return err
		}
		id = oid
	}
	m.ID = id
	m.Name = e.Name
	m.Description = e.Description
	m.HostID = e.HostID
	m.Port = e.Port
	m.Proto = string(e.Proto)
	m.Action = uint8(e.Action)
	m.Host = e.Host
	return nil
}

func (m *RuleModel) ToEntity() (*entity.Rule, error) {
	return &entity.Rule{
		ID:          m.ID.Hex(),
		Type:        entity.RuleType(m.Type),
		CreatedAt:   m.CreatedAt,
		Name:        m.Name,
		Description: m.Description,
		HostID:      m.HostID,
		Port:        m.Port,
		Proto:       entity.RuleProto(m.Proto),
		Action:      entity.RuleAction(m.Action),
		Host:        m.Host,
	}, nil
}

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

func (m *RuleMongodbRepository) Create(ctx context.Context, rule *entity.Rule) (*entity.Rule, error) {
	model := &RuleModel{}
	if err := model.FromEntity(rule); err != nil {
		return nil, err // 返回错误，以及 nil 指针作为结果
	}

	currentTime := ctime.CurrentTimestampMillis()
	model.CreatedAt = currentTime
	model.UpdatedAt = currentTime
	model.ID = primitive.NewObjectID()

	_, err := m.collection.InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}

	e, err := model.ToEntity()
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (m *RuleMongodbRepository) Get(ctx context.Context, id string) (*entity.Rule, error) {
	filter := bson.M{"_id": id}
	var model RuleModel
	err := m.collection.FindOne(ctx, filter).Decode(&model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrorHostNotFound
		}
		return nil, err
	}
	rule, err := model.ToEntity()
	if err != nil {
		return nil, err
	}
	return rule, nil
}

func (m *RuleMongodbRepository) Update(ctx context.Context, rule *entity.Rule) error {
	// 将 api.Host 转换为 HostDBModel
	model := &RuleModel{}
	if err := model.FromEntity(rule); err != nil {
		return err
	}

	filter := bson.M{"_id": rule.ID}
	update := bson.M{"$set": model}

	//opts := options.Update().SetUpsert(true)
	_, err := m.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (m *RuleMongodbRepository) Delete(ctx context.Context, id string) error {
	// 将字符串转换为 ObjectId
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// 使用 ObjectId 构建过滤条件
	filter := bson.M{"_id": objectID}
	_, err = m.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (m *RuleMongodbRepository) Find(ctx context.Context, options *entity.FindOptions) ([]*entity.Rule, error) {
	filter := bson.M{}

	if options.HostID != "" {
		filter["host_id"] = options.HostID
	}
	if options.Name != "" {
		filter["name"] = options.Name
	}

	findOptions := optionsToMongoFindOptions(options)

	fmt.Println("filter => ", filter)

	var rules []*entity.Rule

	cursor, err := m.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var rule RuleModel
		if err := cursor.Decode(&rule); err != nil {
			return nil, err
		}
		e, err := rule.ToEntity()
		if err != nil {
			return nil, err
		}
		rules = append(rules, e)
	}

	return rules, nil
}

func optionsToMongoFindOptions(efo *entity.FindOptions) *options.FindOptions {
	findOptions := &options.FindOptions{}

	if efo.PageSize != 0 {
		findOptions.SetLimit(int64(efo.PageSize))
	}

	if efo.PageNum != 0 {
		findOptions.SetSkip(int64((efo.PageNum - 1) * efo.PageSize))
	}

	return findOptions
}
