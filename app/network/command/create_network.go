package command

import (
	"context"
	"errors"
	"github.com/am6737/headnexus/app/network"
	"github.com/am6737/headnexus/domain/network/entity"
	"github.com/segmentio/ksuid"
)

func (h *NetworkHandler) Create(ctx context.Context, cmd *network.CreateNetwork) (*entity.Network, error) {
	if cmd == nil {
		return nil, errors.New("cmd is nil")
	}
	//if cmd.Cidr == "" {
	//	return nil, errors.New("cidr is empty")
	//}
	//if cmd.Name == "" {
	//	cmd.Name = ksuid.New().String()
	//}
	//namePattern := regexp.MustCompile(`^\S{1,30}$`)
	//if !namePattern.MatchString(cmd.Name) {
	//	return nil, errors.New("network name should be 1 to 30 characters long and cannot contain spaces")
	//}

	// 检查是否已经存在具有相同CIDR或名称的网络
	//existingNetworks, err := n.repo.Find(ctx, &api.QueryNetwork{Name: cmd.Name, Cidr: cmd.Cidr})
	//if err != nil {
	//	return nil, err
	//}
	//if len(existingNetworks) > 0 {
	//	return nil, errors.New("network with the same CIDR or name already exists")
	//}

	//fmt.Println("existingNetworks => ", existingNetworks)

	id := ksuid.New().String()
	if cmd.Name == "" {
		cmd.Name = id
	}
	network := &entity.Network{
		Name: cmd.Name,
		Cidr: cmd.Cidr,
		ID:   id,
	}
	r, err := h.repo.Create(ctx, network)
	if err != nil {
		return nil, err
	}

	return &entity.Network{
		ID:   r.ID,
		Name: r.Name,
		Cidr: r.Cidr,
		//CreatedAt: time.Unix(r.CreatedAt, 0).Format("2006-01-02 15:04:05"),
		CreatedAt: r.CreatedAt,
	}, nil
}
