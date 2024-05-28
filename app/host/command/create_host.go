package command

import (
	"context"
	"errors"
	"fmt"
	v1 "github.com/am6737/headnexus/api/http/v1"
	"github.com/am6737/headnexus/config"
	"github.com/am6737/headnexus/domain/host/entity"
	"github.com/am6737/headnexus/infra/persistence"
	"github.com/am6737/headnexus/pkg/decorator"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
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

type CreateHostHandler decorator.CommandHandler[*CreateHost, *entity.Host]

func NewCreateHostHandler(
	logger *logrus.Logger,
	repos persistence.Repositories,
) CreateHostHandler {
	return &createHostHandler{
		logger: logger,
		repos:  repos,
	}
}

type createHostHandler struct {
	logger *logrus.Logger
	repos  persistence.Repositories
}

func (h *createHostHandler) Handle(ctx context.Context, cmd *CreateHost) (*entity.Host, error) {
	h.logger.WithField("cmd", cmd).Info("Create Request")

	if cmd.IsLighthouse && cmd.PublicIP == "" {
		return nil, errors.New("lighthouse主机必须有公网ip地址")
	}

	find, err := h.repos.HostRepo.Find(ctx, &entity.HostFindOptions{
		UserID: cmd.UserID,
		Name:   cmd.Name,
	})
	if err != nil {
		return nil, err
	}

	if len(find) > 0 {
		return nil, errors.New("host name already exists")
	}

	if cmd.IPAddress != "" {
		if err := h.repos.NetworkRepo.AllocateIP(ctx, cmd.NetworkID, cmd.IPAddress); err != nil {
			h.logger.WithError(err).Error("allocate static address failed")
			return nil, err
		}
	} else {
		ip, err := h.repos.NetworkRepo.AutoAllocateIP(ctx, cmd.NetworkID)
		if err != nil {
			h.logger.WithError(err).Error("allocate auto address failed")
			return nil, err
		}
		cmd.IPAddress = ip
		h.logger.WithField("ip", ip).
			WithField("hostname", cmd.Name).
			WithField("network_id", cmd.NetworkID).
			Info("allocate auto address")
	}

	id := ksuid.New().String()

	net, err := h.repos.NetworkRepo.Get(ctx, cmd.NetworkID)
	if err != nil {
		h.logger.WithError(err).Error("get network failed")
		return nil, err
	}

	var hc config.HostConfig
	if cmd.Role == string(v1.HostRoleLighthouse) {
		hc = config.GenerateLighthouseConfigTemplate()
	} else {
		hc = config.GenerateConfigTemplate()

	}

	if len(cmd.Rules) > 0 {
		rules, err := h.repos.RuleRepo.Gets(ctx, cmd.UserID, cmd.Rules)
		if err != nil {
			h.logger.WithError(err).Error("get rules failed")
			return nil, err
		}
		for _, rule := range rules {
			if rule.Type == entity.RuleTypeInbound {
				hc.Inbound = append(hc.Inbound, config.InboundRule{
					Port:   rule.Port,
					Proto:  rule.Proto.String(),
					Host:   rule.Host,
					Action: rule.Action.String(),
				})
			}
		}
	}

	hc.Listen.Port = cmd.Port
	hc.Tun.IP = cmd.IPAddress
	hc.Tun.Mask = net.Mask()

	yamlData, err := yaml.Marshal(&hc)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	fmt.Println("----------------")
	fmt.Println(string(yamlData))
	fmt.Println("----------------")

	e := &entity.Host{
		Owner:     cmd.UserID,
		ID:        id,
		Name:      cmd.Name,
		NetworkID: cmd.NetworkID,
		Role:      cmd.Role,
		IPAddress: cmd.IPAddress,
		PublicIP:  cmd.PublicIP,
		Port:      cmd.Port,
		Tags:      cmd.Tags,
		Config:    hc,
	}

	h.logger.WithField("entity", e).Info("create host")

	r, err := h.repos.HostRepo.Create(ctx, e)
	if err != nil {
		if err := h.repos.NetworkRepo.ReleaseIP(ctx, cmd.NetworkID, cmd.IPAddress); err != nil {
			h.logger.WithError(err).Error("recycle address failed")
		}
		return nil, err
	}

	if len(cmd.Rules) > 0 {
		if err := h.repos.HostRuleRepo.AddHostRule(ctx, r.ID, cmd.Rules...); err != nil {
			h.logger.WithError(err).Error("add host rule failed")
			return nil, err
		}
	}

	return r, nil
}
