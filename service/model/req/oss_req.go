package req

// 获取文件上传策略信息请求参数
type GetOssPostPolicyReq struct {
	// 策略类型, 1: 项目封面，2：任务资源（需要callback）, 3：导入任务的excel，4：项目文件（需要callback），5：兼容测试，6:用户头像，9:备忘录
	PolicyType int `json:"policyType"`
	// 如果policyType为1和2和3，那么projectId必传(创建场景传0)。没有则传 0。
	ProjectID *int64 `json:"projectId"`
	// 如果policyType为2，那么issueId必传。没有则传 0。
	IssueID *int64 `json:"issueId"`
	// 目录id, policy为4的时候必填。没有则传 0。
	FolderID *int64 `json:"folderId"`
}
