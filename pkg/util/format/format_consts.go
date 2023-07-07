package format

import "regexp"

const (
	EmailPattern        = `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //用户邮箱
	PasswordPattern     = `^[a-zA-Z]\w*$`                               //用户密码
	AccountPattern      = `^[0-9A-Za-z@_.]{2,20}$`                      //账号
	MobileRegionPattern = `^\+[0-9]{1,6}$`                              //手机区号
	MobilePattern       = `^1[0-9]{10}$`                                //手机号

	OrgNamePattern             = "^[\u4e00-\u9fa5|0-9|a-zA-Z]{1,20}$" //组织名
	OrgCodePattern             = "^[0-9|a-zA-Z]{1,20}$"               //网址后缀编号
	OrgAdressPattern           = `^.{0,100}$`                         //组织地址
	ProjectNamePattern         = `^.{1,30}$`                          //项目名
	ProjectPreviousCodePattern = `^[a-zA-Z|0-9]{0,50}$`               //项目前缀编号
	ProjectRemarkPattern       = `^[\s\S]{0,500}$`                    //项目描述(简介)
	//ProjectNoticePattern       = `^.{0,2000}$`                                 //项目公告
	IssueNamePattern = `^.{1,500}$` //任务名
	//IssueRemarkPattern           = `^.{0,10000}$`                                //任务描述(详情)
	IssueCommenPattern           = `^.{1,200}$` //任务评论
	ProjectObjectTypeNamePattern = `^.{1,30}$`  //标题栏名
	//ResourceNamePattern          = `^[^\\\\/:*?\"<>|]{1,300}(\.[a-zA-Z0-9]+)?$` //资源名
	ResourceNamePattern = `^.{1,300}$`               //资源名
	FolderNamePattern   = `^[^\\\\/:*?\"<>|]{1,30}$` //文件夹名

	RoleNamePattern        = "^[\\pP\\pL\\pS\\pN ]{1,20}$" // RoleNamePattern 角色名。允许空格，因为英文名会带空格。
	RoleGroupNamePattern   = "^[\\pP\\pL\\pS\\pN ]{1,20}$" // RoleGroupNamePattern 角色组名。允许空格，因为英文名会带空格。
	ManageGroupNamePattern = "^[\\pP\\pL\\pS\\pN ]{1,10}$" // ManageGroupNamePattern 管理员组名。允许空格，因为英文名会带空格。

	PositionNamePattern = "^[\\pP\\pL\\pS\\pN]{1,10}$" // positionNamePattern 职级名 1-10位

)

const (
	ChinesePattern         = "[\u4e00-\u9fa5]+?" // 中文
	SpacePattern           = " +?"               // 空格
	UpperCaseLetterPattern = "[A-Z]+"            // 大写字母
	LowerCaseLetterPattern = "[a-z]+"            // 小写字母
	NumberPattern          = "\\d+"              // 数字

)

var ChineseReg = regexp.MustCompile(ChinesePattern)
var SpaceReg = regexp.MustCompile(SpacePattern)
var UpperCaseLetterReg = regexp.MustCompile(UpperCaseLetterPattern)
var LowerCaseLetterReg = regexp.MustCompile(LowerCaseLetterPattern)
var NumberReg = regexp.MustCompile(NumberPattern)
