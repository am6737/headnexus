package converter

import (
	"github.com/am6737/headnexus/domain/host/entity"
	"github.com/am6737/headnexus/infra/persistence/po"
)

func HostEntityToPO(host *entity.Host) (*po.Host, error) {
	m := &po.Host{}
	m.ID = host.ID
	m.Name = host.Name
	m.NetworkID = host.NetworkID
	m.IPAddress = host.IPAddress
	m.Role = host.Role
	m.Port = host.Port
	m.IsLighthouse = host.IsLighthouse
	m.StaticAddresses = host.StaticAddresses
	m.Tags = host.Tags
	m.Config = host.Config
	return m, nil
}

func HostPOToEntity(po *po.Host) (*entity.Host, error) {
	host := &entity.Host{
		ID:              po.ID,
		Name:            po.Name,
		NetworkID:       po.NetworkID,
		IPAddress:       po.IPAddress,
		Role:            po.Role,
		Port:            po.Port,
		IsLighthouse:    po.IsLighthouse,
		StaticAddresses: po.StaticAddresses,
		LastSeenAt:      po.LastSeenAt,
		Tags:            po.Tags,
		Config:          po.Config,
		Status:          entity.HostStatus(po.Status),
		//CreatedAt:       po.CreatedAt,
	}
	return host, nil
}
