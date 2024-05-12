package po

type Network struct {
	ID        string   `bson:"_id,omitempty"`
	Name      string   `bson:"name"`
	Cidr      string   `bson:"cidr"`
	UsedIPs   []string `bson:"used_ips"`
	Status    uint     `bson:"status"`
	CreatedAt int64    `bson:"created_at"`
	UpdatedAt int64    `bson:"updated_at"`
	DeletedAt int64    `bson:"deleted_at"`

	//AvailableIPs []string `bson:"available_ips"`
}
