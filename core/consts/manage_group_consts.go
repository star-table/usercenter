package consts

/**
管理组
*/
const (
	DeptOpt       = int64(1)
	RoleOpt       = int64(2)
	AppPackageOpt = int64(3)
)

// 全部操作权限
var AllManageGroupOptAuths = []int64{
	DeptOpt, RoleOpt, AppPackageOpt,
}

// 特殊管理组名称
var DefaultManageGroupName = []string{
	"超级管理员",
}

// 极星，管理组的权限项
const (
	OperationOrgConfigView               = "Permission.Org.Config.View"
	OperationOrgConfigModify             = "Permission.Org.Config.Modify"
	OperationOrgConfigTransfer           = "Permission.Org.Config.Transfer"
	OperationOrgMessageConfigView        = "Permission.Org.MessageConfig.View"
	OperationOrgMessageConfigModify      = "Permission.Org.MessageConfig.Modify"
	OperationOrgUserView                 = "Permission.Org.User.View"
	OperationOrgUserWatch                = "Permission.Org.User.Watch"
	OperationOrgUserModifyStatus         = "Permission.Org.User.ModifyStatus"
	OperationOrgUserInvite               = "Permission.Org.User.Invite"
	OperationOrgUserModifyUserAdminGroup = "Permission.Org.User.ModifyUserAdminGroup"
	OperationOrgUserModifyUserDept       = "Permission.Org.User.ModifyUserDept"
	OperationOrgTeamView                 = "Permission.Org.Team.View"
	OperationOrgTeamCreate               = "Permission.Org.Team.Create"
	OperationOrgTeamModify               = "Permission.Org.Team.Modify"
	OperationOrgTeamDelete               = "Permission.Org.Team.Delete"
	OperationOrgTeamModifyStatus         = "Permission.Org.Team.ModifyStatus"
	OperationOrgTeamBind                 = "Permission.Org.Team.Bind"
	OperationOrgDeptCreate               = "Permission.Org.Department.Create"
	OperationOrgDeptModify               = "Permission.Org.Department.Modify"
	OperationOrgDeptDelete               = "Permission.Org.Department.Delete"
	OperationOrgAdminGroupView           = "Permission.Org.AdminGroup.View"
	OperationOrgAdminGroupCreate         = "Permission.Org.AdminGroup.Create"
	OperationOrgAdminGroupModify         = "Permission.Org.AdminGroup.Modify"
	OperationOrgAdminGroupDelete         = "Permission.Org.AdminGroup.Delete"
	OperationOrgRoleGroupView            = "Permission.Org.RoleGroup.View"
	OperationOrgRoleGroupCreate          = "Permission.Org.RoleGroup.Create"
	OperationOrgRoleGroupModify          = "Permission.Org.RoleGroup.Modify"
	OperationOrgRoleGroupDelete          = "Permission.Org.RoleGroup.Delete"
	OperationOrgRoleView                 = "Permission.Org.Role.View"
	OperationOrgRoleCreate               = "Permission.Org.Role.Create"
	OperationOrgRoleModify               = "Permission.Org.Role.Modify"
	OperationOrgRoleDelete               = "Permission.Org.Role.Delete"
	OperationOrgRoleBind                 = "Permission.Org.Role.Bind"
	OperationOrgProjectView              = "Permission.Org.Project.View"
	OperationOrgProjectCreate            = "Permission.Org.Project.Create"
	OperationOrgProjectModify            = "Permission.Org.Project.Modify"
	OperationOrgProjectDelete            = "Permission.Org.Project.Delete"
	OperationOrgProjectAttention         = "Permission.Org.Project.Attention"
	OperationOrgProjectFiling            = "Permission.Org.Project.Filing"
	OperationOrgProjectTypeView          = "Permission.Org.ProjectType.View"
	OperationOrgProjectTypeModify        = "Permission.Org.ProjectType.Modify"
	OperationOrgProjectTypeCreate        = "Permission.Org.ProjectType.Create"
	OperationOrgProjectTypeDelete        = "Permission.Org.ProjectType.Delete"
	OperationOrgIssueSourceView          = "Permission.Org.IssueSource.View"
	OperationOrgIssueSourceModify        = "Permission.Org.IssueSource.Modify"
	OperationOrgIssueSourceCreate        = "Permission.Org.IssueSource.Create"
	OperationOrgIssueSourceDelete        = "Permission.Org.IssueSource.Delete"
	OperationOrgProjectObjectTypeView    = "Permission.Org.ProjectObjectType.View"
	OperationOrgProjectObjectTypeModify  = "Permission.Org.ProjectObjectType.Modify"
	OperationOrgProjectObjectTypeCreate  = "Permission.Org.ProjectObjectType.Create"
	OperationOrgProjectObjectTypeDelete  = "Permission.Org.ProjectObjectType.Delete"
	OperationOrgPriorityView             = "Permission.Org.Priority.View"
	OperationOrgPriorityModify           = "Permission.Org.Priority.Modify"
	OperationOrgPriorityCreate           = "Permission.Org.Priority.Create"
	OperationOrgPriorityDelete           = "Permission.Org.Priority.Delete"
	OperationOrgProcessStatusView        = "Permission.Org.ProcessStatus.View"
	OperationOrgProcessStatusModify      = "Permission.Org.ProcessStatus.Modify"
	OperationOrgProcessStatusCreate      = "Permission.Org.ProcessStatus.Create"
	OperationOrgProcessStatusDelete      = "Permission.Org.ProcessStatus.Delete"
	OperationOrgProcessView              = "Permission.Org.Process.View"
	OperationOrgProcessModify            = "Permission.Org.Process.Modify"
	OperationOrgProcessCreate            = "Permission.Org.Process.Create"
	OperationOrgProcessDelete            = "Permission.Org.Process.Delete"
	OperationOrgProcessStepView          = "Permission.Org.ProcessStep.View"
	OperationOrgProcessStepModify        = "Permission.Org.ProcessStep.Modify"
	OperationOrgProcessStepCreate        = "Permission.Org.ProcessStep.Create"
	OperationOrgProcessStepDelete        = "Permission.Org.ProcessStep.Delete"
)

const (
	ManageGroupKeyApp     = "app_ids"
	ManageGroupKeyOptAuth = "opt_auth"
)

const (
	OptAuthList = `[
    {
        "code":"Permission.Org.Config-Modify",
        "name":"编辑团队设置",
        "type":"FuncPermission",
        "group":"团队设置",
        "groupCode":"Permission.Org.Config",
        "showStatus":1
    },
    {
        "code":"Permission.Org.Config-ModifyField",
        "name":"管理自定义字段",
        "type":"FuncPermission",
        "group":"团队设置",
        "groupCode":"Permission.Org.Config",
        "showStatus":1
    },
    {
        "code":"Permission.Org.Config-TplSaveAs",
        "name":"另存为模板",
        "type":"FuncPermission",
        "group":"团队设置",
        "groupCode":"Permission.Org.Config",
        "showStatus":1
    },
    {
        "code":"Permission.Org.Config-TplDelete",
        "name":"删除模板",
        "type":"FuncPermission",
        "group":"团队设置",
        "groupCode":"Permission.Org.Config",
        "showStatus":1
    },
    {
        "code":"Permission.Org.User-ModifyStatus",
        "name":"编辑/审核成员状态",
        "type":"FuncPermission",
        "group":"成员管理",
        "groupCode":"Permission.Org.User",
        "showStatus":1
    },
    {
        "code":"Permission.Org.User-ModifyUserAdminGroup",
        "name":"修改成员角色",
        "type":"FuncPermission",
        "group":"成员管理",
        "groupCode":"Permission.Org.User",
        "showStatus":1
    },
    {
        "code":"Permission.Org.User-ModifyUserDept",
        "name":"修改成员部门",
        "type":"FuncPermission",
        "group":"成员管理",
        "groupCode":"Permission.Org.User",
        "showStatus":1
    },
    {
        "code":"Permission.Org.Department-Create",
        "name":"创建部门",
        "type":"FuncPermission",
        "group":"部门管理",
        "groupCode":"Permission.Org.Department",
        "showStatus":1
    },
    {
        "code":"Permission.Org.Department-Modify",
        "name":"编辑部门",
        "type":"FuncPermission",
        "group":"部门管理",
        "groupCode":"Permission.Org.Department",
        "showStatus":1
    },
    {
        "code":"Permission.Org.Department-Delete",
        "name":"删除部门",
        "type":"FuncPermission",
        "group":"部门管理",
        "groupCode":"Permission.Org.Department",
        "showStatus":1
    },
    {
        "code":"Permission.Org.InviteUser-Invite",
        "name":"邀请成员",
        "type":"FuncPermission",
        "group":"邀请成员",
        "groupCode":"Permission.Org.InviteUser",
        "showStatus":1
    },
	{
		"code": "Permission.Org.AddUser-Add",
		"name": "添加/导入成员",
		"type": "FuncPermission",
		"group": "添加/导入成员",
		"groupCode": "Permission.Org.AddUser",
		"showStatus": 1
	},
    {
        "code":"Permission.Org.AdminGroup-View",
        "name":"查看角色",
        "type":"FuncPermission",
        "group":"角色管理",
        "groupCode":"Permission.Org.AdminGroup",
        "showStatus":1
    },
    {
        "code":"Permission.Org.AdminGroup-Create",
        "name":"创建角色",
        "type":"FuncPermission",
        "group":"角色管理",
        "groupCode":"Permission.Org.AdminGroup",
        "showStatus":1
    },
    {
        "code":"Permission.Org.AdminGroup-Modify",
        "name":"编辑角色",
        "type":"FuncPermission",
        "group":"角色管理",
        "groupCode":"Permission.Org.AdminGroup",
        "showStatus":1
    },
    {
        "code":"Permission.Org.AdminGroup-Delete",
        "name":"删除角色",
        "type":"FuncPermission",
        "group":"角色管理",
        "groupCode":"Permission.Org.AdminGroup",
        "showStatus":1
    },
	{
		"code": "Permission.Org.PersonInfo-Manage",
		"name": "个人信息管理",
		"type": "FuncPermission",
		"group": "个人信息管理",
		"groupCode": "Permission.Org.PersonInfo",
		"showStatus": 1
	},
    {
        "code":"Permission.Org.Project-Create",
        "name":"创建应用",
        "type":"AppPermission",
        "group":"创建应用",
        "groupCode":"Permission.Org.CreateProject",
        "showStatus":1
    },
    {
        "code":"Permission.Org.Project-Manage",
        "name":"应用管理",
        "type":"AppPermission",
        "group":"应用管理",
        "groupCode":"Permission.Org.ManageProject",
        "showStatus":1
    },
    {
        "code":"Permission.Org.Config-Transfer",
        "name":"转让组织",
        "type":"FuncPermission",
        "group":"组织设置管理",
        "groupCode":"Permission.Org.Config",
        "showStatus":2
    },
    {
        "code":"Permission.Org.User-View",
        "name":"查看成员列表",
        "type":"FuncPermission",
        "group":"组织设置管理",
        "groupCode":"Permission.Org.User",
        "showStatus":2
    },
    {
        "code":"MenuPermission.Org-Workspace",
        "name":"工作台",
        "type":"MenuPermission",
        "group":"工作台",
        "groupCode":"MenuPermission.Org-Workspace",
        "showStatus":1
    },
    {
        "code":"MenuPermission.Org-Issue",
        "name":"全部任务",
        "type":"MenuPermission",
        "group":"全部任务",
        "groupCode":"MenuPermission.Org-Issue",
        "showStatus":1
    },
    {
        "code":"MenuPermission.Org-Project",
        "name":"项目",
        "type":"MenuPermission",
        "group":"项目",
        "groupCode":"MenuPermission.Org-Project",
        "showStatus":1
    },
    {
        "code":"MenuPermission.Org-PolarisTpl",
        "name":"模板市场",
        "type":"MenuPermission",
        "group":"模板市场",
        "groupCode":"MenuPermission.Org-PolarisTpl",
        "showStatus":1
    },
    {
        "code":"MenuPermission.Org-Member",
        "name":"成员",
        "type":"MenuPermission",
        "group":"成员",
        "groupCode":"MenuPermission.Org-Member",
        "showStatus":1
    },
    {
        "code":"MenuPermission.Org-Trend",
        "name":"动态",
        "type":"MenuPermission",
        "group":"动态",
        "groupCode":"MenuPermission.Org-Trend",
        "showStatus":1
    },
    {
        "code":"MenuPermission.Org-WorkHour",
        "name":"工时",
        "type":"MenuPermission",
        "group":"工时",
        "groupCode":"MenuPermission.Org-WorkHour",
        "showStatus":1
    },
    {
        "code":"MenuPermission.Org-Trash",
        "name":"回收站",
        "type":"MenuPermission",
        "group":"回收站",
        "groupCode":"MenuPermission.Org-Trash",
        "showStatus":1
    },
    {
        "code":"MenuPermission.Org-Setting",
        "name":"设置",
        "type":"MenuPermission",
        "group":"设置",
        "groupCode":"MenuPermission.Org-Setting",
        "showStatus":1
    },
	{
		"code": "MenuPermission.Org-CreateButton",
		"name": "显示全局新建按钮",
		"type": "MenuPermission",
		"group": "显示全局新建按钮",
		"groupCode": "MenuPermission.Org-CreateButton",
		"showStatus": 1
	}
]`

	OptAuthExtraInfo = `{
    "Permission.Org.User":{
        "desc":{
            "en":"Audit members and modify their status, roles, and departments",
            "zh":"审核成员，修改成员的姓名、登录密码、角色、所属部门和状态"
        },
        "name":{
            "en":"Member management"
        }
    },
    "Permission.Org.Config":{
        "desc":{
            "en":"Basic information, settings for managing your team",
            "zh":"团队基本信息设置、团队字段管理与团队模板管理"
        },
        "name":{
            "en":"Team settings"
        }
    },
    "MenuPermission.Org-Issue":{
        "desc":{
            "en":"Issue",
            "zh":"全部任务"
        },
        "name":{
            "en":"Issue"
        }
    },
    "MenuPermission.Org-Trash":{
        "desc":{
            "en":"Trash for issue and file",
            "zh":"回收站"
        },
        "name":{
            "en":"Trash"
        }
    },
    "MenuPermission.Org-Trend":{
        "desc":{
            "en":"Trend",
            "zh":"动态"
        },
        "name":{
            "en":"Trend"
        }
    },
    "MenuPermission.Org-Member":{
        "desc":{
            "en":"Member",
            "zh":"成员"
        },
        "name":{
            "en":"Member"
        }
    },
    "Permission.Org.AdminGroup":{
        "desc":{
            "en":"Create, edit, delete roles",
            "zh":"创建、编辑、删除角色"
        },
        "name":{
            "en":"Role management"
        }
    },
    "Permission.Org.Department":{
        "desc":{
            "en":"Create, edit, delete departments",
            "zh":"创建、编辑、删除部门"
        },
        "name":{
            "en":"Department management"
        }
    },
    "Permission.Org.InviteUser":{
        "desc":{
            "en":"Invite external users to join the team",
            "zh":"邀请外部用户加入团队"
        },
        "name":{
            "en":"Invite members"
        }
    },
	"Permission.Org.AddUser":{
		"desc":{
			"en":"Add Team Members",
			"zh":"添加单个成员/excel批量导入成员"
		},
		"name":{
			"en":"Add Team Members"
		}
	},
	"Permission.Org.PersonInfo":{
		"desc":{
			"en":"Edit personal information",
			"zh":"编辑个人信息"
		},
		"name":{
			"en":"Edit personal information"
		}
	},
    "MenuPermission.Org-Project":{
        "desc":{
            "en":"Project",
            "zh":"项目"
        },
        "name":{
            "en":"Project"
        }
    },
    "MenuPermission.Org-Setting":{
        "desc":{
            "en":"Organization Settings",
            "zh":"设置"
        },
        "name":{
            "en":"Settings"
        }
    },
    "MenuPermission.Org-WorkHour":{
        "desc":{
            "en":"Work Time",
            "zh":"工时"
        },
        "name":{
            "en":"Work Time"
        }
    },
    "MenuPermission.Org-Workspace":{
        "desc":{
            "en":"Workspace",
            "zh":"工作台"
        },
        "name":{
            "en":"Workspace"
        }
    },
	"MenuPermission.Org-CreateButton":{
		"desc":{
			"en":"Blue plus sign in the upper right corner",
			"zh":"右上角蓝色➕"
		},
		"name":{
			"en":"Blue plus sign in the upper right corner"
		}
	},
    "Permission.Org.CreateProject":{
        "desc":{
            "en":"Create a new project and managing template",
            "zh":"创建新的应用"
        },
        "name":{
            "en":"Create a project"
        }
    },
    "Permission.Org.ManageProject":{
        "desc":{
            "en":"Sync as administrator of the selected app",
            "zh":"成为所选应用的管理员"
        },
        "name":{
            "en":"become administrators of apps"
        }
    },
    "MenuPermission.Org-PolarisTpl":{
        "desc":{
            "en":"Template Market",
            "zh":"模板市场"
        },
        "name":{
            "en":"Template Market"
        }
    }
}`
)
