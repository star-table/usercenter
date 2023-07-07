package bo

import (
	"sort"
	"time"
)

type PriorityBo struct {
	Id         int64     `db:"id,omitempty" json:"id"`
	OrgId      int64     `db:"org_id,omitempty" json:"orgId"`
	LangCode   string    `db:"lang_code,omitempty" json:"langCode"`
	Name       string    `db:"name,omitempty" json:"name"`
	Type       int       `db:"type,omitempty" json:"type"`
	Sort       int       `db:"sort,omitempty" json:"sort"`
	BgStyle    string    `db:"bg_style,omitempty" json:"bgStyle"`
	FontStyle  string    `db:"font_style,omitempty" json:"fontStyle"`
	IsDefault  int       `db:"is_default,omitempty" json:"isDefault"`
	Remark     string    `db:"remark,omitempty" json:"remark"`
	Status     int       `db:"status,omitempty" json:"status"`
	Creator    int64     `db:"creator,omitempty" json:"creator"`
	CreateTime time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator    int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version    int       `db:"version,omitempty" json:"version"`
	IsDelete   int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PriorityBo) TableName() string {
	return "ppm_prs_priority"
}

type By func(p1, p2 *PriorityBo) bool

func (by By) Sort(planets []PriorityBo) {
	ps := &PriorityBoSorter{
		priorityBos: planets,
		by:          by,
	}
	sort.Sort(ps)
}

type PriorityBoSorter struct {
	priorityBos []PriorityBo
	by          func(p1, p2 *PriorityBo) bool
}

func (s *PriorityBoSorter) Len() int {
	return len(s.priorityBos)
}

func (s *PriorityBoSorter) Swap(i, j int) {
	s.priorityBos[i], s.priorityBos[j] = s.priorityBos[j], s.priorityBos[i]
}

func (s *PriorityBoSorter) Less(i, j int) bool {
	return s.by(&s.priorityBos[i], &s.priorityBos[j])
}

//升序
func SortPriorityBo(prioritiesBos []PriorityBo) {
	s := func(p1, p2 *PriorityBo) bool {
		return p1.Sort < p2.Sort
	}
	By(s).Sort(prioritiesBos)
}
