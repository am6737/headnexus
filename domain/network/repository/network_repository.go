package repository

import (
	"context"
	"github.com/am6737/headnexus/domain/network/entity"
)

type NetworkRepository interface {
	Create(ctx context.Context, network *entity.Network) (*entity.Network, error)
	Get(ctx context.Context, id string) (*entity.Network, error)
	Find(ctx context.Context, query *entity.QueryNetwork) ([]*entity.Network, error)
	Update(ctx context.Context, network *entity.Network) (*entity.Network, error)
	Delete(ctx context.Context, id string) error

	// UsedIPs 获取指定网络的已使用的地址
	UsedIPs(ctx context.Context, id string) ([]string, error)
	// AvailableIPs 获取指定网络的可用地址
	AvailableIPs(ctx context.Context, id string) ([]string, error)
	// AllocateIP 分配地址给指定网络
	AllocateIP(ctx context.Context, id string, ip string) error
	// ReleaseIP 释放指定网络的地址
	ReleaseIP(ctx context.Context, id string, ip string) error

	// AutoAllocateIP 自动分配一个地址给指定网络
	AutoAllocateIP(ctx context.Context, id string) (string, error)
}
