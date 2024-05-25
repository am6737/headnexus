package persistence

import (
	"context"
	"errors"
	ctime "github.com/am6737/headnexus/common/time"
	"github.com/am6737/headnexus/domain/host/entity"
	"github.com/am6737/headnexus/domain/host/repository"
	"github.com/am6737/headnexus/infra/persistence/converter"
	"github.com/am6737/headnexus/infra/persistence/po"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrorHostNotFound = errors.New("host not found")

var _ repository.HostRepository = &HostMongodbRepository{}

var _ repository.HostRuleRepository = &HostMongodbRepository{}

const (
	hostCollectionName     = "hosts"
	hostRuleCollectionName = "host_rules"
)

func NewHostMongodbRepository(client *mongo.Client, dbName string) *HostMongodbRepository {
	db := client.Database(dbName)
	collection := db.Collection(hostCollectionName)
	hostRuleCollection := db.Collection(hostRuleCollectionName)
	m := &HostMongodbRepository{
		client:             client,
		db:                 db,
		collection:         collection,
		hostRuleCollection: hostRuleCollection,
	}
	return m
}

type HostMongodbRepository struct {
	client             *mongo.Client
	db                 *mongo.Database
	collection         *mongo.Collection
	hostRuleCollection *mongo.Collection
}

func (h *HostMongodbRepository) AddHostRule(ctx context.Context, hostID string, ruleIDs ...string) error {
	var hostRules []interface{}
	for _, ruleID := range ruleIDs {
		hostRule := po.HostRuleRelation{
			ID:     primitive.NewObjectID().Hex(),
			HostID: hostID,
			RuleID: ruleID,
		}

		currentTime := ctime.CurrentTimestampMillis()
		hostRule.CreatedAt = currentTime
		hostRule.UpdatedAt = currentTime

		hostRules = append(hostRules, hostRule)
	}

	// 批量插入关联文档
	_, err := h.hostRuleCollection.InsertMany(ctx, hostRules)
	if err != nil {
		return err
	}

	return nil
}

func (h *HostMongodbRepository) ListHostRule(ctx context.Context, opts *entity.ListHostRuleOptions) ([]*entity.HostRuleRelation, error) {
	if opts == nil || opts.HostID == "" {
		return nil, errors.New("hostID is required")
	}

	filter := bson.M{"host_id": opts.HostID}

	// 添加其他过滤条件
	if opts.Type != nil {
		filter["type"] = opts.Type
	}
	if opts.Proto != nil {
		filter["proto"] = opts.Proto
	}
	if opts.Action != nil {
		filter["action"] = opts.Action
	}

	cursor, err := h.hostRuleCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rules []*entity.HostRuleRelation
	for cursor.Next(ctx) {
		model := &po.HostRuleRelation{}
		if err := cursor.Decode(model); err != nil {
			return nil, err
		}

		rule := converter.HostRulePOToEntity(model)
		rules = append(rules, rule)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return rules, nil
}

func (h *HostMongodbRepository) HostOnline(ctx context.Context, hostOnline *entity.HostOnline) (*entity.HostOnline, error) {
	filter := bson.M{"_id": hostOnline.ID}

	at := ctime.CurrentTimestampMillis()

	update := bson.M{
		"$set": bson.M{
			"status":       entity.Online,
			"last_seen_at": at,
		},
	}

	_, err := h.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &entity.HostOnline{
		ID:         hostOnline.ID,
		Status:     entity.Online,
		LastSeenAt: at,
	}, nil
}

func (h *HostMongodbRepository) HostOffline(ctx context.Context, hostOffline *entity.HostOffline) (*entity.HostOffline, error) {
	filter := bson.M{"_id": hostOffline.ID}

	at := ctime.CurrentTimestampMillis()

	update := bson.M{
		"$set": bson.M{
			"status":       entity.Offline,
			"last_seen_at": at,
		},
	}

	_, err := h.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &entity.HostOffline{
		ID:         hostOffline.ID,
		Status:     entity.Offline,
		LastSeenAt: at,
	}, nil
}

func (h *HostMongodbRepository) GetEnrollHost(ctx context.Context, getEnrollHost *entity.GetEnrollHost) (*entity.EnrollHost, error) {
	filter := bson.M{"_id": getEnrollHost.HostID}
	model := &po.Host{}
	if err := h.collection.FindOne(ctx, filter).Decode(model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrorHostNotFound
		}
		return nil, err
	}
	return &entity.EnrollHost{
		HostID:              model.ID,
		Code:                model.EnrollCode,
		EnrollAt:            model.EnrollAt,
		CreatedAt:           model.CreatedAt,
		EnrollCodeExpiredAt: model.EnrollCodeExpiredAt,
	}, nil
}

func (h *HostMongodbRepository) EnrollHost(ctx context.Context, enrollHost *entity.EnrollHost) error {
	filter := bson.M{"_id": enrollHost.HostID}

	update := bson.M{
		"$set": bson.M{
			"enroll_code":            enrollHost.Code,
			"enroll_code_expired_at": enrollHost.EnrollCodeExpiredAt,
			"enroll_at":              enrollHost.EnrollAt,
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := h.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func ToFindOptions(listOptions *entity.HostFindOptions) *options.FindOptions {
	findOptions := options.FindOptions{}

	// Translate FindOptions.Sort into FindOptions.Sort
	if len(listOptions.Sort) > 0 {
		sortMap := make(map[string]interface{})
		for key, value := range listOptions.Sort {
			sortMap[key] = value
		}
		findOptions.SetSort(sortMap)
	}

	// Set other options based on FindOptions fields
	if listOptions.Limit > 0 {
		findOptions.SetLimit(int64(listOptions.Limit))
	}
	if listOptions.Offset > 0 {
		findOptions.SetSkip(int64(listOptions.Offset))
	}

	// Add filter based on FindOptions.Filters
	filter := bson.M{}
	for key, value := range listOptions.Filters {
		filter[key] = value
	}
	findOptions.SetProjection(filter)

	return &findOptions
}

type FindOptions struct {
	Filters map[string]interface{}
	Sort    map[string]int // 1 for ascending, -1 for descending
	Limit   int
	Offset  int
}

func (h *HostMongodbRepository) Find(ctx context.Context, options *entity.HostFindOptions) ([]*entity.Host, error) {
	filter := bson.M{}
	if options.NetworkID != "" {
		filter["network_id"] = options.NetworkID
	}
	if options.IPAddress != "" {
		filter["ip_address"] = options.IPAddress
	}
	if options.Role != "" {
		filter["role"] = options.Role
	}
	if options.Name != "" {
		filter["name"] = options.Name
	}
	filter["is_lighthouse"] = options.IsLighthouse

	cursor, err := h.collection.Find(ctx, filter, ToFindOptions(options))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var hosts []*entity.Host
	for cursor.Next(ctx) {
		model := &po.Host{}
		if err := cursor.Decode(model); err != nil {
			return nil, err
		}

		host, err := converter.HostPOToEntity(model)
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, host)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return hosts, nil
}

func (h *HostMongodbRepository) Create(ctx context.Context, host *entity.Host) (*entity.Host, error) {
	model, err := converter.HostEntityToPO(host)
	if err != nil {
		return nil, err
	}

	currentTime := ctime.CurrentTimestampMillis()
	model.CreatedAt = currentTime
	model.UpdatedAt = currentTime
	_, err = h.collection.InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}
	return host, nil
}

func (h *HostMongodbRepository) Get(ctx context.Context, id string) (*entity.Host, error) {
	filter := bson.M{"_id": id}

	model := &po.Host{}

	if err := h.collection.FindOne(ctx, filter).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrorHostNotFound
		}
		return nil, err
	}

	host, err := converter.HostPOToEntity(model)
	if err != nil {
		return nil, err
	}

	return host, nil
}

func (h *HostMongodbRepository) Update(ctx context.Context, host *entity.Host) (*entity.Host, error) {
	model, err := converter.HostEntityToPO(host)
	if err != nil {
		return nil, err
	}

	// 构建更新操作所需的 filter 和 update
	filter := bson.M{"_id": host.ID}
	update := bson.M{"$set": model}

	// 执行更新操作
	opts := options.Update().SetUpsert(true)
	_, err = h.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}

	// 返回更新后的 api.Host
	return host, nil
}

func (h *HostMongodbRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	_, err := h.collection.DeleteOne(ctx, filter)
	return err
}
