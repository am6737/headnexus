package persistence

import (
	"context"
	"errors"
	ctime "github.com/am6737/headnexus/common/time"
	"github.com/am6737/headnexus/domain/network/entity"
	"github.com/am6737/headnexus/domain/network/repository"
	net2 "github.com/am6737/headnexus/pkg/net"
	"net"
	"sort"
	"sync"
)

var (
	ErrorNetworkNotFound          = errors.New("network not found")
	ErrorNetworkExists            = errors.New("network already exists")
	ErrorAddressNotInCIDRRange    = errors.New("address is not within network CIDR range")
	ErrorAddressAlreadyAllocated  = errors.New("address already allocated")
	ErrorAddressNotFoundInUsedIPs = errors.New("address not found in used IPs")
	ErrorCIDRInUse                = errors.New("cidr already in use")
	ErrNoAvailableIP              = errors.New("no available IP addresses in the network")
)

type NetWorkModel struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Cidr      string `json:"cidr"`
	Status    uint   `gorm:"comment:网络状态(0=正常状态, 1=锁定状态, 2=删除状态)"`
	CreatedAt int64  `gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt int64  `gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt int64  `gorm:"default:0;comment:删除时间"`

	UsedIPs      []string // 存储分配的地址
	AvailableIPs []string // 存储可用的地址
}

func (m *NetWorkModel) From(e *entity.Network) error {
	if e == nil {
		return errors.New("input network is nil")
	}

	m.ID = e.ID
	m.Name = e.Name
	m.Cidr = e.Cidr
	return nil
}

func (m *NetWorkModel) To() (*entity.Network, error) {
	if m == nil {
		return nil, errors.New("network model is nil")
	}

	entityNetwork := &entity.Network{
		ID:   m.ID,
		Name: m.Name,
		Cidr: m.Cidr,
	}
	return entityNetwork, nil
}

var _ repository.NetworkRepository = &NetworkMemoryRepository{}

func NewNetworkMemoryRepository() *NetworkMemoryRepository {
	return &NetworkMemoryRepository{
		lock:          &sync.RWMutex{},
		networks:      make(map[string]*NetWorkModel),
		cidrToNetwork: make(map[string]string),
	}
}

type NetworkMemoryRepository struct {
	lock     *sync.RWMutex
	networks map[string]*NetWorkModel

	cidrToNetwork map[string]string // CIDR到网络名称的映射
}

func (m *NetworkMemoryRepository) AutoAllocateIP(ctx context.Context, id string) (string, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// 检查网络是否存在
	model, exists := m.networks[id]
	if !exists {
		return "", ErrorNetworkNotFound
	}

	// 检查网络的可用地址列表是否为空
	if len(model.AvailableIPs) == 0 {
		return "", ErrNoAvailableIP
	}

	// 从可用地址列表中获取第一个地址并从列表中移除
	allocatedIP := model.AvailableIPs[0]
	model.AvailableIPs = model.AvailableIPs[1:]

	// 更新网络的信息
	model.UpdatedAt = ctime.CurrentTimestampMillis()
	m.networks[id] = model

	return allocatedIP, nil
}

func (m *NetworkMemoryRepository) UsedIPs(ctx context.Context, id string) ([]string, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	network, exists := m.networks[id]
	if !exists {
		return nil, ErrorNetworkNotFound
	}

	return network.UsedIPs, nil
}

func (m *NetworkMemoryRepository) AllocateIP(ctx context.Context, id string, ip string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	network, exists := m.networks[id]
	if !exists {
		return ErrorNetworkNotFound
	}

	if !net2.IsIPInCIDR(ip, network.Cidr) {
		return ErrorAddressNotInCIDRRange
	}

	for _, usedIP := range network.UsedIPs {
		if usedIP == ip {
			return ErrorAddressAlreadyAllocated
		}
	}

	// 将已分配的 IP 地址从可用 IP 地址列表中移除
	for i, availableIP := range network.AvailableIPs {
		if availableIP == ip {
			network.AvailableIPs = append(network.AvailableIPs[:i], network.AvailableIPs[i+1:]...)
			break
		}
	}

	network.UsedIPs = append(network.UsedIPs, ip)
	return nil
}

func (m *NetworkMemoryRepository) ReleaseIP(ctx context.Context, id string, ip string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	network, exists := m.networks[id]
	if !exists {
		return ErrorNetworkNotFound
	}

	if !net2.IsIPInCIDR(ip, network.Cidr) {
		return ErrorAddressNotInCIDRRange
	}

	found := false
	for i, usedIP := range network.UsedIPs {
		if usedIP == ip {
			network.UsedIPs = append(network.UsedIPs[:i], network.UsedIPs[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return ErrorAddressNotFoundInUsedIPs
	}

	// 将地址添加回可用地址列表中
	network.AvailableIPs = append(network.AvailableIPs, ip)

	return nil
}

func (m *NetworkMemoryRepository) AvailableIPs(ctx context.Context, id string) ([]string, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	network, exists := m.networks[id]
	if !exists {
		return nil, ErrorNetworkNotFound
	}

	return network.AvailableIPs, nil
}

func (m *NetworkMemoryRepository) Create(ctx context.Context, network *entity.Network) (*entity.Network, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	model := &NetWorkModel{}
	if err := model.From(network); err != nil {
		return nil, err
	}

	if _, exists := m.networks[model.ID]; exists {
		return nil, ErrorNetworkExists
	}

	for _, existingNetwork := range m.networks {
		if existingNetwork.Cidr == network.Cidr {
			return nil, ErrorCIDRInUse
		}
	}

	// 生成可用的 IP 地址列表
	availableIPs, err := generateAvailableIPs(network.Cidr)
	if err != nil {
		return nil, err
	}
	model.AvailableIPs = availableIPs

	currentTime := ctime.CurrentTimestampMillis()
	model.CreatedAt = currentTime
	model.UpdatedAt = currentTime
	m.networks[model.ID] = model
	network.CreatedAt = currentTime
	return network, nil
}

func (m *NetworkMemoryRepository) Get(ctx context.Context, id string) (*entity.Network, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	model, exists := m.networks[id]
	if !exists {
		return nil, ErrorNetworkNotFound
	}

	network, err := model.To()
	if err != nil {
		return nil, err
	}

	return network, nil
}

func (m *NetworkMemoryRepository) Find(ctx context.Context, query *entity.QueryNetwork) ([]*entity.Network, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// 先将所有网络转换为模型
	var models []*NetWorkModel
	for _, model := range m.networks {
		models = append(models, model)
	}

	// 根据名称过滤网络
	if query.Name != "" {
		filteredModels := make([]*NetWorkModel, 0)
		for _, model := range models {
			if model.Name == query.Name {
				filteredModels = append(filteredModels, model)
			}
		}
		models = filteredModels
	}

	// 根据CIDR过滤网络
	if query.Cidr != "" {
		filteredModels := make([]*NetWorkModel, 0)
		for _, model := range models {
			if model.Cidr == query.Cidr {
				filteredModels = append(filteredModels, model)
			}
		}
		models = filteredModels
	}

	// 根据排序方向进行排序
	if query.SortDirection == "asc" {
		sort.Slice(models, func(i, j int) bool {
			return models[i].CreatedAt < models[j].CreatedAt
		})
	} else if query.SortDirection == "desc" {
		sort.Slice(models, func(i, j int) bool {
			return models[i].CreatedAt > models[j].CreatedAt
		})
	}

	if len(models) == 0 {
		return []*entity.Network{}, nil
	}

	// 根据分页参数进行分页
	start := query.PageSize * (query.PageNumber - 1)
	if start < 0 || start >= len(models) {
		return []*entity.Network{}, nil
	}

	end := start + query.PageSize
	if end > len(models) {
		end = len(models)
	}

	// 如果开始索引大于等于结束索引，则返回空结果集
	if start >= end {
		return []*entity.Network{}, nil
	}

	// 将模型转换为 entity.Network 类型并返回
	var result []*entity.Network
	for _, model := range models[start:end] {
		network, err := model.To()
		if err != nil {
			return nil, err
		}
		result = append(result, network)
	}

	return result, nil
}

func (m *NetworkMemoryRepository) Update(ctx context.Context, e *entity.Network) (*entity.Network, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	model := &NetWorkModel{}
	if err := model.From(e); err != nil {
		return nil, err
	}

	network, exists := m.networks[e.ID]
	if !exists {
		return nil, ErrorNetworkNotFound
	}

	network.Name = e.Name
	network.UpdatedAt = ctime.CurrentTimestampMillis()

	r, err := network.To()
	if err != nil {
		return nil, err
	}

	//m.networks[network.ID] = model
	return r, nil
}

func (m *NetworkMemoryRepository) Delete(ctx context.Context, id string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, exists := m.networks[id]; !exists {
		return ErrorNetworkNotFound
	}

	delete(m.networks, id)
	return nil
}

// isIPAllocated checks if the given IP address is allocated in the network.
func (m *NetworkMemoryRepository) isIPAllocated(network *NetWorkModel, ip string) bool {
	for _, usedIP := range network.UsedIPs {
		if usedIP == ip {
			return true
		}
	}
	return false
}

// 生成可用的 IP 地址列表
func generateAvailableIPs(cidr string) ([]string, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	availableIPs := make([]string, 0)
	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		availableIPs = append(availableIPs, ip.String())
	}

	// 排除网络地址和广播地址
	availableIPs = availableIPs[1 : len(availableIPs)-1]

	return availableIPs, nil
}

// inc increments IP address by one.
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
