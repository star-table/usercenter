package service

import (
	"time"

	"github.com/star-table/usercenter/core/conf"
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/pkg/util"
	"github.com/star-table/usercenter/pkg/util/copyer"
	"github.com/star-table/usercenter/pkg/util/oss"
	"github.com/star-table/usercenter/pkg/util/random"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/resp"
)

// 获取 oss 上传的策略信息
func GetOssPostPolicy(orgId, userId int64, req req.GetOssPostPolicyReq) (*resp.GetOssPostPolicyResp, errs.SystemErrorInfo) {
	oc := conf.Cfg.OSS
	if oc == nil {
		logger.Error("oss缺少配置")
		return nil, errs.OssConfigError
	}
	var policyConfig *conf.OSSPolicyInfo = nil
	// 策略类型，文件上传所用场景
	policyType := req.PolicyType

	toDay := time.Now()
	year, month, day := toDay.Date()

	//文件名
	fileName := random.RandomFileName()
	// 默认一个文件大小限制
	fileSizeLimit := conf.Cfg.OSS.Policies.UserAvatar.MaxFileSize
	//定义callback
	callbackJson := ""
	switch policyType {
	case consts.OssPolicyTypeUserAvatar: // 用户头像
		c := conf.Cfg.OSS.Policies.UserAvatar
		c.Dir, _ = util.ParseCacheKey(c.Dir, map[string]interface{}{
			consts.CacheKeyOrgIdConstName: orgId,
			consts.CacheKeyYearConstName:  year,
			consts.CacheKeyMonthConstName: int(month),
			consts.CacheKeyDayConstName:   day,
		})
		fileSizeLimit = c.MaxFileSize
		policyConfig = &c
	case consts.OssPolicyTypeUserMemo: // 备忘录
		folderName := "user_memo"
		c := conf.Cfg.OSS.Policies.MixResource
		c.Dir, _ = util.ParseCacheKey(c.Dir, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:      orgId,
			consts.CacheKeyFolderNameConstName: folderName,
			consts.CacheKeyYearConstName:       year,
			consts.CacheKeyMonthConstName:      int(month),
			consts.CacheKeyDayConstName:        day,
		})
		fileSizeLimit = c.MaxFileSize
		policyConfig = &c
	}

	if policyConfig == nil {
		return nil, errs.BuildSystemErrorInfo(errs.OssPolicyTypeError)
	}
	policyConfig.MaxFileSize = fileSizeLimit

	//policyBo := oss.PostPolicyWithCallback(policyConfig.Dir, policyConfig.Expire, policyConfig.MaxFileSize, callbackJson)
	policyBo := oss.PostPolicy(policyConfig.Dir, policyConfig.Expire, policyConfig.MaxFileSize)

	//oss绑定自定义域名，所以覆盖host
	policyBo.Host = oc.EndPoint

	resp := &resp.GetOssPostPolicyResp{}
	copyErr := copyer.Copy(policyBo, resp)
	if copyErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	resp.FileName = fileName
	resp.Callback = callbackJson
	return resp, nil
}
