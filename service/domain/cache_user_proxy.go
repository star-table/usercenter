package domain

import (
	"github.com/star-table/usercenter/core/consts"
	sconsts "github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/util"
	"github.com/star-table/usercenter/pkg/util/copyer"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/po"
	"upper.io/db.v3"
)

func GetBaseUserOutInfoBatch(sourceChannel string, orgId int64, userIds []int64) ([]bo.BaseUserOutInfoBo, errs.SystemErrorInfo) {
	keys := make([]interface{}, len(userIds))
	for i, userId := range userIds {
		key, keyErr := util.ParseCacheKey(sconsts.CacheBaseUserOutInfo, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:         orgId,
			consts.CacheKeySourceChannelConstName: sourceChannel,
			consts.CacheKeyUserIdConstName:        userId,
		})
		if keyErr != nil {
			logger.Error(keyErr)
			return nil, errs.TemplateRenderError
		}
		keys[i] = key
	}
	resultList := make([]string, 0)
	if len(keys) > 0 {
		list, redisErr := store.Redis.MGet(keys...)
		if redisErr != nil {
			logger.Error(redisErr)
			return nil, errs.RedisOperateError
		}
		resultList = list
	}
	baseUserOutInfoList := make([]bo.BaseUserOutInfoBo, 0)
	validUserIds := map[int64]bool{}
	for _, userInfoJson := range resultList {
		userOutInfoBo := &bo.BaseUserOutInfoBo{}
		jsonErr := json.FromJson(userInfoJson, userOutInfoBo)
		if jsonErr != nil {
			logger.Error(jsonErr)
			return nil, errs.JSONConvertError
		}
		baseUserOutInfoList = append(baseUserOutInfoList, *userOutInfoBo)
		validUserIds[userOutInfoBo.UserId] = true
	}

	missUserIds := make([]int64, 0)
	//找不存在的
	if len(userIds) != len(validUserIds) {
		for _, userId := range userIds {
			if _, ok := validUserIds[userId]; !ok {
				missUserIds = append(missUserIds, userId)
			}
		}
	}

	//批量查外部信息
	outInfos, err := GetBaseUserOutInfoByUserIds(sourceChannel, orgId, missUserIds)
	if err != nil {
		return nil, err
	}

	if len(outInfos) > 0 {
		baseUserOutInfoList = append(baseUserOutInfoList, outInfos...)
	}

	return baseUserOutInfoList, nil
}

func GetBaseUserInfoBatch(sourceChannel string, orgId int64, userIds []int64) ([]bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	//去重
	userIds = slice.SliceUniqueInt64(userIds)

	keys := make([]interface{}, len(userIds))
	for i, userId := range userIds {
		key, keyErr := util.ParseCacheKey(sconsts.CacheBaseUserInfo, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:  orgId,
			consts.CacheKeyUserIdConstName: userId,
		})
		if keyErr != nil {
			logger.Error(keyErr)
			return nil, errs.TemplateRenderError
		}
		keys[i] = key
	}
	resultList := make([]string, 0)
	if len(keys) > 0 {
		list, redisErr := store.Redis.MGet(keys...)
		if redisErr != nil {
			logger.Error(redisErr)
			return nil, errs.RedisOperateError
		}
		resultList = list
	}
	baseUserInfoList := make([]bo.BaseUserInfoBo, 0)
	validUserIds := map[int64]bool{}
	for _, userInfoJson := range resultList {
		userInfoBo := &bo.BaseUserInfoBo{}
		jsonErr := json.FromJson(userInfoJson, userInfoBo)
		if jsonErr != nil {
			logger.Error(jsonErr)
			return nil, errs.JSONConvertError
		}
		baseUserInfoList = append(baseUserInfoList, *userInfoBo)
		validUserIds[userInfoBo.UserId] = true
	}

	logger.InfoF("from cache %s", json.ToJsonIgnoreError(baseUserInfoList))
	missUserIds := make([]int64, 0)
	//找不存在的
	if len(userIds) != len(validUserIds) {
		for _, userId := range userIds {
			if _, ok := validUserIds[userId]; !ok {
				missUserIds = append(missUserIds, userId)
			}
		}
	}

	missUserInfos, userErr := getLocalBaseUserInfoBatch(orgId, missUserIds)
	if userErr != nil {
		logger.Error(userErr)
		return nil, userErr
	}
	if len(missUserInfos) > 0 {
		baseUserInfoList = append(baseUserInfoList, missUserInfos...)
	}

	if sourceChannel != "" {
		//获取用户外部信息
		baseUserOutInfos, err := GetBaseUserOutInfoBatch(sourceChannel, orgId, userIds)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		outInfoMap := make(map[int64]bo.BaseUserOutInfoBo)
		for _, outInfo := range baseUserOutInfos {
			outInfoMap[outInfo.UserId] = outInfo
		}
		for i, _ := range baseUserInfoList {
			userInfo := baseUserInfoList[i]
			if outInfo, ok := outInfoMap[userInfo.UserId]; ok {
				userInfo.OutUserId = outInfo.OutUserId
				userInfo.OutOrgUserId = outInfo.OutOrgUserId
				userInfo.OutOrgId = outInfo.OutOrgId
				userInfo.HasOutInfo = outInfo.OutUserId != ""
				userInfo.HasOrgOutInfo = outInfo.OutOrgId != ""
			}
			baseUserInfoList[i] = userInfo
		}
	}
	return baseUserInfoList, nil
}

func ClearBaseUserInfo(orgId, userId int64) errs.SystemErrorInfo {
	key, keyErr := util.ParseCacheKey(sconsts.CacheBaseUserInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return errs.TemplateRenderError
	}
	_, err := store.Redis.Del(key)
	if err != nil {
		logger.Error(err)
		return errs.RedisOperateError
	}
	return nil
}

//批量清楚用户缓存信息
func ClearBaseUserInfoBatch(orgId int64, userIds []int64) errs.SystemErrorInfo {
	keys := make([]interface{}, 0)
	for _, userId := range userIds {
		key, keyErr := util.ParseCacheKey(sconsts.CacheBaseUserInfo, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:  orgId,
			consts.CacheKeyUserIdConstName: userId,
		})
		if keyErr != nil {
			logger.Error(keyErr)
			return errs.TemplateRenderError
		}
		keys = append(keys, key)
	}
	_, redisErr := store.Redis.Del(keys...)
	if redisErr != nil {
		logger.Error(redisErr)
		return errs.RedisOperateError
	}
	return nil
}

//sourceChannel可以为空
func GetBaseUserInfo(sourceChannel string, orgId int64, userId int64) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	if userId == 0 {
		//系统创建
		return &bo.BaseUserInfoBo{
			OrgId: orgId,
			Name:  "系统",
		}, nil
	}

	key, keyErr := util.ParseCacheKey(sconsts.CacheBaseUserInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return nil, errs.TemplateRenderError
	}

	baseUserInfoJson, redisErr := store.Redis.Get(key)
	if redisErr != nil {
		logger.Error(redisErr)
		return nil, errs.RedisOperateError
	}
	baseUserInfo := &bo.BaseUserInfoBo{}
	if baseUserInfoJson != "" {
		jsonErr := json.FromJson(baseUserInfoJson, baseUserInfo)
		if jsonErr != nil {
			logger.Error(jsonErr)
			return nil, errs.JSONConvertError
		}
	} else {
		userInfo, err := getLocalBaseUserInfo(orgId, userId, key)
		if err != nil {
			return nil, err
		}
		baseUserInfo = userInfo
	}
	//这里不存缓存，动态获取
	baseUserOutInfo, err := GetBaseUserOutInfo(orgId, userId)
	if err != nil {
		return nil, err
	}
	baseUserInfo.OutUserId = baseUserOutInfo.OutUserId
	baseUserInfo.OutOrgId = baseUserOutInfo.OutOrgId
	baseUserInfo.HasOutInfo = baseUserInfo.OutUserId != ""
	baseUserInfo.HasOrgOutInfo = baseUserInfo.OutOrgId != ""
	baseUserInfo.OutOrgUserId = baseUserOutInfo.OutOrgUserId

	return baseUserInfo, nil
}

func GetBaseUserOutInfo(orgId int64, userId int64) (*bo.BaseUserOutInfoBo, errs.SystemErrorInfo) {
	if userId == 0 {
		//系统创建
		return &bo.BaseUserOutInfoBo{
			OrgId: orgId,
		}, nil
	}

	key, keyErr := util.ParseCacheKey(sconsts.CacheBaseUserOutInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return nil, errs.TemplateRenderError
	}
	baseUserOutInfoJson, redisErr := store.Redis.Get(key)
	if redisErr != nil {
		logger.Error(redisErr)
		return nil, errs.RedisOperateError
	}
	if baseUserOutInfoJson != "" {
		baseUserOutInfo := &bo.BaseUserOutInfoBo{}
		jsonErr := json.FromJson(baseUserOutInfoJson, baseUserOutInfo)
		if jsonErr != nil {
			logger.Error(jsonErr)
			return nil, errs.JSONConvertError
		}
		return baseUserOutInfo, nil
	} else {
		//用户外部信息
		userOutInfo := &po.PpmOrgUserOutInfo{}
		_ = store.Mysql.SelectOneByCond(consts.TableUserOutInfo, db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcOrgId:    orgId,
			consts.TcUserId:   userId,
		}, userOutInfo)

		outInfo := bo.BaseUserOutInfoBo{
			UserId:       userId,
			OrgId:        orgId,
			OutUserId:    userOutInfo.OutUserId,
			OutOrgUserId: userOutInfo.OutOrgUserId,
		}
		//组织外部信息
		orgOutInfo := &po.PpmOrgOrganizationOutInfo{}
		err := store.Mysql.SelectOneByCond(consts.TableOrganizationOutInfo, db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcOrgId:    orgId,
		}, orgOutInfo)
		if err != nil {
			if err == db.ErrNoMoreRows {

			} else {
				logger.Error(err)
				return nil, errs.MysqlOperateError
			}
		} else {
			outInfo.OutOrgId = orgOutInfo.OutOrgId
		}
		//baseUserOutInfoJson := json.ToJsonIgnoreError(outInfo)
		//
		//redisErr = store.Redis.SetEx(key, baseUserOutInfoJson, consts.GetCacheBaseExpire())
		//if redisErr != nil {
		//	logger.Error(redisErr)
		//	return nil, errs.RedisOperateError
		//}
		return &outInfo, nil
	}
}

func GetBaseUserOutInfoByUserIds(sourceChannel string, orgId int64, userIds []int64) ([]bo.BaseUserOutInfoBo, errs.SystemErrorInfo) {
	logger.InfoF("批量获取用户外部信息 %d, %s", orgId, json.ToJsonIgnoreError(userIds))

	resultList := make([]bo.BaseUserOutInfoBo, 0)

	if userIds == nil || len(userIds) == 0 {
		return resultList, nil
	}

	//用户外部信息
	var userOutInfos []po.PpmOrgUserOutInfo
	dbErr := store.Mysql.SelectAllByCond(consts.TableUserOutInfo, db.Cond{
		consts.TcIsDelete:      consts.AppIsNoDelete,
		consts.TcOrgId:         orgId,
		consts.TcSourceChannel: sourceChannel,
		consts.TcUserId:        db.In(userIds),
	}, &userOutInfos)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}

	//组织外部信息
	var orgOutInfo po.PpmOrgOrganizationOutInfo
	dbErr = store.Mysql.SelectOneByCond(consts.TableOrganizationOutInfo, db.Cond{
		consts.TcIsDelete:      consts.AppIsNoDelete,
		consts.TcOrgId:         orgId,
		consts.TcSourceChannel: sourceChannel,
	}, &orgOutInfo)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return nil, errs.OrgOutInfoNotExist
		}
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}
	msetArgs := map[string]string{}
	keys := make([]string, 0)
	for _, userOutInfo := range userOutInfos {
		key, keyErr := util.ParseCacheKey(sconsts.CacheBaseUserOutInfo, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:         orgId,
			consts.CacheKeySourceChannelConstName: sourceChannel,
			consts.CacheKeyUserIdConstName:        userOutInfo.UserId,
		})
		if keyErr != nil {
			logger.Error(keyErr)
			return nil, errs.TemplateRenderError
		}
		keys = append(keys, key)

		outInfo := bo.BaseUserOutInfoBo{
			UserId:       userOutInfo.UserId,
			OrgId:        orgId,
			OutUserId:    userOutInfo.OutUserId,
			OutOrgId:     orgOutInfo.OutOrgId,
			OutOrgUserId: userOutInfo.OutOrgUserId,
		}

		resultList = append(resultList, outInfo)
		msetArgs[key] = json.ToJsonIgnoreError(outInfo)
	}

	if len(msetArgs) > 0 {
		redisErr := store.Redis.MSet(msetArgs)
		if redisErr != nil {
			logger.Error(redisErr)
			return nil, errs.RedisOperateError
		}
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.ErrorF("捕获到的错误：%s", r)
			}
		}()

		//for _, key := range keys {
		//	_, _ = store.Redis.Expire(key, consts.GetCacheBaseUserInfoExpire())
		//}
	}()

	return resultList, nil
}

// getLocalBaseUserInfo sourceChannel可以为空
func getLocalBaseUserInfo(orgId, userId int64, key string) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	user := &po.PpmOrgUser{}
	dbErr := store.Mysql.SelectById(user.TableName(), userId, user)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return nil, errs.UserNotExist
		}
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}

	baseUserInfo := &bo.BaseUserInfoBo{
		UserId:             user.Id,
		Name:               user.Name,
		NamePy:             user.NamePinyin,
		Avatar:             user.Avatar,
		OrgId:              orgId,
		OrgUserIsDelete:    2,
		OrgUserStatus:      1,
		OrgUserCheckStatus: 1,
	}

	if orgId > 0 {
		newestUserOrganization, err1 := GetUserOrganizationNewestRelation(orgId, userId)
		if err1 != nil {
			logger.Error(err1)
			return nil, err1
		}
		baseUserInfo.OrgUserIsDelete = newestUserOrganization.IsDelete
		baseUserInfo.OrgUserStatus = newestUserOrganization.Status
		baseUserInfo.OrgUserCheckStatus = newestUserOrganization.CheckStatus
	}

	return baseUserInfo, nil
}

//用来获取用户最新的组织关系
func GetUserOrganizationNewestRelation(orgId, userId int64) (*po.PpmOrgUserOrganization, errs.SystemErrorInfo) {
	UserOrganizationPo := &po.PpmOrgUserOrganization{}

	conn, err := store.Mysql.GetConnect()
	if err != nil {
		return nil, errs.MysqlOperateError
	}
	err = conn.Collection(consts.TableUserOrganization).Find(db.Cond{
		consts.TcOrgId:  orgId,
		consts.TcUserId: userId,
	}).OrderBy("id desc").Limit(1).One(UserOrganizationPo)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, errs.BuildSystemErrorInfo(errs.UserOrgNotRelation)
		} else {
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		}
	}
	return UserOrganizationPo, nil
}

func getLocalBaseUserInfoBatch(orgId int64, userIds []int64) ([]bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	logger.InfoF("批量获取用户信息 %d, %s", orgId, json.ToJsonIgnoreError(userIds))

	baseUserInfos := make([]bo.BaseUserInfoBo, 0)

	if userIds == nil || len(userIds) == 0 {
		return baseUserInfos, nil
	}

	var users []po.PpmOrgUser
	dbErr := store.Mysql.SelectAllByCond(consts.TableUser, db.Cond{
		consts.TcId: db.In(userIds),
	}, &users)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}

	//获取关联列表，要做去重
	var userOrganizationPos []po.PpmOrgUserOrganization
	_, dbErr = store.Mysql.SelectAllByCondWithPageAndOrder(consts.TableUserOrganization, db.Cond{
		consts.TcOrgId:  orgId,
		consts.TcUserId: db.In(userIds),
	}, nil, 0, -1, "id asc", &userOrganizationPos)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}

	//id升序，保留最新: key: userId, value: po
	userOrgMap := map[int64]po.PpmOrgUserOrganization{}
	for _, userOrg := range userOrganizationPos {
		userOrgMap[userOrg.UserId] = userOrg
	}

	for _, user := range users {
		baseUserInfo := bo.BaseUserInfoBo{
			UserId: user.Id,
			Name:   user.Name,
			NamePy: user.NamePinyin,
			Avatar: user.Avatar,
			OrgId:  orgId,
		}

		if userOrg, ok := userOrgMap[user.Id]; ok {
			baseUserInfo.OrgUserIsDelete = userOrg.IsDelete
			baseUserInfo.OrgUserStatus = userOrg.Status
			baseUserInfo.OrgUserCheckStatus = userOrg.CheckStatus
		}

		baseUserInfos = append(baseUserInfos, baseUserInfo)
	}

	msetArgs := map[string]string{}
	keys := make([]string, 0)
	for _, baseUserInfo := range baseUserInfos {
		key, keyErr := util.ParseCacheKey(sconsts.CacheBaseUserInfo, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:  orgId,
			consts.CacheKeyUserIdConstName: baseUserInfo.UserId,
		})
		if keyErr != nil {
			logger.Error(keyErr)
			return nil, errs.TemplateRenderError
		}
		msetArgs[key] = json.ToJsonIgnoreError(baseUserInfo)
		keys = append(keys, key)
	}

	if len(msetArgs) > 0 {
		redisErr := store.Redis.MSet(msetArgs)
		if redisErr != nil {
			logger.Error(redisErr)
			return nil, errs.RedisOperateError
		}
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.ErrorF("捕获到的错误：%s", r)
			}
		}()

		for _, key := range keys {
			_, _ = store.Redis.Expire(key, consts.GetCacheBaseUserInfoExpire())
		}
	}()
	return baseUserInfos, nil
}

func GetUserConfigInfo(orgId int64, userId int64) (*bo.UserConfigBo, errs.SystemErrorInfo) {
	key, keyErr := util.ParseCacheKey(sconsts.CacheUserConfig, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return nil, errs.TemplateRenderError
	}

	userConfigJson, redisErr := store.Redis.Get(key)
	if redisErr != nil {
		logger.Error(redisErr)
		return nil, errs.RedisOperateError
	}
	userConfigBo := &bo.UserConfigBo{}
	if userConfigJson != "" {
		jsonErr := json.FromJson(userConfigJson, userConfigBo)
		if jsonErr != nil {
			logger.Error(jsonErr)
			return nil, errs.JSONConvertError
		}
		return userConfigBo, nil
	} else {
		userConfig := &po.PpmOrgUserConfig{}
		dbErr := store.Mysql.SelectOneByCond(userConfig.TableName(), db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcUserId:   userId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, userConfig)
		if dbErr != nil {
			logger.Error(dbErr)
			return nil, errs.MysqlOperateError
		}
		_ = copyer.Copy(userConfig, userConfigBo)
		userConfigJson = json.ToJsonIgnoreError(userConfigBo)

		redisErr := store.Redis.Set(key, userConfigJson)
		if redisErr != nil {
			logger.Error(redisErr)
			return nil, errs.RedisOperateError
		}
		return userConfigBo, nil
	}
}

func DeleteUserConfigInfo(orgId int64, userId int64) errs.SystemErrorInfo {
	key, keyErr := util.ParseCacheKey(sconsts.CacheUserConfig, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return errs.TemplateRenderError
	}
	_, redisErr := store.Redis.Del(key)
	if redisErr != nil {
		logger.Error(redisErr)
		return errs.RedisOperateError
	}
	return nil
}

func ClearUserCacheInfo(token string) errs.SystemErrorInfo {
	userCacheKey := sconsts.CacheUserToken + token
	_, redisErr := store.Redis.Del(userCacheKey)
	if redisErr != nil {
		logger.Error(redisErr)
		return errs.RedisOperateError
	}
	return nil
}

func SetDingCodeCache(outUserId string) errs.SystemErrorInfo {
	key, keyErr := util.ParseCacheKey(sconsts.LoginByDingCode, map[string]interface{}{
		consts.CacheKeySourceChannelConstName: consts.AppSourceChannelDingTalk,
		consts.CacheKeyOutUserIdConstName:     outUserId,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return errs.TemplateRenderError
	}

	redisErr := store.Redis.SetEx(key, outUserId, int64(60*5))
	if redisErr != nil {
		logger.Error(redisErr)
		return errs.RedisOperateError
	}

	return nil
}

func GetDingCodeCache(outUserId string) errs.SystemErrorInfo {
	key, keyErr := util.ParseCacheKey(sconsts.LoginByDingCode, map[string]interface{}{
		consts.CacheKeySourceChannelConstName: consts.AppSourceChannelDingTalk,
		consts.CacheKeyOutUserIdConstName:     outUserId,
	})
	if keyErr != nil {
		logger.Error(keyErr)
		return errs.TemplateRenderError
	}

	cacheJson, redisErr := store.Redis.Get(key)
	if redisErr != nil {
		logger.Error(redisErr)
		return errs.RedisOperateError
	}

	if cacheJson == consts.BlankString {
		return errs.DingCodeCacheInvalid
	}

	//查到之后清掉
	_, redisErr = store.Redis.Del(key)
	if redisErr != nil {
		logger.Error(redisErr)
		return errs.RedisOperateError
	}

	return nil
}
