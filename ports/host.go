package ports

import (
	v1 "github.com/am6737/headnexus/api/http/v1"
	"github.com/am6737/headnexus/app/host/command"
	"github.com/am6737/headnexus/app/host/query"
	ctime "github.com/am6737/headnexus/common/time"
	"github.com/am6737/headnexus/domain/host/entity"
	"github.com/am6737/headnexus/pkg/http"
	"github.com/gin-gonic/gin"
)

func (h *HttpHandler) CreateEnroll(c *gin.Context, hostId string) {
	//TODO implement me
	panic("implement me")
}

func (h *HttpHandler) EnrollHost(c *gin.Context, code string) {
	enrollHost, err := h.app.Host.Commands.CreateEnrollHost.Handle(c, &command.CreateEnrollHost{
		EnrollCode: code,
		//UserID:     c.GetString("user_id"),
	})
	if err != nil {
		http.FailedResponse(c, err.Error())
		return
	}

	http.SuccessResponse(c, "注册主机成功", &v1.EnrollHostResponse{
		EnrollAt:   enrollHost.EnrollAt,
		HostId:     enrollHost.HostID,
		LastSeenAt: enrollHost.LastSeenAt,
		Online:     enrollHost.Online,
		Config:     enrollHost.Config,
	})
}

func (h *HttpHandler) GetHostConfig(c *gin.Context, hostId string) {
	getHostConfig, err := h.app.Host.Queries.GetHostConfig.Handle(c, &query.GetHostConfig{
		UserID: c.GetString("user_id"),
		HostID: hostId,
	})
	if err != nil {
		http.FailedResponse(c, err.Error())
		return
	}

	http.SuccessResponse(c, "获取主机配置成功", &v1.HostConfig{Config: getHostConfig.Config})
}

func (h *HttpHandler) DeleteHostRule(c *gin.Context, hostId string, ruleId string) {
	err := h.app.Host.Commands.DeleteHostRule.Handle(c, &command.DeleteHostRule{
		UserID: c.Value("user_id").(string),
		HostID: hostId,
		Rules:  []string{ruleId},
	})
	if err != nil {
		http.FailedResponse(c, err.Error())
		return
	}

	http.SuccessResponse(c, "删除主机规则成功", nil)
}

func (h *HttpHandler) ListHostRules(c *gin.Context, hostId string, params v1.ListHostRulesParams) {
	rules, err := h.app.Host.Queries.ListHostRules.Handle(c, &query.ListHostRules{
		UserID:   c.Value("user_id").(string),
		HostID:   hostId,
		PageSize: *params.PageSize,
		PageNum:  *params.PageNum,
	})
	if err != nil {
		http.FailedResponse(c, err.Error())
		return
	}

	http.SuccessResponse(c, "success", listHostRulesToResponse(rules))
}

func listHostRulesToResponse(rules []*entity.HostRule) []*v1.HostRule {
	var res []*v1.HostRule
	for _, rule := range rules {
		res = append(res, &v1.HostRule{
			Action:      v1.HostRuleAction(rule.Action),
			CreatedAt:   rule.CreatedAt,
			Description: rule.Description,
			Host:        rule.Host,
			HostId:      rule.HostID,
			Id:          rule.ID,
			Name:        rule.Name,
			Port:        rule.Port,
			Proto:       v1.HostRuleProto(rule.Proto),
			Type:        v1.HostRuleType(rule.Type),
		})
	}
	return res
}

func hostRuleToResponse(rule *entity.HostRule) *v1.HostRule {
	return &v1.HostRule{
		Action:      v1.HostRuleAction(rule.Action),
		CreatedAt:   rule.CreatedAt,
		Description: rule.Description,
		Host:        rule.Host,
		HostId:      rule.HostID,
		Id:          rule.ID,
		Name:        rule.Name,
		Port:        rule.Port,
		Proto:       v1.HostRuleProto(rule.Proto),
		Type:        v1.HostRuleType(rule.Type),
	}
}

func (h *HttpHandler) AddHostRule(c *gin.Context, hostId string) {
	req := &v1.AddHostRuleJSONRequestBody{}
	if err := c.ShouldBindJSON(req); err != nil {
		http.FailedResponse(c, "参数错误")
		return
	}
	addHostRule, err := h.app.Host.Commands.AddHostRule.Handle(c, &command.AddHostRule{
		UserID: c.GetString("user_id"),
		HostID: hostId,
		Rules:  req.Rules,
	})
	if err != nil {
		http.FailedResponse(c, err.Error())
		return
	}

	http.SuccessResponse(c, "添加主机规则成功", listHostRulesToResponse(addHostRule))
}

func (h *HttpHandler) CheckEnrollCode(c *gin.Context, hostId string) {
	req := &v1.CheckEnrollCodeJSONRequestBody{}
	if err := c.ShouldBindJSON(req); err != nil {
		http.FailedResponse(c, "参数错误")
		return
	}

	check, err := h.app.Host.Commands.EnrollCodeCheck.Handle(c, &command.GenerateEnrollCodeCheck{
		Code:   req.Code,
		HostID: hostId,
		UserID: c.GetString("user_id"),
	})
	if err != nil {
		http.FailedResponse(c, err.Error())
		return
	}

	http.SuccessResponse(c, "success", &v1.CheckEnrollCodeResponse{
		Exists: &check.Exists,
	})
}

func (h *HttpHandler) CreateEnrollCode(c *gin.Context, hostId string) {
	code, err := h.app.Host.Commands.GenerateEnrollCode.Handle(c, &command.GenerateEnrollCode{
		HostID: hostId,
		UserID: c.GetString("user_id"),
	})
	if err != nil {
		http.FailedResponse(c, err.Error())
		return
	}

	http.SuccessResponse(c, "success", &v1.CreateEnrollCodeResponse{
		Code:            code.Code,
		LifetimeSeconds: code.LifetimeSeconds,
	})
}

func (h *HttpHandler) ListHost(c *gin.Context, params v1.ListHostParams) {
	if params.FindOptions == nil {
		http.FailedResponse(c, "参数错误")
		return
	}

	hosts, err := h.app.Host.Queries.FindHost.Handle(c, &query.FindHost{
		UserID: c.GetString("user_id"),
		//Filters:      params.FindOptions.Filters,
		//Sort:         params.FindOptions.Sort,
		Limit:     params.FindOptions.Limit,
		Offset:    params.FindOptions.Offset,
		NetworkID: params.FindOptions.NetworkId,
		IPAddress: params.FindOptions.IpAddress,
		Role:      string(params.FindOptions.Role),
		Name:      params.FindOptions.Name,
	})
	if err != nil {
		http.FailedResponse(c, err.Error())
		return
	}

	var resp = make([]*v1.ShortHost, 0)
	for _, host := range hosts {
		resp = append(resp, convertApiNetworkToShortHost(host))
	}

	http.SuccessResponse(c, "获取主机列表", resp)
}

func convertApiNetworkToShortHost(h *query.Host) *v1.ShortHost {
	return &v1.ShortHost{
		CreatedAt:  h.CreatedAt,
		Id:         h.ID,
		IpAddress:  h.IPAddress,
		LastSeenAt: h.LastSeenAt,
		Name:       h.Name,
		Online:     h.Online,
		PublicIp:   h.PublicIP,
		Port:       h.Port,
	}
}

func (h *HttpHandler) CreateHost(c *gin.Context) {
	req := &v1.CreateHostJSONRequestBody{}
	if err := c.ShouldBindJSON(req); err != nil {
		http.FailedResponse(c, "参数错误")
		return
	}

	host, err := h.app.Host.Commands.CreateHost.Handle(c, &command.CreateHost{
		UserID:       c.GetString("user_id"),
		Name:         req.Name,
		NetworkID:    req.NetworkId,
		IPAddress:    req.IpAddress,
		PublicIP:     req.PublicIp,
		Role:         string(req.Role),
		IsLighthouse: req.Role == v1.Lighthouse,
		Port:         req.Port,
		//StaticAddresses: req.StaticAddresses,
		Rules: req.Rules,
		//Tags:            req.Tags,
	})
	if err != nil {
		http.FailedResponse(c, err.Error())
		return
	}

	http.SuccessResponse(c, "创建主机成功", convertApiNetworkToHost(host))
}

func (h *HttpHandler) DeleteHost(c *gin.Context, hostId string) {
	err := h.app.Host.Commands.DeleteHost.Handle(c, &command.DeleteHost{
		UserID: c.GetString("user_id"),
		ID:     hostId,
	})
	if err != nil {
		http.FailedResponse(c, err.Error())
		return
	}
	http.SuccessResponse(c, "删除主机成功", nil)
}

func (h *HttpHandler) GetHost(c *gin.Context, hostId string) {
	host, err := h.app.Host.Queries.GetHost.Handle(c, &query.GetHost{
		HostID: hostId,
		UserID: c.GetString("user_id"),
	})
	if err != nil {
		http.FailedResponse(c, err.Error())
		return
	}
	http.SuccessResponse(c, "获取主机信息成功", convertApiNetworkToHost(host))
}

func convertApiNetworkToHost(h *entity.Host) *v1.Host {
	return &v1.Host{
		Id:        h.ID,
		IpAddress: h.IPAddress,
		PublicIp:  h.PublicIP,
		Online:    h.Status == entity.Online,
		Name:      h.Name,
		NetworkId: h.NetworkID,
		Port:      h.Port,
		Role:      v1.HostRole(h.Role),
		//StaticAddresses: h.StaticAddresses,
		Tags:       h.Tags,
		CreatedAt:  ctime.FormatTimestamp(h.CreatedAt),
		LastSeenAt: ctime.FormatTimestamp(h.LastSeenAt),
	}
}

func (h *HttpHandler) UpdateHost(c *gin.Context, hostId string) {
	req := &v1.UpdateHostJSONRequestBody{}
	if err := c.ShouldBindJSON(req); err != nil {
		http.FailedResponse(c, "参数错误")
		return
	}

	updatedHost, err := h.app.Host.Commands.UpdateHost.Handle(c, &command.UpdateHost{
		UserID:    c.GetString("user_id"),
		HostID:    hostId,
		Name:      req.Name,
		IPAddress: req.IpAddress,
		PublicIP:  req.PublicIp,
		Port:      req.Port,
	})
	if err != nil {
		http.FailedResponse(c, err.Error())
		return
	}

	http.SuccessResponse(c, "更新主机成功", updatedHost)
}

func convertMapToStringMap(inputMap map[string]interface{}) map[string]string {
	outputMap := make(map[string]string)
	for key, value := range inputMap {
		stringValue, ok := value.(string)
		if !ok {
			stringValue = ""
		}
		outputMap[key] = stringValue
	}
	return outputMap
}
