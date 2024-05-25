package entity

import (
	"github.com/am6737/headnexus/config"
)

// Host 表示一个主机
type Host struct {
	ID        string
	Name      string
	NetworkID string
	Role      string
	IPAddress string
	// StaticAddresses A list of static addresses for the host
	StaticAddresses []string

	Port         int
	IsLighthouse bool
	CreatedAt    int64
	LastSeenAt   int64
	EnrollAt     int64
	Status       HostStatus

	Tags map[string]interface{}

	// Config 主机的配置文件
	Config config.Config
}

// HostStatus 是主机的状态枚举类型。
type HostStatus int

const (
	// Offline 表示主机处于离线状态。
	Offline HostStatus = iota

	// Online 表示主机处于在线状态。
	Online
)

// String 方法返回 HostStatus 的字符串表示形式。
func (hs HostStatus) String() string {
	switch hs {
	case Online:
		return "Online"
	case Offline:
		return "Offline"
	default:
		return "Unknown"
	}
}

type GetEnrollHost struct {
	HostID string
}

type HostOnline struct {
	ID         string
	Status     HostStatus
	LastSeenAt int64
}

type HostOffline struct {
	ID         string
	Status     HostStatus
	LastSeenAt int64
}

type EnrollHost struct {
	HostID              string
	Code                string
	EnrollCodeExpiredAt int64
	EnrollAt            int64
	CreatedAt           int64
}

// HostFindOptions 定义了查询主机数据时的过滤和排序选项
type HostFindOptions struct {
	// 可以添加各种过滤条件,如按名称、IP、标签等进行过滤
	Filters map[string]interface{}

	// 可以添加排序选项,如按创建时间、名称等进行排序
	Sort map[string]interface{} // 1 for ascending, -1 for descending

	// 分页选项
	Limit  int
	Offset int

	NetworkID    string
	IPAddress    string
	Role         string
	Name         string
	IsLighthouse bool
}

type ListHostRuleOptions struct {
	HostID   string
	Port     string
	Type     *RuleType
	Proto    *RuleProto
	Action   *RuleAction
	PageSize int
	PageNum  int
}
