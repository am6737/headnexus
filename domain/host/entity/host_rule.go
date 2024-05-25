package entity

type HostRuleRelation struct {
	ID        string `bson:"_id"`
	HostID    string `bson:"host_id"`
	RuleID    string `bson:"rule_id"`
	CreatedAt int64  `bson:"created_at"`
}
