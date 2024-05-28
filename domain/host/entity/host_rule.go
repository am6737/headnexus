package entity

type HostRuleRelation struct {
	ID        string `bson:"_id"`
	HostID    string `bson:"host_id"`
	RuleID    string `bson:"rule_id"`
	CreatedAt int64  `bson:"created_at"`
}

type HostRule struct {
	Type        RuleType
	CreatedAt   string
	ID          string
	HostID      string
	UserID      string
	Name        string
	Description string
	Port        string
	Proto       RuleProto
	Action      RuleAction
	Host        []string
}
