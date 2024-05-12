package ports

import (
	"fmt"
	v1 "github.com/am6737/headnexus/api/http/v1"
	"github.com/am6737/headnexus/app/host"
	"github.com/am6737/headnexus/pkg/http"
	"github.com/gin-gonic/gin"
)

func (h *HttpHandler) CreateEnroll(c *gin.Context, hostId string) {
	req := &v1.CreateEnrollJSONRequestBody{}
	if err := c.ShouldBindJSON(req); err != nil {
		http.FailedResponse(c, "参数错误")
		return
	}

	enrollHost, err := h.app.Host.Commands.Handler.EnrollHost(c, &host.EnrollHost{
		HostID:     hostId,
		EnrollCode: req.Code,
		//UserID:     "",
	})
	if err != nil {
		http.FailedResponse(c, err.Error())
		return
	}

	http.SuccessResponse(c, "success", enrollHost)
}

func (h *HttpHandler) CheckEnrollCode(c *gin.Context, hostId string) {
	req := &v1.CheckEnrollCodeJSONRequestBody{}
	if err := c.ShouldBindJSON(req); err != nil {
		http.FailedResponse(c, "参数错误")
		return
	}

	check, err := h.app.Host.Commands.Handler.EnrollCodeCheck(c, &host.GenerateEnrollCodeCheck{
		Code:   req.Code,
		HostID: hostId,
		//UserID: "",
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
	code, err := h.app.Host.Commands.Handler.GenerateEnrollCode(c, &host.GenerateEnrollCode{
		HostID: hostId,
		//UserID: c.GetString("user_id"),
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
	//fmt.Println("params => ", params.FindOptions.IsLighthouse)

	hosts, err := h.app.Host.Queries.Handler.Find(c.Request.Context(), &host.FindHost{
		//Filters:      params.FindOptions.Filters,
		//Sort:         params.FindOptions.Sort,
		Limit:        params.FindOptions.Limit,
		Offset:       params.FindOptions.Offset,
		NetworkID:    params.FindOptions.NetworkId,
		IPAddress:    params.FindOptions.IpAddress,
		Role:         params.FindOptions.Role,
		Name:         params.FindOptions.Name,
		IsLighthouse: params.FindOptions.IsLighthouse,
	})
	if err != nil {
		http.FailedResponse(c, err.Error())
		return
	}

	var resp = make([]*v1.Host, 0)
	for _, host := range hosts {
		resp = append(resp, convertApiNetworkToHost(host))
	}

	http.SuccessResponse(c, "获取所有主机列表", resp)
}

func (h *HttpHandler) CreateHost(c *gin.Context) {
	req := &v1.Host{}
	if err := c.ShouldBindJSON(req); err != nil {
		fmt.Printf("create host error: %v", err)
		http.FailedResponse(c, "参数错误")
		return
	}

	host, err := h.app.Host.Commands.Handler.Create(c, &host.CreateHost{
		Name:            req.Name,
		NetworkID:       req.NetworkId,
		IPAddress:       req.IpAddress,
		Role:            string(req.Role),
		Port:            req.Port,
		IsLighthouse:    req.IsLighthouse,
		StaticAddresses: req.StaticAddresses,
		Tags:            req.Tags,
	})
	if err != nil {
		http.FailedResponse(c, err.Error())
		return
	}

	http.SuccessResponse(c, "创建主机成功", convertApiNetworkToHost(host))
}

func (h *HttpHandler) DeleteHost(c *gin.Context, hostId string) {
	//err := h.hc.Delete(c, hostId)
	err := h.app.Host.Commands.Handler.Delete(c, &host.DeleteHost{
		UserID: "",
		ID:     hostId,
	})
	if err != nil {
		http.FailedResponse(c, err.Error())
		return
	}
	http.SuccessResponse(c, "删除主机成功", nil)
}

func (h *HttpHandler) GetHost(c *gin.Context, hostId string) {
	//host, err := h.hc.Get(c, hostId)
	host, err := h.app.Host.Queries.Handler.Get(c, &host.GetHost{
		ID:     hostId,
		UserID: "",
	})
	if err != nil {
		http.FailedResponse(c, err.Error())
		return
	}
	http.SuccessResponse(c, "获取主机信息成功", convertApiNetworkToHost(host))
}

func convertApiNetworkToHost(h *host.Host) *v1.Host {
	return &v1.Host{
		Id:              h.ID,
		IpAddress:       h.IPAddress,
		IsLighthouse:    h.IsLighthouse,
		Online:          h.Online,
		Name:            h.Name,
		NetworkId:       h.NetworkID,
		Port:            h.Port,
		Role:            v1.HostRole(h.Role),
		StaticAddresses: h.StaticAddresses,
		Tags:            h.Tags,
		CreatedAt:       h.CreatedAt,
		LastSeenAt:      h.LastSeenAt,
	}
}

func (h *HttpHandler) UpdateHost(c *gin.Context, hostId string) {
	req := &v1.Host{}
	if err := c.ShouldBindJSON(req); err != nil {
		http.FailedResponse(c, "参数错误")
		return
	}

	updatedHost, err := h.app.Host.Commands.Handler.Update(c, &host.UpdateHost{
		ID:              hostId,
		Name:            req.Name,
		NetworkID:       req.NetworkId,
		IPAddress:       req.IpAddress,
		Role:            string(req.Role),
		Port:            req.Port,
		IsLighthouse:    req.IsLighthouse,
		StaticAddresses: req.StaticAddresses,
		Tags:            req.Tags,
	})
	if err != nil {
		return
	}
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
