package ports

import (
	"github.com/am6737/headnexus/api/http/v1"
	"github.com/am6737/headnexus/app/network"
	ctime "github.com/am6737/headnexus/common/time"
	"github.com/gin-gonic/gin"
)

func (h *HttpHandler) GetAllNetwork(c *gin.Context) {
	networks, err := h.app.Network.Queries.Handler.Find(c, &network.FindNetwork{})
	if err != nil {
		FailedResponse(c, err.Error())
		return
	}

	var resp []*v1.Network
	for _, v := range networks {
		resp = append(resp, &v1.Network{
			Cidr:      v.Cidr,
			Id:        v.ID,
			Name:      v.Name,
			CreatedAt: ctime.FormatTimestamp(v.CreatedAt),
		})
	}

	SuccessResponse(c, "查询网络成功", resp)
}

func (h *HttpHandler) CreateNetwork(c *gin.Context) {
	req := &v1.CreateNetworkRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		FailedResponse(c, err.Error())
		return
	}

	createNetwork, err := h.app.Network.Commands.Handler.Create(c, &network.CreateNetwork{
		Name: req.Name,
		Cidr: req.Cidr,
	})
	if err != nil {
		FailedResponse(c, err.Error())
		return
	}

	SuccessResponse(c, "创建网络成功", &v1.CreateNetworkResponse{
		Cidr: createNetwork.Cidr,
		Id:   createNetwork.ID,
		Name: createNetwork.Name,
		//CreatedAt: createNetwork.CreatedAt,
		CreatedAt: ctime.FormatTimestamp(createNetwork.CreatedAt),
	})
}

func (h *HttpHandler) DeleteNetwork(c *gin.Context, networkId string) {
	err := h.app.Network.Commands.Handler.Delete(c, &network.DeleteNetwork{
		ID: networkId,
	})
	if err != nil {
		FailedResponse(c, err.Error())
		return
	}

	SuccessResponse(c, "删除网络成功", nil)
}

func (h *HttpHandler) GetNetwork(c *gin.Context, networkId string) {
	network, err := h.app.Network.Queries.Handler.Get(c, &network.GetNetwork{
		ID: networkId,
	})
	if err != nil {
		return
	}
	if err != nil {
		FailedResponse(c, err.Error())
		return
	}

	SuccessResponse(c, "查询网络成功", &v1.Network{
		Id:        network.ID,
		Name:      network.Name,
		Cidr:      network.Cidr,
		CreatedAt: ctime.FormatTimestamp(network.CreatedAt),
	})
}

func (h *HttpHandler) UpdateNetwork(c *gin.Context, networkId string) {
	req := &v1.Network{}
	if err := c.ShouldBindJSON(req); err != nil {
		FailedResponse(c, err.Error())
		return
	}

	// TODO 只允许修改网络名称
	updateNetwork, err := h.app.Network.Commands.Handler.Update(c, &network.UpdateNetwork{
		ID:   networkId,
		Name: req.Name,
	})
	if err != nil {
		FailedResponse(c, err.Error())
		return
	}

	SuccessResponse(c, "更新网络成功", &v1.Network{
		Cidr: updateNetwork.Cidr,
		Id:   updateNetwork.ID,
		Name: updateNetwork.Name,
	})
}
