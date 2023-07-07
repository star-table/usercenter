package bo

type PageBo struct {
	Page  int    `json:"page"`
	Size  int    `json:"size"`
	Order string `json:"order"`
}
