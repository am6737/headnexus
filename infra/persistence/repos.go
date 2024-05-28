package persistence

import (
	"context"
	"github.com/am6737/headnexus/domain/host/repository"
	networkRepo "github.com/am6737/headnexus/domain/network/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	HostRepo     repository.HostRepository
	RuleRepo     repository.RuleRepository
	HostRuleRepo repository.HostRuleRepository
	NetworkRepo  networkRepo.NetworkRepository

	client *mongo.Client
	dbName string
}

func NewRepositories(client *mongo.Client, dbName string) *Repositories {
	return &Repositories{
		HostRepo:     NewHostMongodbRepository(client, dbName),
		RuleRepo:     NewRuleMongodbRepository(client, dbName),
		HostRuleRepo: NewHostMongodbRepository(client, dbName),
		NetworkRepo:  NewNetworkMongoDBRepository(client, dbName),
	}
}

func (r *Repositories) TXRepositories(ctx context.Context, fc func(txr *Repositories, sessionCtx mongo.SessionContext) error) error {
	// 开始事务
	session, err := r.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	// 开始事务
	if err := session.StartTransaction(); err != nil {
		return err
	}

	// 设置事务选项
	//txnOpts := options.Transaction().SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

	// 执行带有会话上下文的操作
	err = mongo.WithSession(ctx, session, func(sessionCtx mongo.SessionContext) error {
		// 传递拥有相同会话的新 Repository 对象
		txr := &Repositories{
			HostRepo:     NewHostMongodbRepository(r.client, r.dbName),
			RuleRepo:     NewRuleMongodbRepository(r.client, r.dbName),
			HostRuleRepo: NewHostMongodbRepository(r.client, r.dbName),
		}

		// 执行提供的函数，并传递带有会话上下文的 Repository 对象
		if err := fc(txr, sessionCtx); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		// 回滚事务
		if err := session.AbortTransaction(ctx); err != nil {
			return err
		}
		return err
	}

	// 提交事务
	if err := session.CommitTransaction(ctx); err != nil {
		return err
	}

	return nil
}

func (r *Repositories) Automigrate() error {
	return nil
}
