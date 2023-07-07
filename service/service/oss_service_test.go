package service

import (
	"log"
	"os"
	"testing"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/service/model/req"
)

// 测试获取 oss 上传策略
func TestGetOssPostPolicy(t *testing.T) {
	env := os.Getenv(consts.RunEnv)
	log.Printf("%s, %v\n", env, env == "")
	panic("test")
	// todo config init
	var (
		orgId     = int64(1001)
		userId    = int64(1305473375526715392)
		projectID = int64(0)
		issueID   = int64(0)
		folderID  = int64(0)
	)
	var resp, err = GetOssPostPolicy(orgId, userId, req.GetOssPostPolicyReq{
		PolicyType: 6,
		ProjectID:  &projectID,
		IssueID:    &issueID,
		FolderID:   &folderID,
	})
	if err != nil {
		t.Error(resp)
		return
	}
	t.Log(resp)
}

func configInit() {

}
