package consts

const (
	TcId                = "id"
	TcName              = "name"
	TcNamePinyin        = "name_pinyin"
	TcCode              = "code"
	TcSecret1           = "secret1"
	TcSecret2           = "secret2"
	TcOwner             = "owner"
	TcOwnerChangeTime   = "owner_change_time"
	TcCheckStatus       = "check_status"
	TcStatus            = "status"
	TcStatDate          = "stat_date"
	TcCreator           = "creator"
	TcCreateTime        = "create_time"
	TcUpdator           = "updator"
	TcUpdateTime        = "update_time"
	TcVersion           = "version"
	TcIsDelete          = "is_delete"
	TcDelFlag           = "del_flag"
	TcSysVersion        = "sys_version"
	TcFeatureInfo       = "feature_info"
	TcBugFixInfo        = "bug_fix_info"
	TcChangeLnfo        = "change_lnfo"
	TcDeprecatedInfo    = "deprecated_info"
	TcExt               = "ext"
	TcReleaseTime       = "release_time"
	TcK                 = "k"
	TcV                 = "v"
	TcOrgId             = "org_id"
	TcMaxId             = "max_id"
	TcStep              = "step"
	TcLangCode          = "lang_code"
	TcStorage           = "storage"
	TcMemberCount       = "member_count"
	TcPrice             = "price"
	TcMemberPrice       = "member_price"
	TcDuration          = "duration"
	TcIsShow            = "is_show"
	TcSort              = "sort"
	TcRemark            = "remark"
	TcTopic             = "topic"
	TcMessageKey        = "message_key"
	TcMessage           = "message"
	TcGroupName         = "group_name"
	TcMessageId         = "message_id"
	TcLastConsumerTime  = "last_consumer_time"
	TcFailCount         = "fail_count"
	TcFailTime          = "fail_time"
	TcTimeZone          = "time_zone"
	TcTimeDifference    = "time_difference"
	TcPayLevel          = "pay_level"
	TcPayStartTime      = "pay_start_time"
	TcPayEndTime        = "pay_end_time"
	TcWebSite           = "web_site"
	TcLanguage          = "language"
	TcRemindSendTime    = "remind_send_time"
	TcDatetimeFormat    = "datetime_format"
	TcPasswordLength    = "password_length"
	TcPasswordRule      = "password_rule"
	TcMaxLoginFailCount = "max_login_fail_count"
	TcEmailStatus       = "email_status"
	TcSmtpServer        = "smtp_server"
	TcSmtpPort          = "smtp_port"
	TcSmtpUserName      = "smtp_user_name"
	TcSmtpPassword      = "smtp_password"
	TcEmailFormat       = "email_format"
	TcSenderAddress     = "sender_address"
	TcEmailEncode       = "email_encode"
	TcIndustryId        = "industry_id"
	TcScale             = "scale"
	TcSourceChannel     = "source_channel"
	TcContinentId       = "continent_id"
	TcCountryId         = "country_id"
	//这个就是省的id和TcProvinceId一样在不同表不一样的叫法
	TcStateId                         = "state_id"
	TcProvinceId                      = "province_id"
	TcCityId                          = "city_id"
	TcAddress                         = "address"
	TcLogoUrl                         = "logo_url"
	TcResourceId                      = "resource_id"
	TcFolderId                        = "folder_id"
	TcIsAuthenticated                 = "is_authenticated"
	TcInitStatus                      = "init_status"
	TcInitVersion                     = "init_version"
	TcOutOrgId                        = "out_org_id"
	TcIndustry                        = "industry"
	TcAuthTicket                      = "auth_ticket"
	TcAuthLevel                       = "auth_level"
	TcLoginName                       = "login_name"
	TcLoginNameEditCount              = "login_name_edit_count"
	TcEmail                           = "email"
	TcMobile                          = "mobile"
	TcMobileRegion                    = "mobile_region"
	TcAvatar                          = "avatar"
	TcBirthday                        = "birthday"
	TcSex                             = "sex"
	TcPassword                        = "password"
	TcPasswordSalt                    = "password_salt"
	TcLastEditPwdTime                 = "last_edit_pwd_time"
	TcMotto                           = "motto"
	TcLastLoginIp                     = "last_login_ip"
	TcLoginIp                         = "login_ip"
	TcLastLoginTime                   = "last_login_time"
	TcLoginFailCount                  = "login_fail_count"
	TcUserId                          = "user_id"
	TcDailyReportMessageStatus        = "daily_report_message_status"
	TcDefaultProjectId                = "default_project_id"
	TcDailyProjectReportMessageStatus = "daily_project_report_message_status"
	TcOwnerRangeStatus                = "owner_range_status"
	TcParticipantRangeStatus          = "participant_range_status"
	TcAttentionRangeStatus            = "attention_range_status"
	TcCreateRangeStatus               = "create_range_status"
	TcRemindMessageStatus             = "remind_message_status"
	TcCommentAtMessageStatus          = "comment_at_message_status"
	TcModifyMessageStatus             = "modify_message_status"
	TcRelationMessageStatus           = "relation_message_status"
	TcDefaultProjectObjectTypeId      = "default_project_object_type_id"
	TcPcNoticeOpenStatus              = "pc_notice_open_status"
	TcPcIssueRemindMessageStatus      = "pc_issue_remind_message_status"
	TcPcOrgMessageStatus              = "pc_org_message_status"
	TcPcProjectMessageStatus          = "pc_project_message_status"
	TcPcCommentAtMessageStatus        = "pc_comment_at_message_status"

	TcUseStatus           = "use_status"
	TcOutOrgUserId        = "out_org_user_id"
	TcOutUserId           = "out_user_id"
	TcIsActive            = "is_active"
	TcJobNumber           = "job_number"
	TcProjectId           = "project_id"
	TcProjectObjectTypeId = "project_object_type_id"
	TcTitle               = "title"
	TcPriorityId          = "priority_id"
	TcSourceId            = "source_id"
	TcIssueObjectTypeId   = "issue_object_type_id"
	TcPlanStartTime       = "plan_start_time"
	TcPlanEndTime         = "plan_end_time"
	TcStartTime           = "start_time"
	TcEndTime             = "end_time"
	TcPlanWorkHour        = "plan_work_hour"
	TcIterationId         = "iteration_id"
	TcVersionId           = "version_id"
	TcModuleId            = "module_id"
	TcParentId            = "parent_id"
	TcIsHide              = "is_hide"
	TcIssueId             = "issue_id"
	TcStoryPoint          = "story_point"
	TcTags                = "tags"
	TcTagId               = "tag_id"
	TcRelationId          = "relation_id"
	TcRelationType        = "relation_type"
	TcPackage             = "package"
	TcIcon                = "icon"
	TcAppId               = "app_id"
	TcAppVersion          = "app_version"
	TcPreCode             = "pre_code"
	TcProjectTypeId       = "project_type_id"
	TcPublicStatus        = "public_status"
	TcIsFiling            = "is_filing"
	TcTeamId              = "team_id"
	TcType                = "type"
	TcBgStyle             = "bg_style"
	TcFontStyle           = "font_style"
	TcIsDefault           = "is_default"

	TcPath               = "path"
	TcSuffix             = "suffix"
	TcMd5                = "md5"
	TcPermissionId       = "permission_id"
	TcOperationCodes     = "operation_codes"
	TcIsModifyPermission = "is_modify_permission"
	TcRoleGroupId        = "role_group_id"
	TcRoleId             = "role_id"
	TcDepartmentId       = "department_id"
	TcOutOrgDepartmentId = "out_org_department_id"

	TcAuditorId        = "auditor_id"
	TcAuditTime        = "audit_time"
	TcStatusChangerId  = "status_changer_id"
	TcStatusChangeTime = "status_change_time"
	TcSourceType       = "source_type" //2019/12/25新增资源文件来源类型
	TcFileType         = "file_type"   //2019/12/30新增文件类型
	TcRelationCode     = "relation_code"
	TcRemarkDetail     = "remark_detail"
	TcStatusId         = "status_id"
	TcPropertyId       = "property_id"
	TcFunctionCode     = "function_code"
	TcRemindBindPhone  = "remind_bind_phone"
	TcIsLeader         = "is_leader"
	TcIsRegister       = "is_register"
	TcLastInviteTime   = "last_invite_time"

	TcUserIds       = "user_ids"
	TcAppPackageIds = "app_package_ids"
	TcAppIds        = "app_ids"
	TcUsageIds      = "usage_ids"
	TcOptAuth       = "opt_auth"
	TcDeptIds       = "dept_ids"
	TcRoleIds       = "role_ids"

	TcDbId = "db_id"
	TcDcId = "dc_id"
	TcDsId = "ds_id"

	TcApiKey = "api_key"

	TcPositionLevel = "position_level"
	TcPositionId    = "position_id"
	TcOrgPositionId = "org_position_id"

	TcEmpNo    = "emp_no"
	TcWeiboIds = "weibo_ids"

	TcGlobalUserId = "global_user_id"
)
