package inner_api

import (
	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/pkg/util/copyer"
	"github.com/star-table/usercenter/service/api"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/req/inner_req"
	"github.com/star-table/usercenter/service/service"
	"github.com/star-table/usercenter/service/service/inner_service"
)

type userInner int

var UserInner userInner

// @Summary 获取成员信息列表（只获取有效的）
// @Description 根据成员ID列表获取成员信息列表
// @Tags 成员（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.UserListByIdsInnerReq true "入参"
// @Success 200 {object} []inner_resp.UserInfoInnerResp
// @Failure 400
// @Router /usercenter/inner/api/v1/user/getListByIds [post]
func (userInner) GetUserListByIds(c *gin.Context) {

	var reqParam inner_req.UserListByIdsInnerReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	respVO, err := inner_service.GetUserListByIds(reqParam.OrgId, reqParam.Ids, consts.AppCheckStatusSuccess, consts.AppStatusEnable, consts.AppIsNoDelete)
	if err != nil {
		api.Fail(c, err)
		return
	}
	api.Suc(c, respVO)

}

// @Summary 获取成员列表
// @Description 根据获取成员列表
// @Tags 成员（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.UserListReq true "入参"
// @Success 200 {object} inner_resp.UserListResp
// @Failure 400
// @Router /usercenter/inner/api/v1/user/getUserList [post]
func (userInner) GetOrgMemberList(c *gin.Context) {
	var reqParam inner_req.UserListReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}

	// 内部调用就先不进行权限检查
	//permission, err := service.GetOrgUserPerContext(reqParam.OrgId, reqParam.CurUserId)
	//if err != nil {
	//	api.Fail(c, err)
	//	return
	//}
	queryParam := &req.UserListReq{}
	copyer.Copy(reqParam, queryParam)
	res, err := service.GetOrgMemberList(reqParam.OrgId, *queryParam, nil)
	if err != nil {
		api.Fail(c, err)
		return
	}
	//resp := inner_resp.UserListResp{}
	//copyer.Copy(res, &resp)

	api.Suc(c, res)
}

// @Summary 获取所有成员信息列表
// @Description 根据成员ID列表获取成员信息列表
// @Tags 成员（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.UserListByIdsInnerReq true "入参"
// @Success 200 {object} []inner_resp.UserInfoInnerResp
// @Failure 400
// @Router /usercenter/inner/api/v1/user/getAllListByIds [post]
func (userInner) GetAllUserListByIds(c *gin.Context) {
	var reqParam inner_req.UserListByIdsInnerReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.Fail(c, err)
		return
	}
	respVO, err := inner_service.GetUserListByIds(reqParam.OrgId, reqParam.Ids, 0, 0, 0)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, respVO)
}

// @Summary 获取成员权限信息
// @Description 根据成员ID获取成员权限信息
// @Tags 成员（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.UserAuthorityInnerReq true "入参"
// @Success 200 {object} []inner_resp.UserAuthorityInnerResp
// @Failure 400
// @Router /usercenter/inner/api/v1/user/getUserAuthority [post]
func (userInner) GetUserAuthority(c *gin.Context) {

	var reqParam inner_req.UserAuthorityInnerReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	respVO, err := inner_service.GetUserAuthorityByUserId(reqParam.OrgId, reqParam.UserId)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, respVO)

}

// GetUserAuthoritySimple 获取成员权限信息（简化版）
// @Summary  获取成员权限信息（简化版），主要用于少量的授权信息查询。
// @Description 获取成员权限信息（简化版），主要用于少量的授权信息查询。
// @Tags 成员（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.UserAuthorityInnerReq true "入参"
// @Success 200 {object} []inner_resp.UserAuthorityInnerResp
// @Failure 400
// @Router /usercenter/inner/api/v1/user/getUserAuthoritySimple [post]
func (userInner) GetUserAuthoritySimple(c *gin.Context) {
	var reqParam inner_req.UserAuthorityInnerReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	respVO, err := inner_service.GetUserAuthorityByUserIdSimple(reqParam.OrgId, reqParam.UserId)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	api.InnerSuc(c, respVO)

}

// @Summary 获取成员简单信息
// @Description 获取成员简单信息
// @Tags 成员（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.GetMemberSimpleInfoReq true "入参"
// @Success 200 {object} []inner_resp.GetMemberSimpleInfoResp
// @Failure 400
// @Router /usercenter/inner/api/v1/user/getMemberSimpleInfo [post]
func (userInner) GetMemberSimpleInfo(c *gin.Context) {
	var reqParam inner_req.GetMemberSimpleInfoReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}

	respVo, respErr := inner_service.GetMemberSimpleInfo(reqParam.OrgId, reqParam.Type, reqParam.NeedDelete == 1)
	if respErr != nil {
		api.InnerFail(c, err)
		return
	}

	api.InnerSuc(c, respVo)
}

// @Security token
// @Summary 成员列表
// @Description 成员列表接口
// @Tags 成员
// @accept application/json
// @Produce application/json
// @param input body req.MemberSimpleInfoListReq true "入参"
// @Success 200 {object} resp.UserListResp
// @Failure 400
// @Router /usercenter/inner/api/v1/user/memberSimpleInfoList [post]
func (userInner) GetMemberSimpleInfoList(c *gin.Context) {
	var reqParam inner_req.MemberSimpleInfoListReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}

	respVo, err := inner_service.GetMemberSimpleInfoList(reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}

	api.InnerSuc(c, respVo)
}

// @Summary 获取重名成员/部门/角色
// @Description 获取重名成员/部门/角色
// @Tags 成员（内部调用）
// @accept application/json
// @Produce application/json
// @param input body inner_req.SimpleReq true "入参"
// @Success 200 {object} []inner_resp.RepeatMemberResp
// @Failure 400
// @Router /usercenter/inner/api/v1/user/getRepeatMember [post]
func (userInner) GetRepeatMember(c *gin.Context) {
	var reqParam inner_req.SimpleReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}

	respVo, respErr := inner_service.GetRepeatMember(reqParam.OrgId)
	if respErr != nil {
		api.InnerFail(c, err)
		return
	}

	api.InnerSuc(c, respVo)
}

// GetUsersCouldManage 获取有权限管理项目的用户列表
// @Security token
// @Summary 获取有权限管理项目的用户列表
// @Description 成员管理组权限
// @Tags 成员（内部调用）
// @accept application/json
// @Produce application/json
// @Success 200 {object} inner_resp.GetUsersCouldManageResp
// @Failure 400
// @Router /usercenter/inner/api/v1/user/get-users-could-manage [get]
func (userInner) GetUsersCouldManage(c *gin.Context) {
	var reqParam inner_req.GetManageUserReq
	err := api.ParseBody(c, &reqParam)
	if err != nil {
		api.InnerFail(c, err)
		return
	}
	resp, resErr := inner_service.GetUsersCouldManage(reqParam.OrgId, reqParam.AppId)
	if resErr != nil {
		api.InnerFail(c, resErr)
		return
	}

	api.InnerSuc(c, resp)
}
