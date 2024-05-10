package persistence

import (
	"context"
	"github.com/am6737/headnexus/domain/user/entity"
	"github.com/am6737/headnexus/domain/user/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserModel struct {
	ID        string `bson:"_id,omitempty"`
	Name      string `bson:"name"`
	Email     string `bson:"email"`
	Token     string `bson:"token"`
	Password  string `bson:"password"`
	CreatedAt int64  `bson:"created_at"`
	UpdatedAt int64  `bson:"updated_at"`
	DeletedAt int64  `bson:"deleted_at"`
}

//func (m *RuleModel) FromEntity(e *entity.Rule) error {
//	var id primitive.ObjectID
//	if e.ID != "" {
//		oid, err := primitive.ObjectIDFromHex(e.ID)
//		if err != nil {
//			return err
//		}
//		id = oid
//	}
//	m.ID = id
//	m.Name = e.Name
//	m.Description = e.Description
//	m.HostID = e.HostID
//	m.Port = e.Port
//	m.Proto = string(e.Proto)
//	m.Action = uint8(e.Action)
//	m.Host = e.Host
//	return nil
//}
//
//func (m *RuleModel) ToEntity() (*entity.Rule, error) {
//	return &entity.Rule{
//		ID:          m.ID.Hex(),
//		Type:        entity.RuleType(m.Type),
//		CreatedAt:   m.CreatedAt,
//		Name:        m.Name,
//		Description: m.Description,
//		HostID:      m.HostID,
//		Port:        m.Port,
//		Proto:       entity.RuleProto(m.Proto),
//		Action:      entity.RuleAction(m.Action),
//		Host:        m.Host,
//	}, nil
//}

func (m *UserModel) FromEntity(e *entity.User) error {
	var id primitive.ObjectID
	if e.ID != "" {
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
	m.Password = e.Password
	return nil
}

var _ repository.UserRepository = &UserMongodbRepository{}

type UserMongodbRepository struct {
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
}

func NewUserMongodbRepository(client *mongo.Client, dbName string) *UserMongodbRepository {
	db := client.Database(dbName)
	collection := db.Collection(ruleCollectionName)
	m := &UserMongodbRepository{
		client:     client,
		db:         db,
		collection: collection,
	}
	return m
}

func (u *UserMongodbRepository) Create(ctx context.Context, network *entity.User) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserMongodbRepository) Get(ctx context.Context, id string) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserMongodbRepository) Update(ctx context.Context, network *entity.User) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserMongodbRepository) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
