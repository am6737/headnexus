package persistence

import (
	"context"
	"fmt"
	"github.com/am6737/headnexus/domain/network/entity"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestNetworkMemoryRepository_CreateAndGet(t *testing.T) {
	repo := NewNetworkMemoryRepository()

	// 测试创建网络
	network := &entity.Network{
		ID:   "1",
		Name: "TestNetwork",
		Cidr: "192.168.0.0/24",
	}
	createdNetwork, err := repo.Create(context.Background(), network)
	assert.NoError(t, err)
	assert.NotNil(t, createdNetwork)

	// 测试获取网络
	retrievedNetwork, err := repo.Get(context.Background(), "1")
	assert.NoError(t, err)
	assert.NotNil(t, retrievedNetwork)
	assert.Equal(t, "1", retrievedNetwork.ID)
	assert.Equal(t, "TestNetwork", retrievedNetwork.Name)
	assert.Equal(t, "192.168.0.0/24", retrievedNetwork.Cidr)
}

func TestNetworkMemoryRepository_Update(t *testing.T) {
	repo := NewNetworkMemoryRepository()

	// 创建网络
	network := &entity.Network{
		ID:   "1",
		Name: "TestNetwork",
		Cidr: "192.168.0.0/24",
	}
	_, err := repo.Create(context.Background(), network)
	assert.NoError(t, err)

	// 更新网络
	updatedNetwork := &entity.Network{
		ID:   "1",
		Name: "UpdatedNetwork",
	}
	updatedNetwork, err = repo.Update(context.Background(), updatedNetwork)
	assert.NoError(t, err)
	assert.NotNil(t, updatedNetwork)
	assert.Equal(t, "1", updatedNetwork.ID)
	assert.Equal(t, "UpdatedNetwork", updatedNetwork.Name)
}

func TestNetworkMemoryRepository_Delete(t *testing.T) {
	repo := NewNetworkMemoryRepository()

	// 创建网络
	network := &entity.Network{
		ID:   "1",
		Name: "TestNetwork",
		Cidr: "192.168.0.0/24",
	}
	_, err := repo.Create(context.Background(), network)
	assert.NoError(t, err)

	// 删除网络
	err = repo.Delete(context.Background(), "1")
	assert.NoError(t, err)

	// 确认网络已删除
	_, err = repo.Get(context.Background(), "1")
	assert.Error(t, err)
	assert.Equal(t, ErrorNetworkNotFound, err)
}

func TestNetworkMemoryRepository_Find(t *testing.T) {
	repo := NewNetworkMemoryRepository()

	// 创建多个网络
	for i := 0; i < 20; i++ {
		network := &entity.Network{
			ID:   strconv.Itoa(i),
			Name: fmt.Sprintf("Network%d", i),
			Cidr: fmt.Sprintf("10.0.%d.0/24", i),
		}
		_, err := repo.Create(context.Background(), network)
		assert.NoError(t, err)
	}

	// 分页查询网络
	query := &entity.QueryNetwork{
		SortDirection: "asc",
		PageSize:      10,
		PageNumber:    1,
	}
	networks, err := repo.Find(context.Background(), query)
	assert.NoError(t, err)
	assert.Len(t, networks, 10)

	query.PageNumber = 2
	networks, err = repo.Find(context.Background(), query)
	assert.NoError(t, err)
	assert.Len(t, networks, 10)

	query.PageNumber = 3
	networks, err = repo.Find(context.Background(), query)
	assert.NoError(t, err)
	assert.Len(t, networks, 0)
}

func TestNetworkMemoryRepository_CreateDuplicate(t *testing.T) {
	repo := NewNetworkMemoryRepository()

	// 创建网络
	network := &entity.Network{
		ID:   "1",
		Name: "TestNetwork",
		Cidr: "192.168.0.0/24",
	}
	_, err := repo.Create(context.Background(), network)
	assert.NoError(t, err)

	// 再次创建相同ID的网络，应该返回错误
	_, err = repo.Create(context.Background(), network)
	assert.Error(t, err)
	assert.Equal(t, ErrorNetworkExists, err)
}

func TestNetworkMemoryRepository_UsedIPs(t *testing.T) {
	repo := NewNetworkMemoryRepository()

	// 创建网络
	network := &entity.Network{
		ID:   "1",
		Name: "TestNetwork",
		Cidr: "192.168.0.0/24",
	}
	_, err := repo.Create(context.Background(), network)
	assert.NoError(t, err)

	// 分配几个地址
	allocatedIPs := []string{
		"192.168.0.2",
		"192.168.0.5",
		"192.168.0.10",
	}
	for _, ip := range allocatedIPs {
		err := repo.AllocateIP(context.Background(), "1", ip)
		assert.NoError(t, err)
	}

	// 获取已分配的地址
	usedIPs, err := repo.UsedIPs(context.Background(), "1")
	assert.NoError(t, err)
	assert.ElementsMatch(t, allocatedIPs, usedIPs)
}

func TestNetworkMemoryRepository_AvailableIPs(t *testing.T) {
	repo := NewNetworkMemoryRepository()

	// 创建网络
	network := &entity.Network{
		ID:   "1",
		Name: "TestNetwork",
		Cidr: "192.168.0.0/24",
	}
	_, err := repo.Create(context.Background(), network)
	assert.NoError(t, err)

	// 初始化可用地址列表
	ps, err := generateAvailableIPs(network.Cidr)
	assert.NoError(t, err)
	network.AvailableIPs = ps
	assert.NoError(t, err)

	// 分配几个地址
	allocatedIPs := []string{
		"192.168.0.2",
		"192.168.0.5",
		"192.168.0.10",
	}
	fmt.Println(1)
	for _, ip := range allocatedIPs {
		err := repo.AllocateIP(context.Background(), "1", ip)
		assert.NoError(t, err)
	}
	fmt.Println(2)

	// 获取可用地址
	availableIPs, err := repo.AvailableIPs(context.Background(), "1")
	assert.NoError(t, err)
	fmt.Println(3)

	// 确保所有已分配的地址不在可用地址列表中
	for _, allocatedIP := range allocatedIPs {
		assert.NotContains(t, availableIPs, allocatedIP)
	}
	fmt.Println(4)

	// 确保可用地址的数量等于可分配地址数量减去已分配地址数量
	ps2, err := generateAvailableIPs(network.Cidr)
	assert.NoError(t, err)
	expectedAvailableIPCount := len(ps2) - len(allocatedIPs)
	assert.Len(t, availableIPs, expectedAvailableIPCount)
}

func TestNetworkMemoryRepository_AllocateIP(t *testing.T) {
	repo := NewNetworkMemoryRepository()

	// 创建网络
	network := &entity.Network{
		ID:   "1",
		Name: "TestNetwork",
		Cidr: "192.168.0.0/24",
	}
	_, err := repo.Create(context.Background(), network)
	assert.NoError(t, err)

	// 分配地址
	err = repo.AllocateIP(context.Background(), "1", "192.168.0.2")
	assert.NoError(t, err)

	// 再次分配相同的地址，应该返回错误
	err = repo.AllocateIP(context.Background(), "1", "192.168.0.2")
	assert.Error(t, err)
	assert.Equal(t, ErrorAddressAlreadyAllocated, err)

	// 分配不在 CIDR 范围内的地址，应该返回错误
	err = repo.AllocateIP(context.Background(), "1", "10.0.0.2")
	assert.Error(t, err)
	assert.Equal(t, ErrorAddressNotInCIDRRange, err)
}

func TestNetworkMemoryRepository_ReleaseIP(t *testing.T) {
	repo := NewNetworkMemoryRepository()

	// 创建网络
	network := &entity.Network{
		ID:   "1",
		Name: "TestNetwork",
		Cidr: "192.168.0.0/24",
	}
	_, err := repo.Create(context.Background(), network)
	assert.NoError(t, err)

	// 分配地址
	err = repo.AllocateIP(context.Background(), "1", "192.168.0.2")
	assert.NoError(t, err)

	// 释放已分配的地址
	err = repo.ReleaseIP(context.Background(), "1", "192.168.0.2")
	assert.NoError(t, err)

	// 尝试释放未分配的地址，应该返回错误
	err = repo.ReleaseIP(context.Background(), "1", "192.168.0.3")
	assert.Error(t, err)
	assert.Equal(t, ErrorAddressNotFoundInUsedIPs, err)
}
