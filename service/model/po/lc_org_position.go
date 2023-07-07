package po

import "time"

/**
职级信息
@author WangShiChang
@version v1.0
@date 2020-10-21
*/

// LcOrgPosition 职级信息
type LcOrgPosition struct {
	Id            int64     `db:"id,omitempty" json:"id"` //注意此处的ID是全局ID
	OrgId         int64     `db:"org_id,omitempty" json:"orgId"`
	OrgPositionId int64     `db:"org_position_id,omitempty" json:"orgPositionId"` //注意此处是org内的职级ID
	Name          string    `db:"name,omitempty" json:"name"`
	PositionLevel int       `db:"position_level,omitempty" json:"positionLevel"`
	Remark        string    `db:"remark,omitempty" json:"remark"`
	Status        int       `db:"status,omitempty" json:"status"`
	Creator       int64     `db:"creator,omitempty" json:"creator"`
	CreateTime    time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator       int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime    time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version       int       `db:"version,omitempty" json:"version"`
	IsDelete      int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*LcOrgPosition) TableName() string {
	return "lc_org_position"
}

// UserPosition 用户职级信息
type UserPosition struct {
	Id            int64  `db:"id,omitempty" json:"id"`                         //注意此处的ID是全局ID
	OrgPositionId int64  `db:"org_position_id,omitempty" json:"orgPositionId"` //注意此处是org内的职级ID
	PositionName  string `db:"position_name,omitempty" json:"positionName"`
	PositionLevel int    `db:"position_level,omitempty" json:"positionLevel"`
	UserId        int64  `db:"user_id,omitempty" json:"userId"`
}
