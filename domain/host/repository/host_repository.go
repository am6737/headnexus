package repository

import (
	"context"
	"github.com/am6737/headnexus/domain/host/entity"
)

// HostRepository 定义了对主机数据的存储和操作接口
type HostRepository interface {
	Create(ctx context.Context, host *entity.Host) (*entity.Host, error)
	Get(ctx context.Context, id string) (*entity.Host, error)
	Update(ctx context.Context, host *entity.Host) (*entity.Host, error)
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, options *entity.HostFindOptions) ([]*entity.Host, error)

	GetHostByEnrollCode(ctx context.Context, code string) (*entity.Host, error)

	// GetEnrollHost 获取主机的注册信息
	GetEnrollHost(ctx context.Context, getEnrollHost *entity.GetEnrollHost) (*entity.EnrollHost, error)
	// EnrollHost 注册主机
	EnrollHost(ctx context.Context, enrollHost *entity.EnrollHost) error

	// HostOnline 主机上线
	HostOnline(ctx context.Context, hostOnline *entity.HostOnline) (*entity.HostOnline, error)
	// HostOffline 主机下线
	HostOffline(ctx context.Context, hostOffline *entity.HostOffline) (*entity.HostOffline, error)
}
