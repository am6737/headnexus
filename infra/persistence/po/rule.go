package po

type Rule struct {
	ID          string   `bson:"_id"`
	UserID      string   `bson:"user_id"`
	Name        string   `bson:"name"`
	Description string   `bson:"description"`
	HostID      string   `bson:"host_id"`
	Port        string   `bson:"port"`
	Proto       string   `bson:"proto"`
	Host        []string `bson:"host,omitempty"`
	CreatedAt   int64    `bson:"created_at"`
	UpdatedAt   int64    `bson:"updated_at"`
	DeletedAt   int64    `bson:"deleted_at"`
	Type        uint8    `bson:"type"`
	Action      uint8    `bson:"action"`
}
