package consts

const (
	RoleSysAdmin                 = "Role.Sys.Admin"
	RoleSysManager               = "Role.Sys.Manager"
	RoleSysMember                = "Role.Sys.Member"
	RoleGroupSpecialCreator      = "RoleGroup.Special.Creator"
	RoleGroupSpecialOwner        = "RoleGroup.Special.Owner"
	RoleGroupSpecialWorker       = "RoleGroup.Special.Worker"
	RoleGroupSpecialAttention    = "RoleGroup.Special.Attention"
	RoleGroupSpecialMember       = "RoleGroup.Special.Member"
	RoleGroupSpecialVisitor      = "RoleGroup.Special.Visitor"
	RoleGroupOrgAdmin            = "RoleGroup.Org.Admin"
	RoleGroupOrgManager          = "RoleGroup.Org.Manager"
	RoleGroupProProjectManager   = "RoleGroup.Pro.ProjectManager"
	RoleGroupProTechnicalManager = "RoleGroup.Pro.TechnicalManager"
	RoleGroupProProductManager   = "RoleGroup.Pro.ProductManager"
	RoleGroupProDeveloper        = "RoleGroup.Pro.Developer"
	RoleGroupProTester           = "RoleGroup.Pro.Tester"
	RoleGroupProMember           = "RoleGroup.Pro.Member"
)

const (
	ManageGroupSys              = "ManageGroup.Sys" // 超级管理员 1
	ManageGroupSub              = "ManageGroup.Sub"
	ManageGroupSubNormalAdmin   = "ManageGroup.Sub.NormalAdmin" // bjx：普通管理员   3 // 这个数字和创建管理组入参的 groupType 对应
	ManageGroupSubNormalUser    = "ManageGroup.Sub.NormalUser"  // bjx：团队成员  4
	ManageGroupSubNormalUserBjx = "RoleGroup.Special.Member"    // bjx：团队成员 5 (bjx 同步过来的 langCode，同样代表组织成员)
	ManageGroupSubUserCustom    = "ManageGroup.Sub.UserCustom"  // 用户创建的管理组 6
)

const (
	RoleGroupSys     = "RoleGroup.Sys"
	RoleGroupSpecial = "RoleGroup.Special"
	RoleGroupOrg     = "RoleGroup.Org"
	RoleGroupPro     = "RoleGroup.Pro"
)

const (
	RoleOperationView             = "View"
	RoleOperationModify           = "Modify"
	RoleOperationDelete           = "Delete"
	RoleOperationCreate           = "CreateOrgMember"
	RoleOperationCheck            = "Check"
	RoleOperationInvite           = "Invite"
	RoleOperationBind             = "Bind"
	RoleOperationUnbind           = "Unbind"
	RoleOperationAttention        = "Attention"
	RoleOperationUnAttention      = "UnAttention"
	RoleOperationModifyStatus     = "ModifyStatus"
	RoleOperationComment          = "Comment"
	RoleOperationTransfer         = "Transfer"
	RoleOperationInit             = "Init"
	RoleOperationDrop             = "Drop"
	RoleOperationFiling           = "Filing"
	RoleOperationUnFiling         = "UnFiling"
	RoleOperationUpload           = "Upload"
	RoleOperationDownload         = "Download"
	RoleOperationRemove           = "Remove"
	RoleOperationModifyPermission = "ModifyPermission"
	RoleOperationCreateFolder     = "CreateFolder"
	RoleOperationModifyFolder     = "ModifyFolder"
	RoleOperationDeleteFolder     = "DeleteFolder"
)

const (
	RoleOperationPathSys                  = "/Sys"
	RoleOperationPathSysDic               = "/Sys/Dic"
	RoleOperationPathSysSource            = "/Sys/Source"
	RoleOperationPathSysPayLevel          = "/Sys/PayLevel"
	RoleOperationPathOrgProRole           = "/Org/{org_id}/Pro/{pro_id}/Role"
	RoleOperationPathOrgProProjectVersion = "/Org/{org_id}/Pro/{pro_id}/ProjectVersion"
	RoleOperationPathOrgProProjectModule  = "/Org/{org_id}/Pro/{pro_id}/ProjectModule"
	RoleOperationPathOrgPro               = "/Org/{org_id}/Pro/{pro_id}"
	RoleOperationPathOrgProIteration      = "/Org/{org_id}/Pro/{pro_id}/Iteration"
	RoleOperationPathOrgProIssueTT        = "/Org/{org_id}/Pro/{pro_id}/Issue/6"
	RoleOperationPathOrgProIssueT         = "/Org/{org_id}/Pro/{pro_id}/Issue/4"
	RoleOperationPathOrgProIssueF         = "/Org/{org_id}/Pro/{pro_id}/Issue/2"
	RoleOperationPathOrgProIssueD         = "/Org/{org_id}/Pro/{pro_id}/Issue/3"
	RoleOperationPathOrgProIssueB         = "/Org/{org_id}/Pro/{pro_id}/Issue/5"
	RoleOperationPathOrgProIssue          = "/Org/{org_id}/Pro/{pro_id}/Issue"
	RoleOperationPathOrgProProConfig      = "/Org/{org_id}/Pro/{pro_id}/ProConfig"
	RoleOperationPathOrgProComment        = "/Org/{org_id}/Pro/{pro_id}/Comment"
	RoleOperationPathOrgProBan            = "/Org/{org_id}/Pro/{pro_id}/Ban"
	RoleOperationPathOrgUser              = "/Org/{org_id}/User"
	RoleOperationPathOrgTeam              = "/Org/{org_id}/Team"
	RoleOperationPathOrgRoleGroup         = "/Org/{org_id}/RoleGroup"
	RoleOperationPathOrgRole              = "/Org/{org_id}/Role"
	RoleOperationPathOrgProjectType       = "/Org/{org_id}/ProjectType"
	RoleOperationPathOrgProjectObjectType = "/Org/{org_id}/ProjectObjectType"
	RoleOperationPathOrgProject           = "/Org/{org_id}/Project"
	RoleOperationPathOrgProcessStep       = "/Org/{org_id}/ProcessStep"
	RoleOperationPathOrgProcessStatus     = "/Org/{org_id}/ProcessStatus"
	RoleOperationPathOrgProcess           = "/Org/{org_id}/Process"
	RoleOperationPathOrgPriority          = "/Org/{org_id}/Priority"
	RoleOperationPathOrg                  = "/Org/{org_id}"
	RoleOperationPathOrgMessageConfig     = "/Org/{org_id}/MessageConfig"
	RoleOperationPathOrgIssueSource       = "/Org/{org_id}/IssueSource"
	RoleOperationPathOrgOrgConfig         = "/Org/{org_id}/OrgConfig"
	RoleOperationPathOrgProMember         = "/Org/{org_id}/Pro/{pro_id}/Member"
	RoleOperationPathOrgProFile           = "/Org/{org_id}/Pro/{pro_id}/File"
	RoleOperationPathOrgProTag            = "/Org/{org_id}/Pro/{pro_id}/Tag"
	RoleOperationPathOrgProAttachment     = "/Org/{org_id}/Pro/{pro_id}/Attachment"
)

//权限项
const (
	PermissionSysSys               = "Permission.Sys.Sys"
	PermissionSysDic               = "Permission.Sys.Dic"
	PermissionSysSource            = "Permission.Sys.Source"
	PermissionSysPayLevel          = "Permission.Sys.PayLevel"
	PermissionOrgOrg               = "Permission.Org.Org"
	PermissionOrgConfig            = "Permission.Org.Config"
	PermissionOrgMessageConfig     = "Permission.Org.MessageConfig"
	PermissionOrgUser              = "Permission.Org.User"
	PermissionOrgTeam              = "Permission.Org.Team"
	PermissionOrgRoleGroup         = "Permission.Org.RoleGroup"
	PermissionOrgRole              = "Permission.Org.Role"
	PermissionOrgProject           = "Permission.Org.Project"
	PermissionOrgProjectType       = "Permission.Org.ProjectType"
	PermissionOrgIssueSource       = "Permission.Org.IssueSource"
	PermissionOrgProjectObjectType = "Permission.Org.ProjectObjectType"
	PermissionOrgPriority          = "Permission.Org.Priority"
	PermissionOrgProcessStatus     = "Permission.Org.ProcessStatus"
	PermissionOrgProcess           = "Permission.Org.Process"
	PermissionOrgProcessStep       = "Permission.Org.ProcessStep"
	PermissionProPro               = "Permission.Pro.Pro"
	PermissionProConfig            = "Permission.Pro.Config"
	PermissionProBan               = "Permission.Pro.Ban"
	PermissionProIteration         = "Permission.Pro.Iteration"
	PermissionProIssue             = "Permission.Pro.Issue"
	PermissionProIssue2            = "Permission.Pro.Issue.2"
	PermissionProIssue3            = "Permission.Pro.Issue.3"
	PermissionProIssue4            = "Permission.Pro.Issue.4"
	PermissionProIssue5            = "Permission.Pro.Issue.5"
	PermissionProIssue6            = "Permission.Pro.Issue.6"
	PermissionProComment           = "Permission.Pro.Comment"
	PermissionProProjectVersion    = "Permission.Pro.ProjectVersion"
	PermissionProProjectModule     = "Permission.Pro.ProjectModule"
	PermissionProRole              = "Permission.Pro.Role"
	PermissionProTest              = "Permission.Pro.Test"
	PermissionProTestTestApp       = "Permission.Pro.Test.TestApp"
	PermissionProTestTestDevice    = "Permission.Pro.Test.TestDevice"
	PermissionProTestTestReport    = "Permission.Pro.Test.TestReport"
	PermissionProFile              = "Permission.Pro.File"
	PermissionProTag               = "Permission.Pro.Tag"
	PermissionProAttachment        = "Permission.Pro.Attachment"
	PermissionProMember            = "Permission.Pro.Member"

	MenuPermissionOrgTrash = "MenuPermission.Org-Trash" // 菜单权限项：回收站

	PermissionAppManage = "Permission.Org.Project-Manage"

	MenuPermissionTrend = "MenuPermission.Org-Trend"

	PermissionOrgAddUser = "Permission.Org.AddUser-Add"
)
