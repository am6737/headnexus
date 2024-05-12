package persistence

import (
	"context"
	"errors"
	"fmt"
	ctime "github.com/am6737/headnexus/common/time"
	"github.com/am6737/headnexus/domain/user/entity"
	"github.com/am6737/headnexus/domain/user/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	userCollectionName = "users"
)

type UserModel struct {
	ID           string `bson:"_id,omitempty"`
	Name         string `bson:"name"`
	Email        string `bson:"email"`
	Token        string `bson:"token"`
	Password     string `bson:"password"`
	Status       uint   `bson:"status"`
	Verification string `bson:"verification"`
	CreatedAt    int64  `bson:"created_at"`
	UpdatedAt    int64  `bson:"updated_at"`
	DeletedAt    int64  `bson:"deleted_at"`
}

func (m *UserModel) FromEntity(e *entity.User) error {
	var id primitive.ObjectID
	if e.ID != "" {
		fmt.Println("id:", e.ID)
		oid, err := primitive.ObjectIDFromHex(e.ID)
		if err != nil {
			return err
		}
		id = oid
	}
	m.ID = id.Hex()
	m.Name = e.Name
	m.Email = e.Email
	m.Token = e.Token
	m.Status = uint(e.Status)
	m.Verification = e.Verification
	m.Password = e.Password
	return nil
}

func (m *UserModel) ToEntity() (*entity.User, error) {
	return &entity.User{
		ID:           m.ID,
		Name:         m.Name,
		Email:        m.Email,
		Token:        m.Token,
		Password:     m.Password,
		Status:       entity.UserStatus(m.Status),
		Verification: m.Verification,
	}, nil
}

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
	model := &UserModel{}
	if err := model.FromEntity(user); err != nil {
		return nil, err // 返回错误，以及 nil 指针作为结果
	}

	currentTime := ctime.CurrentTimestampMillis()
	fmt.Println("currentTime:", currentTime)
	model.CreatedAt = currentTime
	model.UpdatedAt = currentTime
	model.ID = primitive.NewObjectID().Hex()

	fmt.Println("model:", model.ID)
	_, err := m.collection.InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}

	e, err := model.ToEntity()
	if err != nil {
		return nil, err
	}

	fmt.Println("创建完成")
	return e, nil
}

func (m *UserMongodbRepository) Get(ctx context.Context, id string) (*entity.User, error) {
	filter := bson.M{"_id": id}
	var model UserModel
	err := m.collection.FindOne(ctx, filter).Decode(&model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrorHostNotFound
		}
		return nil, err
	}
	user, err := model.ToEntity()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m *UserMongodbRepository) Update(ctx context.Context, user *entity.User) error {
	// 将 api.Host 转换为 HostDBModel
	model := &UserModel{}
	if err := model.FromEntity(user); err != nil {
		return err
	}

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": model}

	//opts := options.Update().SetUpsert(true)
	_, err := m.collection.UpdateOne(ctx, filter, update)
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
		var user UserModel
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		e, err := user.ToEntity()
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
