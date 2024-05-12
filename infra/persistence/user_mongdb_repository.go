package persistence

import (
	"context"
	"errors"
	"fmt"
	ctime "github.com/am6737/headnexus/common/time"
	"github.com/am6737/headnexus/domain/user/entity"
	"github.com/am6737/headnexus/domain/user/repository"
	"github.com/am6737/headnexus/infra/persistence/converter"
	"github.com/am6737/headnexus/infra/persistence/po"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	userCollectionName = "users"
)

var _ repository.UserRepository = &UserMongodbRepository{}

type UserMongodbRepository struct {
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
}

func NewUserMongodbRepository(client *mongo.Client, dbName string) *UserMongodbRepository {
	db := client.Database(dbName)
	collection := db.Collection(userCollectionName)
	m := &UserMongodbRepository{
		client:     client,
		db:         db,
		collection: collection,
	}
	return m
}

func (m *UserMongodbRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	model, err := converter.UserEntityToPO(user)
	if err != nil {
		return nil, err
	}

	currentTime := ctime.CurrentTimestampMillis()
	model.CreatedAt = currentTime
	model.UpdatedAt = currentTime
	model.ID = primitive.NewObjectID().Hex()

	_, err = m.collection.InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}

	e, err := converter.UserPOToEntity(model)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (m *UserMongodbRepository) Get(ctx context.Context, id string) (*entity.User, error) {
	filter := bson.M{"_id": id}
	model := &po.User{}
	err := m.collection.FindOne(ctx, filter).Decode(&model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrorHostNotFound
		}
		return nil, err
	}
	e, err := converter.UserPOToEntity(model)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (m *UserMongodbRepository) Update(ctx context.Context, user *entity.User) error {
	model, err := converter.UserEntityToPO(user)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": model}

	//opts := options.Update().SetUpsert(true)
	_, err = m.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserMongodbRepository) Delete(ctx context.Context, id string) error {
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

func (m *UserMongodbRepository) Find(ctx context.Context, options *entity.FindOptions) ([]*entity.User, error) {
	filter := bson.M{}

	if options.Email != "" {
		filter["email"] = options.Email
	}
	if options.Token != "" {
		filter["token"] = options.Token
	}

	findOptions := optionsToMongoUserFindOptions(options)

	fmt.Println("filter => ", filter)

	var users []*entity.User

	cursor, err := m.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		model := &po.User{}
		if err := cursor.Decode(model); err != nil {
			return nil, err
		}
		e, err := converter.UserPOToEntity(model)
		if err != nil {
			return nil, err
		}
		users = append(users, e)
	}

	return users, nil
}

func optionsToMongoUserFindOptions(efo *entity.FindOptions) *options.FindOptions {
	findOptions := &options.FindOptions{}

	return findOptions
}
