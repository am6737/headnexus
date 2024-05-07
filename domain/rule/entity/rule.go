package entity

import "fmt"

type Rule struct {
	Type        RuleType
	CreatedAt   int64
	ID          string
	Name        string
	Description string
	HostID      string
	Port        string
	Proto       RuleProto
	Action      RuleAction
	Host        []string
}

type RuleType uint8

const (
	RuleTypeInbound RuleType = iota
	RuleTypeOutbound
)

func ParseRuleType(s string) (RuleType, error) {
	switch s {
	case "inbound":
		return RuleTypeInbound, nil
	case "outbound":
		return RuleTypeOutbound, nil
	default:
		return RuleTypeInbound, fmt.Errorf("unsupported RuleAction string: %s", s)
	}
}

func (rt RuleType) String() string {
	switch rt {
	case RuleTypeInbound:
		return "inbound"
	case RuleTypeOutbound:
		return "outbound"
	default:
		return "unknown"
	}
}

type RuleProto string

const (
	RuleProtoAny  RuleProto = "any"
	RuleProtoTCP  RuleProto = "tcp"
	RuleProtoUDP  RuleProto = "udp"
	RuleProtoICMP RuleProto = "icmp"
)

type RuleAction uint8

func (ra RuleAction) String() string {
	switch ra {
	case RuleActionAllow:
		return "allow"
	case RuleActionDeny:
		return "deny"
	default:
		return "unknown"
	}
}

const (
	RuleActionAllow RuleAction = iota
	RuleActionDeny
)

// ConvertToRuleType 将 uint8 转换为 RuleType
func ConvertToRuleType(value uint8) RuleType {
	switch value {
	case uint8(RuleTypeInbound):
		return RuleTypeInbound
	case uint8(RuleTypeOutbound):
		return RuleTypeOutbound
	default:
		return RuleTypeInbound
	}
}

// ConvertToRuleProto 将 string 转换为 RuleProto
func ConvertToRuleProto(value string) RuleProto {
	switch value {
	case string(RuleProtoAny):
		return RuleProtoAny
	case string(RuleProtoTCP):
		return RuleProtoTCP
	case string(RuleProtoUDP):
		return RuleProtoUDP
	case string(RuleProtoICMP):
		return RuleProtoICMP
	default:
		return RuleProtoAny
	}
}

// ConvertToRuleAction 将 uint8 转换为 RuleAction
func ConvertToRuleAction(value uint8) RuleAction {
	switch value {
	case uint8(RuleActionAllow):
		return RuleActionAllow
	case uint8(RuleActionDeny):
		return RuleActionDeny
	default:
		return RuleActionAllow
	}
}

// ParseRuleAction 将字符串转换为 RuleAction
func ParseRuleAction(s string) (RuleAction, error) {
	switch s {
	case "allow":
		return RuleActionAllow, nil
	case "deny":
		return RuleActionDeny, nil
	default:
		return RuleActionAllow, fmt.Errorf("unsupported RuleAction string: %s", s)
	}
}

type UpdateRule struct {
	Type        uint8
	ID          int64
	Description string
	Port        string
	Proto       string
	Action      string
	Host        []string
}

type FindOptions struct {
	Name     string
	HostID   string
	PageSize int
	PageNum  int
}
