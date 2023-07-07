package server

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/star-table/usercenter/core/errs"
	_ "github.com/star-table/usercenter/docs"
	"github.com/star-table/usercenter/pkg/util/jsonx"
	"github.com/star-table/usercenter/service/api"
	"github.com/star-table/usercenter/service/api/inner_api"
	"github.com/star-table/usercenter/service/api/open_api"
	"github.com/star-table/usercenter/service/model/vo"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/zeta-io/ginx"
	"github.com/zeta-io/zeta"
)

type Server struct {
	engine *gin.Engine
}

func New(engine *gin.Engine) Server {
	return Server{
		engine: engine,
	}
}

// Init service
func (s *Server) Init() {
	s.initRoutes()
	s.initInnerRoutes()
	s.initOpenRoutes()
}

func (s *Server) initRoutes() {
	router := s.engine
	// 无需路由前缀、无需登录的路由
	router.GET("/usercenter/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 无需登录的业务路由
	noLoginGroup := router.Group("/usercenter/api/v1")
	{
		// 注册
		noLoginGroup.POST("/user/register", api.Account.Register)
		// 找回密码
		noLoginGroup.POST("/user/retrievePassword", api.Account.RetrievePassword)
		// 修改密码
		noLoginGroup.POST("/user/update-password-by-username", api.Account.UpdatePasswordByUsername)

		// 登录
		noLoginGroup.POST("/user/login", api.Authenticate.Login)
		//退出登录
		noLoginGroup.POST("/user/quit", api.Authenticate.Logout)

		// 根据邀请码获取邀请信息
		noLoginGroup.POST("/user/getInviteInfo", api.User.GetInviteInfo)
		// 发送手机邮箱验证码
		noLoginGroup.POST("/user/sendAuthCode", api.SendAuthCode.SendAuthCode)
		// 行业列表
		noLoginGroup.GET("/user/industryList", api.Industry.IndustryList)
	}

	// **需要**登录状态的接口
	needLoginGroup := router.Group("/usercenter/api/v1", api.AuthHandler)
	{

		// 验证原邮箱/手机号
		needLoginGroup.POST("/user/verifyOldName", api.Account.VerifyOldName)
		// 绑定新的联系方式
		needLoginGroup.POST("/user/bindLoginName", api.Account.BindLoginName)
		// 解绑登录名
		needLoginGroup.POST("/user/unbindLoginName", api.Account.UnbindLoginName)
		// 设置本人密码
		needLoginGroup.POST("/user/setPassword", api.Account.SetPassword)
		// 修改本人密码
		needLoginGroup.POST("/user/update-password", api.Account.UpdatePassword)
		// 重置密码，超管使用
		needLoginGroup.POST("/user/reset-password", api.Account.ResetPassword)
		// 校验短信验证码是否正确，并返回一个 token
		needLoginGroup.POST("/user/auth-sms-code", api.SendAuthCode.AuthSmsCode) // AuthSmsCode

		// 获取当前成员信息
		needLoginGroup.POST("/user/personalInfo", api.User.PersonalInfo)
		// 更新当前成员信息
		needLoginGroup.POST("/user/updatePersonalInfo", api.User.UpdatePersonalInfo)
		// 组织切换
		needLoginGroup.POST("/user/switchUserOrganization", api.Org.SwitchUserOrganization)

		//创建组织
		needLoginGroup.POST("/user/createOrg", api.Org.CreateOrg)

		// 新建成员
		needLoginGroup.POST("/user/create", api.User.CreateOrgMember)
		// 更新成员
		needLoginGroup.POST("/user/update", api.User.UpdateOrgMemberInfo)
		// 更新成员状态（启用禁用）
		needLoginGroup.POST("/user/updateUserStatus", api.User.ChangeOrgMemberStatus)
		// 移除成员
		needLoginGroup.POST("/user/removeUser", api.User.RemoveOrgMember)
		// 成员列表
		needLoginGroup.POST("/user/userList", api.User.GetOrgMemberList)
		// 成员信息
		needLoginGroup.POST("/user/member-info", api.User.GetOrgMemberInfoById)
		// 邀请时根据邮箱搜索成员是否存在
		needLoginGroup.POST("/user/searchUser", api.User.SearchUser)
		// 获取邀请码
		needLoginGroup.GET("/user/getInviteCode", api.User.GetInviteCode)
		// 邀请成员
		needLoginGroup.POST("/user/inviteUser", api.User.InviteUser)
		// 未接受邀请的成员
		needLoginGroup.POST("/user/inviteUserList", api.User.InviteUserList)
		// 导出通讯录
		needLoginGroup.POST("/user/exportAddressList", api.User.ExportAddressList)
		// 移除未邀请成员
		needLoginGroup.POST("/user/removeInviteUser", api.User.RemoveInviteUser)
		// 成员tab统计
		needLoginGroup.POST("/user/userStat", api.User.OegMemberStat)

		// 成员管理组权限
		needLoginGroup.GET("/user/get-user-manage-auth", api.User.GetUserManageAuth)

		// 用户组织列表
		needLoginGroup.POST("/user/userOrgList", api.User.UserOrgList)

		// 部门列表
		needLoginGroup.POST("/user/departmentList", api.Department.DepartmentList)
		// 新建部门
		needLoginGroup.POST("/user/createDepartment", api.Department.CreateDepartment)
		// 更新部门
		needLoginGroup.POST("/user/updateDepartment", api.Department.UpdateDepartment)
		// 删除部门
		needLoginGroup.POST("/user/deleteDepartment", api.Department.DeleteDepartment)
		// 交换部门顺序
		needLoginGroup.POST("/dept/change-dept-sort", api.Department.ChangeDeptSort)
		// 将多个用户移出某个部门
		needLoginGroup.POST("/dept/remove-users", api.Department.DeptRemoveUsers)

		// 修改成员部门/职级
		needLoginGroup.POST("/user/change-dept-and-position", api.Department.ChangeUserDeptAndPosition)
		// 给成员分配部门
		needLoginGroup.POST("/user/allocate-dept", api.Department.AllocateUserDept)
		// 设置成员部门等级
		needLoginGroup.POST("/user/setUserDepartmentLevel", api.Department.UpdateDeptLeader)
		// 切换成员的管理组
		needLoginGroup.POST("/user/change-user-admin-group", api.Department.ChangeUserAdminGroup)

		// 角色组新增
		needLoginGroup.POST("/roleGroup/create", api.RoleGroup.Create)
		// 角色组修改
		needLoginGroup.PUT("/roleGroup/update", api.RoleGroup.Update)
		// 角色组删除
		needLoginGroup.DELETE("/roleGroup/delete", api.RoleGroup.Delete)
		// 角色组列表
		needLoginGroup.GET("/roleGroup/list", api.RoleGroup.GetGroupList)

		// 角色新增
		needLoginGroup.POST("/role/create", api.Role.Create)
		// 角色修改
		needLoginGroup.PUT("/role/update", api.Role.Update)
		// 角色删除
		needLoginGroup.DELETE("/role/delete", api.Role.Delete)
		// 角色移动
		needLoginGroup.PUT("/role/move", api.Role.Move)
		// 角色列表查询
		needLoginGroup.GET("/role/list", api.Role.GetOrgRoleList)

		// 成员分配角色
		needLoginGroup.POST("/userRole/assignRoles", api.UserRole.AssignRoles)
		// 角色添加成员
		needLoginGroup.POST("/userRole/create", api.UserRole.Create)
		// 角色删除成员
		needLoginGroup.DELETE("/userRole/delete", api.UserRole.Delete)

		// 管理组新增
		needLoginGroup.POST("/adminGroup/create", api.ManageGroup.CreateGroup)
		// 管理组修改
		needLoginGroup.PUT("/adminGroup/update", api.ManageGroup.UpdateGroup)
		// 管理组删除
		needLoginGroup.DELETE("/adminGroup/delete", api.ManageGroup.DeleteGroup)
		// 管理组修改权限 成员\组织\应用包
		needLoginGroup.PUT("/adminGroup/updateContents", api.ManageGroup.UpdateContents)
		// 管理组树
		needLoginGroup.GET("/adminGroup/tree", api.ManageGroup.GetManageGroupTree)
		// 管理组详情查询
		needLoginGroup.GET("/adminGroup/detail", api.ManageGroup.GetManageGroupDetail)
		// 获取管理组的权限项配置
		needLoginGroup.GET("/adminGroup/operationConfigs", api.ManageGroup.GetManageGroupOperationConfig)

		// oss 存储，头像上传，获取策略信息
		needLoginGroup.POST("/oss/getOssPostPolicy", api.Oss.GetOssPostPolicy)

		// 开启OpenApiKey
		needLoginGroup.POST("/org/generate-api-key", api.Org.GenerateApiKey)
		needLoginGroup.POST("/org/reset-api-key", api.Org.ResetApiKey)
		needLoginGroup.DELETE("/org/close-api-key", api.Org.CloseOpenApi)
		needLoginGroup.POST("/org/switch-to-inner-member", api.Org.SwitchToInnerMember)

		// 职级
		needLoginGroup.POST("/positions/create", api.Position.CreatePosition)
		needLoginGroup.POST("/positions/update-info", api.Position.ModifyPositionInfo)
		needLoginGroup.POST("/positions/update-status", api.Position.UpdatePositionStatus)
		needLoginGroup.DELETE("/positions/delete", api.Position.DeletePosition)
		needLoginGroup.POST("/positions/list", api.Position.GetPositionList)
		needLoginGroup.POST("/positions/page", api.Position.GetPositionPageList)

		// 日志
		needLoginGroup.POST("/records/login-records", api.Records.GetLoginRecordsList)           //登陆日志
		needLoginGroup.POST("/records/export-login-records", api.Records.ExportLoginRecordsList) // 登陆日志导出

		zetaRouter := zeta.Router("/usercenter/api/v1", api.AuthHandler)
		contactGroup := zetaRouter.Group("/contact")
		contactGroup.Post("/filter", api.Contact.Filter)
		contactGroup.Post("/search", api.Contact.Search)
		contactGroup.Post("/aggregation", api.Contact.Aggregation)
		roleGroup := zetaRouter.Group("/role")
		roleGroup.Post("/filter", api.Role.Filter)

		zeta.New(zetaRouter, ginx.New(s.engine).Serial(zeta.DefaultSerial(jsonx.JSON())).Response(func(c *gin.Context, data interface{}, err error) {
			msg := errs.OK.Error()
			code := errs.OK.Code()
			if err != nil {
				if ec, ok := err.(*errs.ResultCodeInfo); ok {
					msg = ec.Error()
					code = ec.Code()
				} else if validatorErrs, ok := err.(validator.ValidationErrors); ok && len(validatorErrs) > 0 {
					msg = errs.ReqParamsValidateError.Message() + ":" + validatorErrs[0].Field()
					code = errs.ReqParamsValidateError.Code()
				} else if ec, ok := err.(error); ok {
					msg = ec.Error()
					code = errs.ServerError.Code()
				}
			}
			resp := vo.Resp{
				Data:    data,
				Message: msg,
				Code:    code,
			}
			c.JSON(200, resp)
			c.Abort()
		})).Complete()
	}
}

// initInnerRoutes 初始化内部调用路由
func (s *Server) initInnerRoutes() {
	router := s.engine

	router.GET("/inner/api/v1/health", func(c *gin.Context) {
		c.Data(200, "text/plain;charset=utf8", []byte("ok"))
	})

	// 内部调用
	innerRouters := router.Group("/usercenter/inner/api/v1")
	//校验token
	innerRouters.POST("/user/auth", api.Authenticate.Auth)
	//校验api-key
	innerRouters.POST("/auth/api-key-auth", api.OpenApi.ApiKeyAuth)

	//校验token并且校验状态
	innerRouters.POST("/user/auth-check-status", api.Authenticate.AuthCheckStatus)

	// 根据各种ID(UserId/DeptId/...)获取成员基本信息列表
	innerRouters.POST("/universal/getUserBaseInfoByIds", inner_api.UniversalInner.GetUserBaseInfoByIds)

	// 根据成员ID获取可用成员列表
	innerRouters.POST("/user/getListByIds", inner_api.UserInner.GetUserListByIds)
	// 根据成员ID获取所有成员列表
	innerRouters.POST("/user/getAllListByIds", inner_api.UserInner.GetAllUserListByIds)
	// 根据成员ID获取成员权限信息
	innerRouters.POST("/user/getUserAuthority", inner_api.UserInner.GetUserAuthority)
	// 根据成员ID获取成员权限信息（简化版），主要用于少量的授权信息查询
	innerRouters.POST("/user/getUserAuthoritySimple", inner_api.UserInner.GetUserAuthoritySimple)
	// 获取成员简单信息（1成员，2部门，3角色）
	innerRouters.POST("/user/getMemberSimpleInfo", inner_api.UserInner.GetMemberSimpleInfo)
	// 简易成员信息列表
	innerRouters.POST("/user/memberSimpleInfoList", inner_api.UserInner.GetMemberSimpleInfoList)
	// 获取重复成员信息（成员/部门/橘色）
	innerRouters.POST("/user/getRepeatMember", inner_api.UserInner.GetRepeatMember)
	// 获取有权限管理项目的用户列表
	innerRouters.POST("/user/get-users-could-manage", inner_api.UserInner.GetUsersCouldManage)
	// search user list
	innerRouters.POST("/user/get-user-list", inner_api.UserInner.GetOrgMemberList)

	// 根据角色ID获取角色列表
	innerRouters.POST("/role/getListByIds", inner_api.RoleInner.GetRoleListByIds)
	innerRouters.POST("/role/getAllListByIds", inner_api.RoleInner.GetAllRoleListByIds)
	// 获取角色用户（角色对应用户数组）
	innerRouters.POST("/role/getUserIds", inner_api.RoleInner.GetRoleUserIds)

	// 管理组
	innerRouters.POST("/manage-group/add-pkg", inner_api.ManageGroupInner.AddAppPkgToManageGroup)
	innerRouters.POST("/manage-group/delete-pkg", inner_api.ManageGroupInner.DeleteAppPkgFromManageGroup)
	innerRouters.POST("/manage-group/add-app", inner_api.ManageGroupInner.AddAppToManageGroup)
	innerRouters.POST("/manage-group/delete-app", inner_api.ManageGroupInner.DeleteAppFromManageGroup)
	innerRouters.POST("/manage-group/getManagerInfo", inner_api.ManageGroupInner.GetManager)
	// 创建组织的时候初始化管理组
	innerRouters.POST("/manage-group/init", inner_api.ManageGroupInner.ManageGroupInit)
	// 增加人员到系统管理组
	innerRouters.POST("/manage-group/add-user-to-sys-manage-group", inner_api.ManageGroupInner.AddUserToSysManageGroup)
	// 更换组织的超管
	innerRouters.POST("/manage-group/replace-super-admin", inner_api.ManageGroupInner.ReplaceSuperAdmin)
	// 获取组织的超管（id 列表）
	innerRouters.GET("/manage-group/get-super-admin-ids", inner_api.ManageGroupInner.GetOrgSuperAdminIds)

	// 获取普通管理员可以管理的app
	innerRouters.POST("/manage-group/getCommonAdminManageApps", inner_api.ManageGroupInner.GetCommAdminMangeApps)

	// 管理组树
	innerRouters.POST("/adminGroup/tree", inner_api.ManageGroupInner.GetManageGroupTree)
	// 管理组修改权限 成员\组织\应用包
	innerRouters.PUT("/adminGroup/updateContents", inner_api.ManageGroupInner.UpdateContents)
	// 系统管理组详情
	innerRouters.GET("/adminGroup/detail", inner_api.ManageGroupInner.GetManageGroupDetail)
	// 将一个用户从该组织的所有管理组中移除
	innerRouters.POST("/adminGroup/deleteOneUserFromOrg", inner_api.ManageGroupInner.DeleteOneUserFromOrg)
	// 向组织角色中追加新的菜单权限项
	innerRouters.POST("/adminGroup/addNewMenuToRole", inner_api.ManageGroupInner.AddNewMenuToRole)

	//统计部门下的人数
	innerRouters.POST("/dept/getUserCountByDeptIds", inner_api.DeptInner.GetUserCountByDeptIds)
	//获取成员所属部门
	innerRouters.POST("/dept/getUserDeptIds", inner_api.DeptInner.GetUserDeptIds)
	//获取成员所属部门(批量)
	innerRouters.POST("/dept/getUserDeptIdsBatch", inner_api.DeptInner.GetUserDeptIdsBatch)
	//获取某些部门下的人员id
	innerRouters.POST("/dept/getUserIdsByDeptIds", inner_api.DeptInner.GetUserIdsByDeptIds)
	//获取某些部门下的leaders
	innerRouters.POST("/dept/getLeadersByDeptIds", inner_api.DeptInner.GetLeadersByDeptIds)
	// 根据部门ID获取可用部门列表
	innerRouters.POST("/dept/getListByIds", inner_api.DeptInner.GetDeptListByIds)
	// 根据部门ID获取所有部门列表
	innerRouters.POST("/dept/getAllListByIds", inner_api.DeptInner.GetAllDeptListByIds)
	// 获取部门用户（部门对应用户数组）
	innerRouters.POST("/dept/getUserIds", inner_api.DeptInner.GetDeptUserIds)
	//获取完整路径的部门信息，例如：职能部/技术部/运维组
	innerRouters.POST("/dept/getFullDeptByIds", inner_api.DeptInner.GetDeptFullNameByIds)
	// 获取部门列表
	innerRouters.POST("/dept/getDeptList", inner_api.DeptInner.GetDeptList)
	// 获取用户对应的部门id（包括父部门id）
	innerRouters.POST("/dept/getUserDeptIdsWithParentId", inner_api.DeptInner.GetUserDeptIdsWithParentId)

	//获取组织信息
	innerRouters.GET("/org/info", inner_api.OrgInner.GetOrgInfo)
	//获取文件上传策略信息
	innerRouters.POST("/oss/getOssPostPolicy", inner_api.OssInner.GetOssPostPolicy)
	//新增组织外部协作人
	innerRouters.POST("/org/add-out-collaborator", inner_api.OrgInner.AddOutCollaborator)
	//检查组织拥有者是否是超管，如果不是，则设置为超管。
	innerRouters.POST("/org/check-and-set-super-admin", inner_api.OrgInner.CheckAndSetSuperAdmin)
}

// initInnerRoutes 初始化OpenApi调用路由
func (s *Server) initOpenRoutes() {
	router := s.engine
	// 内部调用
	openRouters := router.Group("/open/usercenter/api/v1", api.AuthHandler)

	// token认证
	openRouters.POST("/auth/check", open_api.AuthOpen.AuthCheckStatus)

	// 获取部门列表
	openRouters.POST("/dept/list", open_api.DeptOpen.GetDeptList)
	// 获取部门列表
	openRouters.POST("/dept/list-by-query-cond", open_api.DeptOpen.GetDeptListByQueryCond)
	// 根据部门ID获取部门列表
	openRouters.POST("/dept/list-by-ids", open_api.DeptOpen.GetDeptListByIds)
	// 获取子部门
	openRouters.POST("/dept/children", open_api.DeptOpen.GetDeptChildrenList)
	// 获取部门信息列表,包含成员
	openRouters.POST("/dept/list-have-member", open_api.DeptOpen.GetDeptHaveMemberList)

	// 获取组织角色列表
	openRouters.POST("/role/list", open_api.RoleOpen.GetRoleList)
	// 根据角色ID获取组织角色列表
	openRouters.POST("/role/list-by-ids", open_api.RoleOpen.GetRoleListByIds)

	// 获取组织成员列表
	openRouters.POST("/member/list", open_api.MemberOpen.GetOrgMemberList)
	// 人员权限基本信息
	openRouters.POST("/member/user-auth-base-info", open_api.MemberOpen.GetUserAuthBaseInfo)
	// 获取组织成员列表
	openRouters.POST("/member/list-by-query-cond", open_api.MemberOpen.GetOrgMemberListByQueryCond)
	// 获取指定组织成员列表
	openRouters.POST("/member/list-by-ids", open_api.MemberOpen.GetOrgMemberListByUserIds)
	// 获取指定组织成员
	openRouters.POST("/member/by-id", open_api.MemberOpen.GetOrgMemberByUserId)
	// 获取组织成员列表，排除指定成员
	openRouters.POST("/member/exclude-ids", open_api.MemberOpen.GetOrgMemberListByExcludeIds)
	// 根据成员ID获取绑定部门信息列表
	openRouters.POST("/member/by-dept", open_api.MemberOpen.GetOrgMemberListByDept)
	// 根据成员ID获取绑定部门信息列表
	openRouters.POST("/member/dept-list-by-user", open_api.MemberOpen.GetUserDeptBindListByUser)
	// 根据成员ID列表获取绑定部门信息列表
	openRouters.POST("/member/dept-list-by-users", open_api.MemberOpen.GetUserDeptBindListByUsers)
	// 根据成员ID获取绑定角色信息列表
	openRouters.POST("/member/role-list-by-user", open_api.MemberOpen.GetUserRoleBindListByUser)
	// 根据成员ID列表获取绑定角色信息列表
	openRouters.POST("/member/role-list-by-users", open_api.MemberOpen.GetUserRoleBindListByUsers)
	// 根据成员ID列表获取绑定部门角色信息列表
	openRouters.POST("/member/dept-role-list-by-users", open_api.MemberOpen.GetUserBindDeptAndRoleListByUsers)
	// 获取平级或者上级的个数
	openRouters.POST("/member/same-superior-count", open_api.MemberOpen.GetSameOrSuperiorCount)
	// 获取平级或者下级的用户ID
	openRouters.POST("/member/get-same-subordinate-users", open_api.MemberOpen.GetSameOrSubordinateMembers)

}

func (s *Server) Run(addr string) {
	if err := s.engine.Run(addr); err != nil {
		panic(err)
	}
}

// Close close the server.
func (s *Server) Close() {

}
