package po

import "github.com/am6737/headnexus/config"

type Host struct {
	ID              string                 `bson:"_id"`
	Name            string                 `bson:"name"`
	NetworkID       string                 `bson:"network_id"`
	IPAddress       string                 `bson:"ip_address"`
	Role            string                 `bson:"role"`
	EnrollCode      string                 `bson:"enroll_code"`
	Owner           string                 `bson:"owner"`
	Port            int                    `bson:"port"`
	IsLighthouse    bool                   `bson:"is_lighthouse"`
	StaticAddresses []string               `bson:"static_addresses"`
	Tags            map[string]interface{} `bson:"tags"`
	CreatedAt       int64                  `bson:"created_at"`
	UpdatedAt       int64                  `bson:"updated_at"`
	DeletedAt       int64                  `bson:"deleted_at"`
	LastSeenAt      int64                  `bson:"last_seen_at"`
	EnrollAt        int64                  `bson:"enroll_at"`
	LifetimeSeconds int64                  `bson:"lifetime_seconds"`
	Status          int8                   `bson:"status"`
	Config          config.Config          `bson:"Config"`
}
