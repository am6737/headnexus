package converter

import (
	"github.com/am6737/headnexus/domain/host/entity"
	"github.com/am6737/headnexus/infra/persistence/po"
)

func HostEntityToPO(host *entity.Host) (*po.Host, error) {
	m := &po.Host{}
	m.Owner = host.Owner
	m.ID = host.ID
	m.Name = host.Name
	m.NetworkID = host.NetworkID
	m.IPAddress = host.IPAddress
	m.PublicIP = host.PublicIP
	m.Role = host.Role
	m.Port = host.Port
	m.Tags = host.Tags
	m.Config = host.Config
	return m, nil
}

func HostPoToEntity(po *po.Host) (*entity.Host, error) {
	host := &entity.Host{
		Owner:      po.Owner,
		ID:         po.ID,
		Name:       po.Name,
		NetworkID:  po.NetworkID,
		IPAddress:  po.IPAddress,
		PublicIP:   po.PublicIP,
		Role:       po.Role,
		Port:       po.Port,
		LastSeenAt: po.LastSeenAt,
		EnrollAt:   po.EnrollAt,
		Tags:       po.Tags,
		Config:     po.Config,
		Status:     entity.HostStatus(po.Status),
		CreatedAt:  po.CreatedAt,
	}
	return host, nil
}
