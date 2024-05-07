package entity

type Network struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Cidr         string   `json:"cidr"`
	UsedIPs      []string `json:"used_ips"`      // 存储分配的地址
	AvailableIPs []string `json:"available_ips"` // 存储可用的地址
	CreatedAt    int64    `json:"created_at"`
}

type QueryNetwork struct {
	Name          string
	Cidr          string
	SortDirection string // 排序方向
	IncludeCounts bool
	PageSize      int // 每页大小
	PageNumber    int // 页码
}
