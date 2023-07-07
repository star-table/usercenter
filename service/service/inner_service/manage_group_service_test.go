package inner_service

import "testing"

func TestGetManager(t *testing.T) {
	manager, err := GetManager(2373)
	if err != nil {
		t.Error(err)
	}
	t.Log(manager.Data)
}

func TestGetUsersCouldManage(t *testing.T)  {
	resp, err := GetUsersCouldManage(2373, -1)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp.List)
}

func TestGetUserAuthorityByUserIdSimple(t *testing.T)  {
	resp, err := GetUserAuthorityByUserIdSimple(2373, 29612)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}

func TestGetCommAdminMangeApps(t *testing.T)  {
	mangeApps, err := GetCommAdminMangeApps(2373)
	if err != nil {
		t.Error(err)
	}
	t.Log(mangeApps)
}

func TestGetUserAuthorityByUserId(t *testing.T)  {
	resp, err := GetUserAuthorityByUserId(2758, 34067)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}

func TestGetUserDeptIdsWithParentId(t *testing.T)  {
	resp, err := GetUserDeptIdsWithParentId(2797, 34135)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}

func TestGetUserIdsByDeptIds(t *testing.T)  {
	deptIds := []int64{96750053334800, 97057345051152}
	resp, err := GetUserIdsByDeptIds(2797, deptIds)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}
