package ports

import (
	v1 "github.com/am6737/headnexus/api/http/v1"
	"github.com/am6737/headnexus/app/rule"
	"github.com/am6737/headnexus/domain/rule/entity"
	"github.com/gin-gonic/gin"
)

func (h *HttpHandler) DeleteRule(c *gin.Context, id string) {
	err := h.app.Rule.Commands.Handler.Delete(c, &rule.DeleteRule{
		//UserID: "",
		ID: id,
	})
	if err != nil {
		FailedResponse(c, err.Error())
		return
	}

	SuccessResponse(c, "删除规则成功", nil)
}

func (h *HttpHandler) ListRule(c *gin.Context, params v1.ListRuleParams) {
	if params.RuleFindOptions == nil {
		params.RuleFindOptions = &v1.RuleFindOptions{}
	}

	listRule, err := h.app.Rule.Queries.Handler.Find(c, &rule.FindRule{
		//Name:     params.RuleFindOptions.HostId,
		//Name:     params.RuleFindOptions.HostId,
		PageSize: params.RuleFindOptions.PageSize,
		PageNum:  params.RuleFindOptions.PageNum,
	})
	if err != nil {
		FailedResponse(c, err.Error())
		return
	}

	var listRuleResponse []*v1.Rule
	for _, v := range listRule {
		listRuleResponse = append(listRuleResponse, convertRuleToResponse(v))
	}

	SuccessResponse(c, "查询成功", listRuleResponse)
}

func (h *HttpHandler) CreateRule(c *gin.Context) {
	req := &v1.Rule{}
	if err := c.ShouldBindJSON(req); err != nil {
		FailedResponse(c, "参数错误")
		return
	}

	createRule, err := h.app.Rule.Commands.Handler.Create(c, &rule.CreateRule{
		UserID:      "",
		Type:        string(req.Type),
		Action:      string(req.Action),
		Name:        req.Name,
		Description: req.Description,
		Port:        req.Port,
		Proto:       string(req.Proto),
		Host:        req.Host,
	})
	if err != nil {
		FailedResponse(c, err.Error())
		return
	}

	SuccessResponse(c, "创建成功", convertRuleToResponse(createRule))
}

func convertRuleToResponse(e *entity.Rule) *v1.Rule {
	return &v1.Rule{
		Id:          e.ID,
		Type:        v1.RuleType(e.Type),
		Direction:   v1.RuleDirection(e.Type.String()),
		Action:      v1.RuleAction(e.Action.String()),
		Name:        e.Name,
		Description: e.Description,
		Port:        e.Port,
		Proto:       v1.RuleProto(e.Proto),
		Host:        e.Host,
	}
}
