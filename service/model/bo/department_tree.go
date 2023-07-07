package bo

type DeptNode struct {
	ID       int64                 `json:"id"`
	Name     string                `json:"name"`
	ParentID int64                 `json:"parentId"`
	Childs   []*DepartmentTreeNode `json:"childs"`
	Parent   *DepartmentTreeNode   `json:"parent"`
}

type DeptTree struct {
	Root *DeptNode
}
