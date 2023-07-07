package api

import (
	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/resp"
	"github.com/star-table/usercenter/service/service"
)

type contact int

var Contact contact = 1

// @Security token
// @Summary 通讯录查询
// @Description 通讯录查询
// @Tags 选人控件
// @accept application/json
// @Produce application/json
// @param input body req.ContactFilterReq true "请求结构体"
// @Success 200 {object} resp.ContactFilterResp
// @Failure 400
// @Router /usercenter/api/v1/contact/filter [post]
func (contact) Filter(c *gin.Context, args struct {
	Body req.ContactFilterReq `param:"body"`
}) (*resp.ContactFilterResp, error) {
	operator, err := GetOperator(c)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return service.ContactFilter(operator.OrgId, args.Body)
}

// @Security token
// @Summary 通讯录搜索
// @Description 通讯录搜索
// @Tags 选人控件
// @accept application/json
// @Produce application/json
// @param input body req.ContactSearchReq true "请求结构体"
// @Success 200 {object} resp.ContactSearchResp
// @Failure 400
// @Router /usercenter/api/v1/contact/search [post]
func (contact) Search(c *gin.Context, args struct {
	Body req.ContactSearchReq `param:"body"`
}) (*resp.ContactSearchResp, error) {
	operator, err := GetOperator(c)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return service.ContactSearch(operator.OrgId, args.Body)
}

// @Security token
// @Summary 通讯录组合搜索（用户、部门、角色）
// @Description 通讯录搜索
// @Tags 选人控件
// @accept application/json
// @Produce application/json
// @param input body req.AggregationReq true "请求结构体"
// @Success 200 {object} resp.AggregationResp
// @Failure 400
// @Router /usercenter/api/v1/contact/aggregation [post]
func (contact) Aggregation(c *gin.Context, args struct {
	Body req.AggregationReq `param:"body"`
}) (*resp.AggregationResp, error) {
	operator, err := GetOperator(c)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	//用户
	user, userErr := service.ContactSearch(operator.OrgId, req.ContactSearchReq{
		SearchType: 1,
		Query:      args.Body.Query,
		Offset:     0,
		Limit:      0,
		Scope:      req.ContactScope{},
		OnlyMember: false,
	})
	if userErr != nil {
		logger.Error(userErr)
		return nil, userErr
	}

	//部门
	dept, deptErr := service.ContactSearch(operator.OrgId, req.ContactSearchReq{
		SearchType: 2,
		Query:      args.Body.Query,
		Offset:     0,
		Limit:      0,
		Scope:      req.ContactScope{},
		OnlyMember: false,
	})
	if deptErr != nil {
		logger.Error(deptErr)
		return nil, deptErr
	}

	//角色
	role, roleErr := service.ContactRoleFilter(operator.OrgId, req.RoleFilterReq{Query: args.Body.Query})
	if roleErr != nil {
		logger.Error(roleErr)
		return nil, roleErr
	}

	return &resp.AggregationResp{
		DepList:  dept.DepList,
		User:     user.User,
		RoleList: role.RoleList,
	}, nil
}
