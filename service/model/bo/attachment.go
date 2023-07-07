package bo

import "github.com/star-table/usercenter/service/model/vo"

type AttachmentBo struct {
	vo.Resource
	IssueList []IssueBo
}
