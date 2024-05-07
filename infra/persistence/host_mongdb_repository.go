package persistence

import (
	"context"
	"errors"
	"fmt"
	ctime "github.com/am6737/headnexus/common/time"
	"github.com/am6737/headnexus/config"
	"github.com/am6737/headnexus/domain/host/entity"
	"github.com/am6737/headnexus/domain/host/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrorHostNotFound = errors.New("host not found")

var _ repository.HostRepository = &HostMongodbRepository{}

const (
	hostCollectionName = "hosts"
)

type HostDBModel struct {
	ID              string                 `bson:"_id"`
	Name            string                 `bson:"name"`
	NetworkID       string                 `bson:"network_id"`
	IPAddress       string                 `bson:"ip_address"`
	Role            string                 `bson:"role"`
	EnrollCode      string                 `bson:"enroll_code"`
	Owner           string                 `bson:"owner"`
	Port            int                    `bson:"port"`
	IsLighthouse    bool                   `bson:"is_lighthouse"`
	StaticAddresses []string               `bson:"static_addresses"`
	Tags            map[string]interface{} `bson:"tags"`
	CreatedAt       int64                  `bson:"created_at"`
	UpdatedAt       int64                  `bson:"updated_at"`
	DeletedAt       int64                  `bson:"deleted_at"`
	LastSeenAt      int64                  `bson:"last_seen_at"`
	EnrollAt        int64                  `bson:"enroll_at"`
	LifetimeSeconds int64                  `bson:"lifetime_seconds"`
	Status          int8                   `bson:"status"`
	Config          config.Config          `bson:"Config"`
}

func (m *HostDBModel) collection() string {
	return "hosts"
}

func (m *HostDBModel) From(host *entity.Host) error {
	m.ID = host.ID
	m.Name = host.Name
	m.NetworkID = host.NetworkID
	m.IPAddress = host.IPAddress
	m.Role = host.Role
	m.Port = host.Port
	m.IsLighthouse = host.IsLighthouse
	m.StaticAddresses = host.StaticAddresses
	m.Tags = host.Tags
	m.Config = host.Config
	return nil
}

func (m *HostDBModel) To() (*entity.Host, error) {
	host := &entity.Host{
		ID:              m.ID,
		Name:            m.Name,
		NetworkID:       m.NetworkID,
		IPAddress:       m.IPAddress,
		Role:            m.Role,
		Port:            m.Port,
		IsLighthouse:    m.IsLighthouse,
		StaticAddresses: m.StaticAddresses,
		LastSeenAt:      m.LastSeenAt,
		Tags:            m.Tags,
		Config:          m.Config,
		Status:          entity.HostStatus(m.Status),
		//CreatedAt:       m.CreatedAt,
	}
	return host, nil
}

func NewHostMongodbRepository(client *mongo.Client, dbName string) *HostMongodbRepository {
	db := client.Database(dbName)
	collection := db.Collection(hostCollectionName)
	m := &HostMongodbRepository{
		client:     client,
		db:         db,
		collection: collection,
	}
	return m
}

type HostMongodbRepository struct {
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
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
	var model HostDBModel
	err := h.collection.FindOne(ctx, filter).Decode(&model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrorHostNotFound
		}
		return nil, err
	}
	return &entity.EnrollHost{
		HostID:          model.ID,
		Code:            model.EnrollCode,
		LifetimeSeconds: model.LifetimeSeconds,
		EnrollAt:        model.EnrollAt,
	}, nil
}

func (h *HostMongodbRepository) EnrollHost(ctx context.Context, enrollHost *entity.EnrollHost) error {
	filter := bson.M{"_id": enrollHost.HostID}

	update := bson.M{
		"$set": bson.M{
			"enroll_code":      enrollHost.Code,
			"lifetime_seconds": enrollHost.LifetimeSeconds,
			"enroll_at":        enrollHost.EnrollAt,
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := h.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func ToFindOptions(listOptions *entity.FindOptions) *options.FindOptions {
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

func (h *HostMongodbRepository) Find(ctx context.Context, options *entity.FindOptions) ([]*entity.Host, error) {
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
		var hostDBModel HostDBModel
		if err := cursor.Decode(&hostDBModel); err != nil {
			return nil, err
		}

		host, err := hostDBModel.To()
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
	model := &HostDBModel{}
	if err := model.From(host); err != nil {
		return nil, err
	}

	fmt.Println("model cfg => ", model.Config)

	currentTime := ctime.CurrentTimestampMillis()
	model.CreatedAt = currentTime
	model.UpdatedAt = currentTime
	_, err := h.collection.InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}
	return host, nil
}

func (h *HostMongodbRepository) Get(ctx context.Context, id string) (*entity.Host, error) {
	filter := bson.M{"_id": id}
	var model HostDBModel
	err := h.collection.FindOne(ctx, filter).Decode(&model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrorHostNotFound
		}
		return nil, err
	}
	host, err := model.To()
	if err != nil {
		return nil, err
	}
	return host, nil
}

func (h *HostMongodbRepository) Update(ctx context.Context, host *entity.Host) (*entity.Host, error) {
	// 将 api.Host 转换为 HostDBModel
	model := &HostDBModel{}
	if err := model.From(host); err != nil {
		return nil, err
	}

	// 构建更新操作所需的 filter 和 update
	filter := bson.M{"_id": host.ID}
	update := bson.M{"$set": model}

	// 执行更新操作
	opts := options.Update().SetUpsert(true)
	_, err := h.collection.UpdateOne(ctx, filter, update, opts)
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
