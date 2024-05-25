package po

type HostRuleRelation struct {
	ID        string `bson:"_id"`
	HostID    string `bson:"host_id"`
	RuleID    string `bson:"rule_id"`
	CreatedAt int64  `bson:"created_at"`
	UpdatedAt int64  `bson:"updated_at"`
	DeletedAt int64  `bson:"deleted_at"`
}
