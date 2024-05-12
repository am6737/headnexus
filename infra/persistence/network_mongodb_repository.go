package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	ctime "github.com/am6737/headnexus/common/time"
	"github.com/am6737/headnexus/domain/network/entity"
	"github.com/am6737/headnexus/infra/persistence/po"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/big"
	"net"
)

const (
	networkCollection = "networks"
)

type MongoDBRepository struct {
	client   *mongo.Client
	db       *mongo.Database
	networks *mongo.Collection
}

func NewNetworkMongoDBRepository(client *mongo.Client, dbName string) *MongoDBRepository {
	db := client.Database(dbName)
	networks := db.Collection(networkCollection)
	return &MongoDBRepository{
		client:   client,
		db:       db,
		networks: networks,
	}
}

func (m *MongoDBRepository) AutoAllocateIP(ctx context.Context, id string) (string, error) {
	network, err := m.getNetworkByID(ctx, id)
	if err != nil {
		return "", err
	}

	// 解析子网的CIDR
	_, subnet, err := net.ParseCIDR(network.Cidr)
	if err != nil {
		return "", err
	}

	// 获取子网的可用地址数量
	ones, bits := subnet.Mask.Size()
	availableIPCount := 1 << (bits - ones)

	// 如果已分配的IP数量等于子网的可用地址数量，则返回错误
	if len(network.UsedIPs) == availableIPCount {
		return "", ErrNoAvailableIP
	}

	// 如果没有已分配的IP，则返回子网的第一个IP地址
	if len(network.UsedIPs) == 0 {
		// 将网络号加1，得到第一个可用IP地址
		firstIPInt := big.NewInt(0).SetBytes(subnet.IP.To4())
		firstIPInt.Add(firstIPInt, big.NewInt(1))
		firstIPBytes := firstIPInt.Bytes()

		// 将第一个可用IP地址转换为字符串形式
		ip := net.IP(firstIPBytes).String()

		// 将新分配的IP添加到UsedIPs中
		network.UsedIPs = append(network.UsedIPs, ip)
		if err := m.updateNetwork(ctx, network); err != nil {
			return "", err
		}

		// 返回子网的第一个IP地址
		return ip, nil
	}

	// 获取最后一个已分配的IP
	lastAllocatedIP := network.UsedIPs[len(network.UsedIPs)-1]

	// 尝试解析字符串为net.IP类型
	ip := net.ParseIP(lastAllocatedIP)
	if ip == nil {
		return "", errors.New("invalid IP address")
	}

	// 从最后一个已分配的IP地址开始，尝试逐个加1，直到找到一个合适的IP地址
	for {
		// 将IP地址加1
		// 递增IP地址的字节表示
		incrementIP(&ip)

		// 检查新的IP地址是否已经存在于UsedIPs中
		ipStr := ip.String()
		if containsIP(network.UsedIPs, ipStr) {
			continue // 如果存在，继续尝试下一个IP地址
		}

		// 检查新的IP地址是否超过了当前子网的范围
		if !isIPInCIDR(ip.String(), network.Cidr) {
			return "", ErrorAddressNotInCIDRRange
		}

		// 将新分配的IP添加到UsedIPs中
		network.UsedIPs = append(network.UsedIPs, ip.String())
		if err := m.updateNetwork(ctx, network); err != nil {
			return "", err
		}

		// 返回新分配的IP
		return ip.String(), nil
	}
}

// 递增IP地址的字节表示
func incrementIP(ip *net.IP) {
	for i := len(*ip) - 1; i >= 0; i-- {
		(*ip)[i]++
		if (*ip)[i] > 0 {
			break
		}
	}
}

// isBroadcastIP 检查IP地址是否是子网的广播地址
//func isBroadcastIP(ip string, subnet *net.IPNet) bool {
//	// 获取子网的广播地址
//	broadcastIP := subnet.IP.Mask(subnet.Mask).To4()
//	if broadcastIP == nil {
//		return false // 广播地址需要为IPv4地址
//	}
//
//	// 获取广播地址的字节表示形式
//	broadcastIPBytes := broadcastIP.To4()
//
//	// 获取子网的广播地址
//	subnetBroadcast := net.IP(broadcastIPBytes)
//
//	// 检查IP地址是否等于广播地址
//	return ip == string(subnetBroadcast)
//}

// containsIP 检查IP地址是否存在于给定的IP地址列表中
func containsIP(ips []string, ip string) bool {
	for _, i := range ips {
		if i == ip {
			return true
		}
	}
	return false
}

//func (m *MongoDBRepository) AutoAllocateIP(ctx context.Context, id string) (string, error) {
//
//	network, err := m.getNetworkByID(ctx, id)
//	if err != nil {
//		return 0, err
//	}
//
//	if len(network.AvailableIPs) == 0 {
//		return 0, ErrNoAvailableIP
//	}
//
//	allocatedIP := network.AvailableIPs[0]
//
//	for _, usedIP := range network.UsedIPs {
//		if usedIP == allocatedIP {
//			return 0, ErrorAddressAlreadyAllocated
//		}
//	}
//	network.UsedIPs = append(network.UsedIPs, allocatedIP)
//	network.AvailableIPs = network.AvailableIPs[1:]
//
//	if err := m.updateNetwork(ctx, network); err != nil {
//		return 0, err
//	}
//
//	ip, err := api.ParseVpnIp(allocatedIP)
//	if err != nil {
//		return 0, err
//	}
//
//	return ip, nil
//}

func (m *MongoDBRepository) UsedIPs(ctx context.Context, id string) ([]string, error) {
	network, err := m.getNetworkByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var usedIPs []string
	for _, usedIP := range network.UsedIPs {
		usedIPs = append(usedIPs, usedIP)
	}

	return usedIPs, nil
}

func (m *MongoDBRepository) AllocateIP(ctx context.Context, id string, ip string) error {
	network, err := m.getNetworkByID(ctx, id)
	if err != nil {
		return err
	}

	if !isIPInCIDR(ip, network.Cidr) {
		return ErrorAddressNotInCIDRRange
	}

	for _, usedIP := range network.UsedIPs {
		if usedIP == ip {
			return ErrorAddressAlreadyAllocated
		}
	}

	// Remove allocated IP from available IPs
	//for i, availableIP := range network.AvailableIPs {
	//	if availableIP == ip.String() {
	//		network.AvailableIPs = append(network.AvailableIPs[:i], network.AvailableIPs[i+1:]...)
	//		break
	//	}
	//}

	network.UsedIPs = append(network.UsedIPs, ip)
	if err := m.updateNetwork(ctx, network); err != nil {
		return err
	}

	return nil
}

func (m *MongoDBRepository) ReleaseIP(ctx context.Context, id string, ip string) error {

	network, err := m.getNetworkByID(ctx, id)
	if err != nil {
		return err
	}

	if !isIPInCIDR(ip, network.Cidr) {
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

	// Add address back to available IPs
	//network.AvailableIPs = append(network.AvailableIPs, ip.String())

	if err := m.updateNetwork(ctx, network); err != nil {
		return err
	}

	return nil
}

func (m *MongoDBRepository) AvailableIPs(ctx context.Context, id string) ([]string, error) {
	//network, err := m.getNetworkByID(ctx, id)
	//if err != nil {
	//	return nil, err
	//}

	var availableIPs []string
	//for _, availableIP := range network.AvailableIPs {
	//	ip, err := api.ParseVpnIp(availableIP)
	//	if err != nil {
	//		return nil, err
	//	}
	//	availableIPs = append(availableIPs, ip)
	//}

	return availableIPs, nil
}

func (m *MongoDBRepository) Create(ctx context.Context, network *entity.Network) (*entity.Network, error) {

	existingNetwork, err := m.findNetworkByName(ctx, network.Name)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if existingNetwork != nil {
		return nil, ErrorNetworkExists
	}

	existingNetwork, err = m.findNetworkByCIDR(ctx, network.Cidr)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if existingNetwork != nil {
		return nil, ErrorCIDRInUse
	}

	currentTime := ctime.CurrentTimestampMillis()
	model := &po.Network{
		ID:        network.ID,
		Name:      network.Name,
		Cidr:      network.Cidr,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Status:    0,
		UsedIPs:   []string{},
		//AvailableIPs: []string{},
	}

	// Generate available IP addresses
	//availableIPs, err := generateAvailableIPs(network.Cidr)
	//if err != nil {
	//	return nil, err
	//}
	//
	//for _, ip := range availableIPs {
	//	model.AvailableIPs = append(model.AvailableIPs, ip.String())
	//}

	//model.AvailableIPs = availableIPs

	marshal, err := json.Marshal(&model)
	if err != nil {
		return nil, err
	}

	fmt.Println("model => ", string(marshal))

	if _, err := m.networks.InsertOne(ctx, model); err != nil {
		return nil, err
	}

	network.CreatedAt = currentTime
	return network, nil
}

func (m *MongoDBRepository) Get(ctx context.Context, id string) (*entity.Network, error) {
	network, err := m.getNetworkByID(ctx, id)
	if err != nil {
		return nil, err
	}

	apiNetwork := &entity.Network{
		ID:        network.ID,
		Name:      network.Name,
		Cidr:      network.Cidr,
		CreatedAt: network.CreatedAt,
	}
	return apiNetwork, nil
}

func (m *MongoDBRepository) Find(ctx context.Context, query *entity.QueryNetwork) ([]*entity.Network, error) {
	filter := bson.M{}
	if query.Name != "" {
		filter["name"] = query.Name
	}
	if query.Cidr != "" {
		filter["cidr"] = query.Cidr
	}

	// Sorting
	options := options.Find()
	if query.SortDirection == "asc" {
		options.SetSort(bson.D{{"created_at", 1}})
	} else if query.SortDirection == "desc" {
		options.SetSort(bson.D{{"created_at", -1}})
	}

	cursor, err := m.networks.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var networks []*entity.Network
	for cursor.Next(ctx) {
		model := &po.Network{}
		if err := cursor.Decode(model); err != nil {
			return nil, err
		}
		apiNetwork := &entity.Network{
			ID:        model.ID,
			Name:      model.Name,
			Cidr:      model.Cidr,
			CreatedAt: model.CreatedAt,
		}
		networks = append(networks, apiNetwork)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return networks, nil
}

func (m *MongoDBRepository) Update(ctx context.Context, e *entity.Network) (*entity.Network, error) {

	network, err := m.getNetworkByID(ctx, e.ID)
	if err != nil {
		return nil, err
	}

	network.Name = e.Name
	network.UpdatedAt = ctime.CurrentTimestampMillis()

	_, err = m.networks.ReplaceOne(ctx, bson.M{"_id": network.ID}, network)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (m *MongoDBRepository) Delete(ctx context.Context, id string) error {
	count, err := m.networks.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if count.DeletedCount == 0 {
		return ErrorNetworkNotFound
	}

	return nil
}

func (m *MongoDBRepository) findNetworkByName(ctx context.Context, name string) (*po.Network, error) {
	filter := bson.M{"name": name}
	model := &po.Network{}
	err := m.networks.FindOne(ctx, filter).Decode(model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return model, nil
}

func (m *MongoDBRepository) findNetworkByCIDR(ctx context.Context, cidr string) (*po.Network, error) {
	filter := bson.M{"cidr": cidr}
	model := &po.Network{}
	err := m.networks.FindOne(ctx, filter).Decode(model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return model, nil
}

func (m *MongoDBRepository) getNetworkByID(ctx context.Context, id string) (*po.Network, error) {
	filter := bson.M{"_id": id}
	model := &po.Network{}
	err := m.networks.FindOne(ctx, filter).Decode(model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrorNetworkNotFound
		}
		return nil, err
	}
	return model, nil
}

func (m *MongoDBRepository) updateNetwork(ctx context.Context, network *po.Network) error {
	_, err := m.networks.ReplaceOne(ctx, bson.M{"_id": network.ID}, network)
	if err != nil {
		return err
	}
	return nil
}

// isIPAllocated checks if the given IP address is allocated in the network.
func (m *MongoDBRepository) isIPAllocated(network *po.Network, ip string) bool {
	for _, usedIP := range network.UsedIPs {
		if usedIP == ip {
			return true
		}
	}
	return false
}
