package service

import (
	"testing"

	"github.com/star-table/usercenter/pkg/util/json"

	"github.com/star-table/usercenter/service/model/req"
)

func TestGetManageGroupDetail(t *testing.T) {
	groupDetail, err := GetManageGroupDetail(2373, 29612, 91513448486416)
	if err != nil {
		t.Error(err)
	}
	t.Log(groupDetail)
}

func TestUpdateManageGroupContents(t *testing.T) {
	data := `["Permission.Org.Config-Modify", "Permission.Org.Config-ModifyField", "Permission.Org.Config-TplSaveAs", "Permission.Org.Config-TplDelete", "Permission.Org.User-ModifyStatus", "Permission.Org.User-ModifyUserAdminGroup", "Permission.Org.User-ModifyUserDept", "Permission.Org.Department-Create", "Permission.Org.Department-Modify", "Permission.Org.Department-Delete", "Permission.Org.InviteUser-Invite", "Permission.Org.AdminGroup-View", "Permission.Org.AdminGroup-Create", "Permission.Org.AdminGroup-Modify", "Permission.Org.AdminGroup-Delete", "Permission.Org.Project-Create", "MenuPermission.Org-Workspace", "MenuPermission.Org-Issue", "MenuPermission.Org-Project", "MenuPermission.Org-PolarisTpl", "MenuPermission.Org-Member", "MenuPermission.Org-Trend", "MenuPermission.Org-WorkHour", "MenuPermission.Org-Trash", "MenuPermission.Org-Setting"]`
	opth := []string{}
	json.FromJson(data, &opth)
	UpdateManageGroupContents(2797, 29612, true, 96663534227216, req.UpdateManageGroupContents{
		Id:         96663534227216,
		Values:     nil,
		ValueIf:    []string{"-1"},
		Key:        "app_ids",
		SourceFrom: "",
		AuthToken:  "",
	})
}
