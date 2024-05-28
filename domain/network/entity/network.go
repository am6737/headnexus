package entity

import "github.com/am6737/headnexus/pkg/net"

type Network struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Cidr         string   `json:"cidr"`
	UsedIPs      []string `json:"used_ips"`      // 存储分配的地址
	AvailableIPs []string `json:"available_ips"` // 存储可用的地址
	CreatedAt    int64    `json:"created_at"`
}

func (n Network) Mask() string {
	mask, err := net.GenerateMask(n.Cidr)
	if err != nil {
		return ""
	}
	return mask
}

type QueryNetwork struct {
	Name          string
	Cidr          string
	SortDirection string // 排序方向
	IncludeCounts bool
	PageSize      int // 每页大小
	PageNumber    int // 页码
}
