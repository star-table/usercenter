package domain

import (
	"time"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/snowflake"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/po"
	"github.com/star-table/usercenter/service/model/req"
	"upper.io/db.v3"
)

// AddLoginRecord 添加登录记录
func AddLoginRecord(orgId, userId int64, sourceChannel string, clientId string, userAgent string, msg string) (int64, error) {
	id := snowflake.Id()
	userLoginRecord := po.PpmOrgUserLoginRecord{
		Id:            id,
		OrgId:         orgId,
		UserId:        userId,
		LoginTime:     time.Now(),
		LoginIp:       clientId,
		UserAgent:     userAgent,
		Msg:           msg,
		SourceChannel: sourceChannel,
		Creator:       userId,
		Updator:       userId,
	}

	dbErr := store.Mysql.Insert(&userLoginRecord)
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}

	return id, nil
}

// GetLoginRecordList 获取登录记录
func GetLoginRecordList(orgId int64, reqParam req.LoginRecordListReq) (int64, []bo.LoginRecordBo, error) {
	conn, dbErr := store.Mysql.GetConnect()
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, nil, dbErr
	}
	cond := db.Cond{
		"l." + consts.TcOrgId:    orgId,
		"l." + consts.TcIsDelete: consts.AppIsNoDelete,
	}
	if !reqParam.StartDate.IsZero() && !reqParam.EndDate.IsZero() {
		cond["l."+consts.TcCreateTime] = db.Between(reqParam.StartDate, reqParam.EndDate)
	}
	if reqParam.IP != "" {
		cond["l."+consts.TcLoginIp] = db.Like(reqParam.IP + "%")
	}
	union := &db.Union{}
	if reqParam.AccountNo != "" {
		union = union.
			Or(db.Cond{
				"u." + consts.TcLoginName: db.Like(reqParam.AccountNo + "%"),
			}).
			Or(db.Cond{
				"u." + consts.TcMobile: db.Like(reqParam.AccountNo + "%"),
			}).
			Or(db.Cond{
				"u." + consts.TcEmail: db.Like(reqParam.AccountNo + "%"),
			})
	}

	mid := conn.Select(
		"l.*",
		"u.login_name",
		"u.mobile",
		"u.email",
	).
		From(consts.TableLoginRecord+" as l ").
		Join(consts.TableUser+" as u ").
		On(db.Cond{
			"l." + consts.TcUserId: db.Raw("u." + consts.TcId),
			"l." + consts.TcOrgId:  db.Raw("u." + consts.TcOrgId),
		}, union).Where(cond).OrderBy("l.create_time desc")
	if union.Sentences() != nil {
		mid.And(union)
	}
	count := uint64(0)
	// 不分页
	if reqParam.PageBo == nil {
		var list []bo.LoginRecordBo
		dbErr = mid.All(&list)
		if dbErr != nil {
			logger.Error(dbErr)
			return 0, nil, dbErr
		}
		return int64(len(list)), list, nil
	}

	paginate := mid.Paginate(uint(reqParam.Size)).Page(uint(reqParam.Page))
	count, dbErr = paginate.TotalEntries()
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, nil, dbErr
	}
	var list []bo.LoginRecordBo
	dbErr = paginate.All(&list)
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, nil, dbErr
	}
	return int64(count), list, nil
}
