package host

import (
	"github.com/am6737/headnexus/app/host/command"
	"github.com/am6737/headnexus/app/host/query"
)

type CreateHost struct {
	Name            string
	NetworkID       string
	Role            string
	StaticAddresses map[string]string
	Rules           []string
	IPAddress       string
	Port            int
	IsLighthouse    bool
	Tags            map[string]interface{}
	UserID          string
	PublicIP        string
}

type AddHostRule struct {
	UserID string
	HostID string
	Rules  []string
}

type GetHost struct {
	ID     string
	UserID string
}

type ListHostRules struct {
	UserID   string
	HostID   string
	PageSize int
	PageNum  int
}

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	AddHostRule        command.AddHostRuleHandler
	CreateHost         command.CreateHostHandler
	DeleteHost         command.DeleteHostHandler
	DeleteHostRule     command.DeleteHostRuleHandler
	GenerateEnrollCode command.GenerateEnrollCodeHandler
	EnrollCodeCheck    command.EnrollCodeCheckHandler
	CreateEnrollHost   command.CreateEnrollHostHandler
	//EnrollHost         command.EnrollHostHandler
	UpdateHost command.UpdateHostHandler
}

type Queries struct {
	GetHostConfig query.GetHostConfigHandler
	ListHostRules query.ListHostRulesHandler
	GetHost       query.GetHostHandler
	FindHost      query.FindHostHandler
}
