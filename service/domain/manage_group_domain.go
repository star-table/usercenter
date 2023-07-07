package domain

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/snowflake"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/store/mysql"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/pkg/util/jsonx"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/service/model/po"
	"github.com/star-table/usercenter/service/model/req"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

var (
	AllOptAuthJson = jsonx.ToJsonIgnoreError(consts.AllManageGroupOptAuths)
	EmptyArrayJson = jsonx.ToJsonIgnoreError([]int64{})
)

// CreateManageGroup 创建管理组
func CreateManageGroup(orgId, operator int64, input req.CreateManageGroup, tx sqlbuilder.Tx) (int64, error) {
	groupLangCode := consts.ManageGroupSub

	appIds := EmptyArrayJson

	switch input.GroupType {
	case 1:
		groupLangCode = consts.ManageGroupSys
		input.Name = consts.DefaultManageGroupName[0]
		appIds = jsonx.ToJsonIgnoreError([]int64{-1})
	case 3: // bjx-普通管理员
		groupLangCode = consts.ManageGroupSubNormalAdmin
		if len(input.OptAuth) == 0 {
			input.OptAuth = GetAdminGroupOptAuthForPolaris("普通管理员")
		}
		// 极星普通管理员默认可以管理所有应用 -1表示可以管理所有应用
		appIds = jsonx.ToJsonIgnoreError([]int64{-1})
	case 4: // bjx-普通成员
		groupLangCode = consts.ManageGroupSubNormalUser
		if len(input.OptAuth) == 0 {
			input.OptAuth = GetAdminGroupOptAuthForPolaris("团队成员")
		}
	case 6: // 用户创建的管理组
		groupLangCode = consts.ManageGroupSubUserCustom
		if len(input.OptAuth) == 0 {
			input.OptAuth = GetAdminGroupOptAuthForPolaris("团队成员")
		}
	}
	opAuthJson := EmptyArrayJson
	if input.OptAuth != nil && len(input.OptAuth) > 0 {
		opAuthJson = json.ToJsonIgnoreError(input.OptAuth)
	}

	group := po.LcPerManageGroup{
		Id:            snowflake.Id(),
		LangCode:      groupLangCode,
		OrgId:         orgId,
		Name:          input.Name,
		UserIds:       EmptyArrayJson, // 初始化防止插入null
		OptAuth:       opAuthJson,     // 初始化防止插入null
		DeptIds:       EmptyArrayJson, // 初始化防止插入null
		RoleIds:       EmptyArrayJson, // 初始化防止插入null
		AppPackageIds: EmptyArrayJson, // 初始化防止插入null
		AppIds:        appIds,         // 初始化防止插入null
		Creator:       operator,
		Updator:       operator,
	}

	if groupLangCode == consts.ManageGroupSub {
		// 默认开启全部操作权限
		group.OptAuth = AllOptAuthJson
	}
	dbErr := store.Mysql.TransInsert(tx, &group)
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}

	return group.Id, nil
}

// UpdateManageGroup 更新管理组
func UpdateManageGroup(orgId, operator, groupId int64, reqParam req.UpdateManageGroup) (int, error) {

	//加入更新人
	count, dbErr := store.Mysql.UpdateSmartWithCond(consts.TableManageGroup, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcId:       groupId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcName:    reqParam.Name,
		consts.TcUpdator: operator,
		consts.TcVersion: db.Raw("version + 1"),
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}
	return int(count), nil
}

// DeleteManageGroup 删除管理组
func DeleteManageGroup(orgId, operator, groupId int64) (int, error) {

	//加入更新人
	count, dbErr := store.Mysql.UpdateSmartWithCond(consts.TableManageGroup, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcId:       groupId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcUpdator:  operator,
		consts.TcVersion:  db.Raw("version + 1"),
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}
	return int(count), nil
}

// UpdateManageGroupContents 更新管理组内容信息
func UpdateManageGroupContents(orgId, operator int64, groupId int64, langCode string, input req.UpdateManageGroupContents) (int, error) {
	valueStr := jsonx.ToJsonIgnoreError(input.Values)
	if len(input.Values) < 1 && input.ValueIf != nil {
		valueStr = jsonx.ToJsonIgnoreError(input.ValueIf)
		if input.Key == consts.TcUserIds {
			// 对于 user_ids， 前端可能传过来 int64、string 混合的数组，因此做一次统一处理
			userIdArr := TransferUserIdsFromUserIdJson(valueStr)
			valueStr = jsonx.ToJsonIgnoreError(userIdArr)
		}
	}

	upd := mysql.Upd{
		input.Key:        valueStr, //此处为了使用mysql的json函数，故把int64转为string
		consts.TcUpdator: operator,
		consts.TcVersion: db.Raw("version + 1"),
	}

	// 如果是管理员组，关闭应用管理权限，appIds需要置为空
	if langCode == consts.ManageGroupSubNormalAdmin && input.Key == consts.ManageGroupKeyOptAuth &&
		!strings.Contains(valueStr, consts.PermissionAppManage) {
		upd[consts.ManageGroupKeyApp] = jsonx.ToJsonIgnoreError([]string{})
	}

	// 加入更新人
	count, dbErr := store.Mysql.UpdateSmartWithCond(consts.TableManageGroup, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcId:       groupId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd)
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}
	return int(count), nil
}

// GetSysManageGroup 获取系统管理组
func GetSysManageGroup(orgId int64) (*po.LcPerManageGroup, error) {
	var group po.LcPerManageGroup
	dbErr := store.Mysql.SelectOneByCond(consts.TableManageGroup, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcLangCode: consts.ManageGroupSys,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, &group)

	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return nil, errs.ManageGroupNotExist
		}
		return nil, dbErr
	}
	return &group, nil
}

// GetManageGroupListByOrg 获取组织的管理组列表
func GetManageGroupListByOrg(orgId int64) ([]po.LcPerManageGroup, error) {
	//拼装条件
	cond := db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
	}
	var groups []po.LcPerManageGroup
	dbErr := store.Mysql.SelectAllByCondWithNumAndOrder(consts.TableManageGroup, cond, nil, 0, -1, "create_time asc", &groups)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	return groups, nil
}

// GetManageGroupByCond 根据条件获取管理组
func GetManageGroupByCond(orgId int64, cond db.Cond) ([]po.LcPerManageGroup, error) {
	//拼装条件
	if cond == nil || len(cond) < 1 {
		cond = db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}
	} else {
		cond[consts.TcOrgId] = orgId
	}
	var groups []po.LcPerManageGroup
	dbErr := store.Mysql.SelectAllByCondWithNumAndOrder(consts.TableManageGroup, cond, nil, 0, -1, "create_time asc", &groups)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}

	return groups, nil
}

// GetManageGroupListByUsers
func GetManageGroupListByUsers(orgId int64, userIds []int64) ([]po.LcPerManageGroup, error) {
	conn, dbErr := store.Mysql.GetConnect()
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	conn.SetLogging(true)
	//拼装条件
	cond := db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
	}
	userIdsStr := make([]string, len(userIds))
	for i, v := range userIds {
		userIdsStr[i] = strconv.FormatInt(v, 10)
	}
	var groups []po.LcPerManageGroup
	dbErr = conn.Select(db.Raw("*")).
		From(consts.TableManageGroup).
		Where(cond).
		And(
			db.Or(
				db.Raw(fmt.Sprintf("JSON_OVERLAPS(`user_ids` -> '$', CAST('%s' AS JSON))", json.ToJsonIgnoreError(userIds))),
				db.Raw(fmt.Sprintf("JSON_OVERLAPS(`user_ids` -> '$', CAST('%s' AS JSON))", json.ToJsonIgnoreError(userIdsStr))),
			),
		).
		All(&groups)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	return groups, nil
}

// GetManageGroupListByUser 根据用户ID获取管理组
func GetManageGroupListByUser(orgId int64, userId int64) (*po.LcPerManageGroup, error) {
	list, err := GetManageGroupListByUsers(orgId, []int64{userId})
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, db.ErrNoMoreRows
	}
	return &list[0], nil
}

// GetManageGroup 获取管理组
func GetManageGroup(orgId, groupId int64) (*po.LcPerManageGroup, error) {
	var group po.LcPerManageGroup
	dbErr := store.Mysql.SelectOneByCond(consts.TableManageGroup, db.Cond{
		consts.TcId:       groupId,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, &group)

	if dbErr != nil {
		return nil, dbErr
	}

	return &group, nil
}

// AppendContentsToManageSubGroup 更新管理组内容信息
func AppendContentsToManageSubGroup(orgId int64, operatorUid int64, groupId int64, key string, value int64) (int, error) {
	count, dbErr := store.Mysql.UpdateSmartWithCond(
		consts.TableManageGroup,
		db.Cond{
			consts.TcId:       groupId,
			consts.TcOrgId:    orgId,
			consts.TcLangCode: consts.ManageGroupSub,
			consts.TcIsDelete: consts.AppIsNoDelete,
			db.Raw(fmt.Sprintf("json_search(`%s`, 'one', ?)", consts.TcUserIds), operatorUid): db.IsNotNull(),
			db.Raw(fmt.Sprintf("json_search(`%s`, 'one', ?)", key), value):                    db.IsNull(),
		},
		mysql.Upd{
			key:              db.Raw(fmt.Sprintf("JSON_ARRAY_APPEND(`%s`,'$', ?)", key), strconv.FormatInt(value, 10)),
			consts.TcUpdator: operatorUid,
			consts.TcVersion: db.Raw("version + 1"),
		},
	)
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}

	return int(count), nil
}

// RemoveContentsFromSubManageGroup 删除组织中所有管理组关联的信息
func RemoveContentsFromSubManageGroup(orgId int64, operatorUid int64, key string, value int64) (int, error) {
	count, dbErr := store.Mysql.UpdateSmartWithCond(
		consts.TableManageGroup,
		db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcLangCode: consts.ManageGroupSub,
			consts.TcIsDelete: consts.AppIsNoDelete,
			db.Raw(fmt.Sprintf("json_search(`%s`, 'one', ?)", key), value): db.IsNotNull(),
		},
		mysql.Upd{
			key:              db.Raw(fmt.Sprintf("json_remove(`%s`,JSON_UNQUOTE(json_search(`%s`, 'one', ?)))", key, key), value),
			consts.TcUpdator: operatorUid,
			consts.TcVersion: db.Raw("version + 1"),
		},
	)
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}
	return int(count), nil
}

// RemoveUserFromAdminGroup 删除管理组中的某个人
// 该用户 id 在 user_ids 中的索引位置。
func RemoveUserFromAdminGroup(id int64, operateUid int64, index int, tx sqlbuilder.Tx) (int64, error) {
	var (
		count int64
		dbErr error
	)
	cond1 := db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	upd1 := mysql.Upd{
		consts.TcUserIds: db.Raw(fmt.Sprintf("json_remove(`%s`, '$[%d]')", consts.TcUserIds, index)),
		consts.TcUpdator: operateUid,
		consts.TcVersion: db.Raw("version + 1"),
	}
	if tx != nil {
		count, dbErr = store.Mysql.TransUpdateSmartWithCond(
			tx,
			consts.TableManageGroup,
			cond1,
			upd1,
		)
	} else {
		count, dbErr = store.Mysql.UpdateSmartWithCond(
			consts.TableManageGroup,
			cond1,
			upd1,
		)
	}

	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}
	return count, nil
}

// AppendUserIntoAdminGroup 向管理组中增加一个人
func AppendUserIntoAdminGroup(id int64, operateUid int64, userId int64, tx sqlbuilder.Tx) (int64, error) {
	var (
		count int64
		dbErr error
	)
	cond1 := db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	upd1 := mysql.Upd{
		consts.TcUserIds: db.Raw(fmt.Sprintf("JSON_ARRAY_APPEND(`%s`, '$', %d)", consts.TcUserIds, userId)),
		consts.TcUpdator: operateUid,
		consts.TcVersion: db.Raw("version + 1"),
	}
	if tx != nil {
		count, dbErr = store.Mysql.TransUpdateSmartWithCond(
			tx,
			consts.TableManageGroup,
			cond1,
			upd1,
		)
	} else {
		count, dbErr = store.Mysql.UpdateSmartWithCond(
			consts.TableManageGroup,
			cond1,
			upd1,
		)
	}

	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}
	return count, nil
}

func ManageGroupInitDefault(orgId, creatorId int64) (int64, error) {
	var id int64
	var dbErr error
	dbErr = store.Mysql.TransX(func(tx sqlbuilder.Tx) error {
		//创建系统管理组
		id, dbErr = CreateManageGroup(orgId, creatorId, req.CreateManageGroup{
			GroupType: 1,
		}, tx)
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}

		return nil
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}

	return id, nil
}

// ManageGroupInitForPolaris 给极星应用中的组织初始化管理组
func ManageGroupInitForPolaris(orgId, creatorId int64) (int64, error) {
	var id int64
	var dbErr error
	// 先查询是否存在这些默认管理组
	//	* 如果存在，则忽略
	//  * 如果不存在，则新增
	list, err := GetManageGroupListByOrg(orgId)
	if err != nil {
		logger.Error(err)
		return 0, err
	}
	defaultManageGroupNameMap := map[string]int{
		"超级管理员": 1,
		"普通管理员": 3,
		"团队成员":  4,
	}
	defaultManageGroupNameList := make([]string, 0)
	for name, _ := range defaultManageGroupNameMap {
		defaultManageGroupNameList = append(defaultManageGroupNameList, name)
	}
	existNames := make([]string, 0)
	for _, poObj := range list {
		existNames = append(existNames, poObj.Name)
	}
	needAddNames := slice.ArrayDiffString(defaultManageGroupNameList, existNames)
	dbErr = store.Mysql.TransX(func(tx sqlbuilder.Tx) error {
		for _, groupName := range needAddNames {
			optAuth := GetAdminGroupOptAuthForPolaris(groupName)
			if groupType, ok := defaultManageGroupNameMap[groupName]; ok {
				// 管理组
				id, dbErr = CreateManageGroup(orgId, -1, req.CreateManageGroup{
					Name:      groupName,
					GroupType: groupType,
					OptAuth:   optAuth,
				}, tx)
				if dbErr != nil {
					logger.Error(dbErr)
					return dbErr
				}
			}
		}
		return nil
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}

	return id, nil
}

func GetAdminGroupOptAuthForPolaris(name string) []string {
	configMap := map[string][]string{
		"超级管理员": []string{
			// 这个就特殊处理，默认就支持所有权限项
		},
		"普通管理员": []string{
			"Permission.Org.AdminGroup-View", "Permission.Org.AdminGroup-Create", "Permission.Org.AdminGroup-Modify", "Permission.Org.AdminGroup-Delete",
			"Permission.Org.Config-Modify", "Permission.Org.Config-ModifyField", "Permission.Org.Config-TplSaveAs", "Permission.Org.Config-TplDelete",
			"Permission.Org.User-ModifyStatus", "Permission.Org.User-ModifyUserAdminGroup", "Permission.Org.User-ModifyUserDept", "Permission.Org.User-View",
			"Permission.Org.Department-Create", "Permission.Org.Department-Modify", "Permission.Org.Department-Delete",
			"Permission.Org.InviteUser-Invite", "Permission.Org.AddUser-Add", "Permission.Org.PersonInfo-Manage",
			"Permission.Org.Project-Create", "Permission.Org.Project-Manage",
			// 菜单权限项 code
			"MenuPermission.Org-Workspace", "MenuPermission.Org-Issue", "MenuPermission.Org-Project", "MenuPermission.Org-PolarisTpl",
			"MenuPermission.Org-Member", "MenuPermission.Org-Trend", "MenuPermission.Org-WorkHour", "MenuPermission.Org-Setting",
			"MenuPermission.Org-Trash", "MenuPermission.Org-CreateButton",
		},
		"团队成员": []string{
			"Permission.Org.Project-Create", "Permission.Org.User-View", "Permission.Org.InviteUser-Invite", "Permission.Org.PersonInfo-Manage",
			// 菜单权限项 code
			"MenuPermission.Org-Workspace", "MenuPermission.Org-Issue", "MenuPermission.Org-Project", "MenuPermission.Org-PolarisTpl",
			"MenuPermission.Org-Member", "MenuPermission.Org-Trend", "MenuPermission.Org-WorkHour", "MenuPermission.Org-Setting",
			"MenuPermission.Org-Trash", "MenuPermission.Org-CreateButton",
		},
	}
	if val, ok := configMap[name]; ok {
		return val
	}
	return make([]string, 0)
}

// GetOrgDefaultAdminGroupForPolaris bjx 获取组织默认管理组。一个成员如果没有分配管理组，则默认一个管理组（“团队成员”管理组）
func GetOrgDefaultAdminGroupForPolaris(orgId int64) (*po.LcPerManageGroup, error) {
	list, dbErr := GetManageGroupByCond(orgId, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcLangCode: db.In([]string{consts.ManageGroupSubNormalUser, consts.ManageGroupSubNormalUserBjx}),
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	if len(list) > 0 {
		return &list[0], nil
	}

	return nil, nil
}

// 将形如：`["Permission.Pro.Tag-Create,Modify,Delete"]` 的权限数组转换为 operation（格式是：`Permission.Pro.Tag.Create`） 一致的格式。
func TransferOperationArr(optAuthArr []string) []string {
	opList := make([]string, 0)
	for _, item := range optAuthArr {
		infos := strings.Split(item, "-")
		if len(infos) > 1 {
			opPrev := infos[0]
			opSuffixArr := strings.Split(infos[1], ",")
			for _, oneSuffix := range opSuffixArr {
				opList = append(opList, fmt.Sprintf("%s.%s", opPrev, oneSuffix))
			}
		} else {
			opList = append(opList, item)
		}
	}
	return slice.SliceUniqueString(opList)
}

// GetSearchedIndexArr 从 int64 数组中，查找一个数，返回所在索引。如果未找到，返回 -1
func GetSearchedIndexArr(list []int64, needle int64) int {
	for index, item := range list {
		if needle == item {
			return index
		}
	}
	return -1
}

func GetAdminGroupsByLangCode(orgId int64, langCodeList []string) ([]po.LcPerManageGroup, error) {
	list, dbErr := GetManageGroupByCond(orgId, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		// []string{consts.ManageGroupSubNormalUser, consts.ManageGroupSubNormalUserBjx}
		consts.TcLangCode: db.In(langCodeList),
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	return list, nil
}

// TransferUserIdsFromUserIdJson 将 userIds json 转换为数组
func TransferUserIdsFromUserIdJson(userIdsJson string) []int64 {
	list := make([]int64, 0)
	if userIdsJson == "" {
		return list
	}

	// err := jsonx.FromJson(userIdsJson, &list)// 这个方法在一些时候有问题
	// 下面这个方法不支持 ["1001", 1002] 这种 case
	//d := stdJson.NewDecoder(bytes.NewBuffer([]byte(userIdsJson)))
	//d.UseNumber()
	//err := d.Decode(&list)
	tmpArr := make([]interface{}, 0)
	err := jsonx.FromJson(userIdsJson, &tmpArr)
	if err != nil {
		logger.Error(err)
		return list
	}
	for _, item := range tmpArr {
		if val, ok1 := item.(string); ok1 {
			tmpUid, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				logger.ErrorF("TransferUserIdsFromUserIdJson err: %v", err)
				continue
			}
			list = append(list, tmpUid)
		} else if val, ok2 := item.(float64); ok2 {
			tmpUid := int64(val)
			list = append(list, tmpUid)
		} else if val, ok2 := item.(int64); ok2 {
			list = append(list, val)
		} else {
			logger.ErrorF("TransferUserIdsFromUserIdJson convert error: %v", item)
		}
	}

	return slice.SliceUniqueInt64(list)
}

// GetOneNormalAdminForUpgrade [bjx专用]获取普通管理员id，用于将其设为超管
func GetOneNormalAdminForUpgrade(orgId int64) ([]int64, error) {
	normalAdminIds := make([]int64, 0)
	// 查询普通管理员
	groups, dbErr := GetAdminGroupsByLangCode(orgId, []string{consts.ManageGroupSubNormalAdmin})
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	if len(groups) < 1 {
		return nil, errors.New("普通管理员不存在。")
	}
	for _, group := range groups {
		tmpUserIds := TransferUserIdsFromUserIdJson(group.UserIds)
		normalAdminIds = append(normalAdminIds, tmpUserIds...)
	}

	return normalAdminIds, nil
}

// GetOneNormalUserForSuperAdmin [bjx专用]获取一个普通用户，用于将其设为超管
func GetOneNormalUserForSuperAdmin(orgId int64) (int64, error) {
	oneEffectiveUid := int64(0)
	// 找出用户不在“超管组”、“普通管理组”中的用户
	groups, dbErr := GetAdminGroupsByLangCode(orgId, []string{consts.ManageGroupSys, consts.ManageGroupSubNormalAdmin})
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}
	exceptUids := make([]int64, 0)
	for _, group := range groups {
		tmpUids := TransferUserIdsFromUserIdJson(group.UserIds)
		exceptUids = append(exceptUids, tmpUids...)
	}
	// 查找一个普通用户
	cond1 := db.Cond{
		consts.TcOrgId:       orgId,
		consts.TcCheckStatus: consts.AppCheckStatusSuccess,
		//consts.TcUseStatus:   consts.AppStatusEnable,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	if len(exceptUids) > 0 {
		cond1[consts.TcUserId] = db.NotIn(exceptUids)
	}
	pos, dbErr := GetUserIdsOfOrgByCond(cond1, 1, 1)
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}
	for _, relation := range pos {
		oneEffectiveUid = relation.UserId
		if oneEffectiveUid > 0 {
			break
		}
	}
	return oneEffectiveUid, nil
}

// FilterExistTrashMenuRoles 过滤掉已经存在回收站的菜单项的角色，并返回最后的角色 Id 列表
func FilterExistTrashMenuRoles(roles []po.LcPerManageGroup) []int64 {
	filteredData := make([]int64, 0)
	for _, role := range roles {
		if strings.Contains(role.OptAuth, consts.MenuPermissionOrgTrash) {
			continue
		}
		filteredData = append(filteredData, role.Id)
	}

	return filteredData
}

// GetAppendPrmStrForSql 组装 sql，用于向角色中追加组织管理组的菜单权限项
func GetAppendPrmStrForSql(items []string) string {
	// 存放 `'$', "1001-MenuPermission.Org-Workspace"`
	frameArr := make([]string, 0)
	for _, code := range items {
		frameArr = append(frameArr, fmt.Sprintf("'$', \"%s\"", code))
	}
	opStr := strings.Join(frameArr, ", ")

	return opStr
}
