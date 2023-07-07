package domain

import (
	"strconv"
	"strings"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/nacos"
	"github.com/star-table/usercenter/core/snowflake"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/store/mysql"
	"github.com/star-table/usercenter/pkg/util/asyn"
	"github.com/star-table/usercenter/pkg/util/copyer"
	"github.com/star-table/usercenter/pkg/util/format"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/pkg/util/uuid"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/po"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/resp"
	"github.com/star-table/usercenter/service/model/resp/inner_resp"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

// GetOrgListByIds 根据组织ID列表获取组织信息
func GetOrgListByIds(orgIds []int64) ([]po.PpmOrgOrganization, error) {
	var organizationList []po.PpmOrgOrganization
	dbErr := store.Mysql.SelectAllByCond(consts.TableOrganization, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcStatus:   consts.AppStatusEnable,
		consts.TcId:       db.In(orgIds),
	}, &organizationList)

	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	return organizationList, nil
}

// GetOrgById
func GetOrgById(orgId int64) (*po.PpmOrgOrganization, error) {
	var org po.PpmOrgOrganization
	dbErr := store.Mysql.SelectOneByCond(consts.TableOrganization, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcStatus:   consts.AppStatusEnable,
		consts.TcId:       orgId,
	}, &org)

	if dbErr != nil {
		return nil, dbErr
	}
	return &org, nil
}

// GetOrgByCode
func GetOrgByCode(code string) (*po.PpmOrgOrganization, error) {
	var org po.PpmOrgOrganization
	dbErr := store.Mysql.SelectOneByCond(consts.TableOrganization, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcCode:     code,
	}, &org)

	if dbErr != nil {
		return nil, dbErr
	}
	return &org, nil
}

// GetOrgByOutOrgId
func GetOrgByOutOrgId(sourceChannel, outOrgId string) (*po.PpmOrgOrganization, error) {
	conn, dbErr := store.Mysql.GetConnect()
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	var orgPo po.PpmOrgOrganization
	dbErr = conn.Select("o.*").From(consts.TableOrganization + " o").
		Join(consts.TableOrganizationOutInfo + " ot").
		On(db.Cond{
			"o." + consts.TcOrgId:         db.Raw("ot." + consts.TcOrgId),
			"o." + consts.TcIsDelete:      consts.AppIsNoDelete,
			"o." + consts.TcSourceChannel: sourceChannel,
			"o." + consts.TcOutOrgId:      outOrgId,
		}).Limit(1).One(&orgPo)
	if dbErr != nil {
		return nil, dbErr
	}
	return &orgPo, nil
}

// CreateOrg 创建组织
func CreateOrg(createOrgBo bo.CreateOrgBo, creatorId int64, sourceChannel, sourcePlatform string, outOrgId string) (int64, errs.SystemErrorInfo) {
	orgId := snowflake.Id()
	orgOutId := snowflake.Id()
	//组织配置
	orgConfigId := snowflake.Id()

	orgName := strings.TrimSpace(createOrgBo.OrgName)
	isOrgNameRight := format.VerifyOrgNameFormat(orgName)
	if !isOrgNameRight {
		return 0, errs.OrgNameLenError
	}

	//组织
	org := &po.PpmOrgOrganization{
		Id:             orgId,
		Status:         consts.AppStatusEnable,
		IsDelete:       consts.AppIsNoDelete,
		Creator:        creatorId,
		Owner:          creatorId,
		Updator:        creatorId,
		Name:           orgName,
		SourceChannel:  sourceChannel,
		SourcePlatform: sourcePlatform,
	}

	if createOrgBo.Scale != nil {
		org.Scale = *createOrgBo.Scale
	}
	if createOrgBo.IndustryId != nil && *createOrgBo.IndustryId != 0 {
		_, err := GetIndustryBo(*createOrgBo.IndustryId)
		if err != nil {
			logger.Error(err)
			return 0, err
		}
		org.IndustryId = *createOrgBo.IndustryId
	}

	dbErr := store.Mysql.TransX(func(tx sqlbuilder.Tx) error {
		dbErr := store.Mysql.TransInsert(tx, org)
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}

		//外部组织信息
		orgOutInfo := &po.PpmOrgOrganizationOutInfo{
			Id:             orgOutId,
			OrgId:          orgId,
			OutOrgId:       outOrgId,
			SourceChannel:  sourceChannel,
			SourcePlatform: sourcePlatform,
			Name:           createOrgBo.OrgName,
			Creator:        creatorId,
			Updator:        creatorId,
		}
		dbErr = store.Mysql.TransInsert(tx, orgOutInfo)
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}

		//组织配置信息
		orgConfig := &po.PpmOrcConfig{
			Id:    orgConfigId,
			OrgId: orgId,
		}
		dbErr = store.Mysql.TransInsert(tx, orgConfig)
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}

		//创建系统管理组
		_, dbErr = CreateManageGroup(orgId, creatorId, req.CreateManageGroup{
			GroupType: 1,
		}, tx)
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}

		// 初始化根部门
		_, dbErr = InitOrgRootDept(orgId, creatorId, orgName, sourceChannel, tx)
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}

		// 初始化职级列表
		dbErr = InitPosition(orgId, creatorId)
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}

		return nil
	})

	if dbErr != nil {
		logger.Error(dbErr)
		return 0, errs.MysqlOperateError
	}

	asyn.Execute(func() {
		//配置数据源
		dataCenterData, err := GetDatabase(orgId)
		if err != nil {
			logger.Error(err)
			return
		}

		_, dbErr := store.Mysql.UpdateSmartWithCond(consts.TableOrgConfig, db.Cond{
			consts.TcId: orgConfigId,
		}, mysql.Upd{
			consts.TcDbId: dataCenterData.DbId,
			consts.TcDcId: dataCenterData.DcId,
			consts.TcDsId: dataCenterData.DsId,
		})
		if dbErr != nil {
			logger.Error(dbErr)
			return
		}
	})
	return orgId, nil
}

func GetDatabase(orgId int64) (*resp.AllocateData, errs.SystemErrorInfo) {
	//获取数据源信息
	msg, code, err := nacos.DoGet("datacenter", "datacenter/inner/api/v1/dc/allot", map[string]interface{}{"orgId": strconv.FormatInt(orgId, 10)})
	if err != nil {
		logger.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.SystemError, err)
	}
	if code != 200 {
		logger.Error("数据中心接口请求错误")
		return nil, errs.BuildSystemErrorInfo(errs.SystemError)
	}

	msgResp := &resp.AllocateDbResp{}
	jsonErr := json.FromJson(msg, msgResp)
	if jsonErr != nil {
		logger.Error(jsonErr)
		return nil, errs.JSONConvertError
	}

	if msgResp.Code != 0 {
		logger.Error(msgResp.Message)
		return nil, errs.SystemError
	}

	return &msgResp.Data, nil
}

func UpdateOrg(updateBo bo.UpdateOrganizationBo) errs.SystemErrorInfo {

	organizationBo := updateBo.Bo
	upds := updateBo.OrganizationUpdateCond

	_, err := store.Mysql.UpdateSmartWithCond(consts.TableOrganization, db.Cond{
		consts.TcId: organizationBo.Id,
	}, upds)

	if err != nil {
		logger.ErrorF("store.Mysql.TransUpdateSmart: %q\n", err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	err = ClearCacheBaseOrgInfo(organizationBo.SourceChannel, organizationBo.Id)

	if err != nil {
		logger.ErrorF("redis err: %q\n", err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}

	return nil
}

// GetOrgConfig 获取组织配置信息
func GetOrgConfig(orgId int64) (*resp.OrgConfig, errs.SystemErrorInfo) {
	var orgConfigPo po.PpmOrcConfig
	dbErr := store.Mysql.SelectOneByCond(consts.TableOrgConfig, db.Cond{
		consts.TcOrgId: orgId,
	}, &orgConfigPo)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return nil, errs.OrgConfigNotExist
		}
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}
	var orgConfigBo resp.OrgConfig
	_ = copyer.Copy(orgConfigPo, &orgConfigBo)
	return &orgConfigBo, nil
}

// GenerateAndSetApiKey 生成并设置ApiKey
func GenerateAndSetApiKey(orgId, operatorUid int64) error {
	logger.InfoF("[生成OpenApiKey] -> orgId:%d , operatorUid:%d", orgId, operatorUid)
	apiKey := generateApiKey()

	_, dbErr := store.Mysql.UpdateSmartWithCond(consts.TableOrganization, db.Cond{
		consts.TcId:       orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcApiKey:  apiKey,
		consts.TcUpdator: operatorUid,
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return dbErr
	}
	return nil
}

// RemoveApiKey 删除ApiKey
func RemoveApiKey(orgId, operatorUid int64) error {
	logger.InfoF("[删除OpenApiKey] -> orgId:%d , operatorUid:%d", orgId, operatorUid)

	_, dbErr := store.Mysql.UpdateSmartWithCond(consts.TableOrganization, db.Cond{
		consts.TcId:       orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcApiKey:  "",
		consts.TcUpdator: operatorUid,
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return dbErr
	}
	return nil
}

// generateApiKey 生成ApiKey
func generateApiKey() string {
	return strings.ReplaceAll(uuid.NewUuid(), "-", "")
}

// GetOrgByApiKey  根据ApiKey获取机构信息
func GetOrgByApiKey(apiKey string) (*po.PpmOrgOrganization, error) {
	var org po.PpmOrgOrganization
	dbErr := store.Mysql.SelectOneByCond(consts.TableOrganization, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcStatus:   consts.AppStatusEnable,
		consts.TcApiKey:   apiKey,
	}, &org)

	if dbErr != nil {
		return nil, dbErr
	}
	return &org, nil
}

// CheckSourceChannelIsPolaris 校验组织的 source_channel，是否是 `polaris`
func CheckSourceChannelIsPolaris(sourceChannel string) bool {
	return inner_resp.CheckSourceChannelIsPolaris(sourceChannel)
}

// AddOrgOutCollaborator 新增组织外部协作人
func AddOrgOutCollaborator(orgId, userId int64) errs.SystemErrorInfo {
	userOrg := po.PpmOrgUserOrganization{}
	err := store.Mysql.SelectOneByCond(consts.TableUserOrganization, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcUserId:   userId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, &userOrg)
	if err != nil {
		if err == db.ErrNoMoreRows {
			userOrg = po.PpmOrgUserOrganization{
				Id:          snowflake.Id(),
				OrgId:       orgId,
				UserId:      userId,
				CheckStatus: 2,
				UseStatus:   2,
				Status:      1,
				Type:        consts.OrgUserTypeOut,
			}
			err = store.Mysql.Insert(&userOrg)
			if err != nil {
				logger.Error(err)
				return errs.MysqlOperateError
			}
			return nil
		}
		logger.Error(err)
		return errs.MysqlOperateError
	}
	return nil
}
