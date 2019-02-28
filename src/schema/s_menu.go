package schema

import (
	"strings"
	"time"
)

// Menu 菜单对象
type Menu struct {
	RecordID   string    `json:"record_id" swaggo:"false,记录ID"`
	Code       string    `json:"code" binding:"required" swaggo:"true,菜单编号"`
	Name       string    `json:"name" binding:"required" swaggo:"true,菜单名称"`
	Type       int       `json:"type" binding:"required,max=3,min=1" swaggo:"true,菜单类型(1：模块 2：功能 3：资源)"`
	Sequence   int       `json:"sequence" swaggo:"false,排序值"`
	Icon       string    `json:"icon" swaggo:"false,菜单图标"`
	Path       string    `json:"path" swaggo:"false,访问路径"`
	Method     string    `json:"method" swaggo:"false,资源请求方式"`
	ParentID   string    `json:"parent_id" swaggo:"false,父级内码"`
	ParentPath string    `json:"parent_path" swaggo:"false,父级路径"`
	Creator    string    `json:"creator" swaggo:"false,创建者"`
	CreatedAt  time.Time `json:"created_at" swaggo:"false,创建时间"`
}

// MenuQueryParam 查询条件
type MenuQueryParam struct {
	RecordIDs  []string // 记录ID列表
	Code       string   // 菜单编号(模糊查询)
	Name       string   // 菜单名称(模糊查询)
	Types      []int    // 菜单类型(1：模块 2：功能 3：资源)
	ParentID   *string  // 父级内码
	UserID     string   // 用户ID（查询用户所拥有的菜单权限）
	ParentPath string   // 父级路径(前缀模糊查询)
}

// MenuQueryOptions 查询可选参数项
type MenuQueryOptions struct {
	PageParam *PaginationParam // 分页参数
}

// MenuQueryResult 查询结果
type MenuQueryResult struct {
	Data       Menus
	PageResult *PaginationResult
}

// Menus 菜单列表
type Menus []*Menu

// SplitAndGetAllRecordIDs 拆分父级路径并获取所有记录ID
func (a Menus) SplitAndGetAllRecordIDs() []string {
	var recordIDs []string
	for _, item := range a {
		recordIDs = append(recordIDs, item.RecordID)
		if item.ParentPath == "" {
			continue
		}

		pps := strings.Split(item.ParentPath, "/")
		for _, pp := range pps {
			var exists bool
			for _, recordID := range recordIDs {
				if pp == recordID {
					exists = true
					break
				}
			}
			if !exists {
				recordIDs = append(recordIDs, pp)
			}
		}
	}
	return recordIDs
}

// ToTrees 转换为菜单树列表
func (a Menus) ToTrees() MenuTrees {
	list := make(MenuTrees, len(a))
	for i, item := range a {
		list[i] = &MenuTree{
			RecordID:   item.RecordID,
			Code:       item.Code,
			Name:       item.Name,
			Type:       item.Type,
			Sequence:   item.Sequence,
			Icon:       item.Icon,
			Path:       item.Path,
			ParentID:   item.ParentID,
			ParentPath: item.ParentPath,
		}
	}
	return list
}

func (a Menus) fillLeafNodeID(tree *[]*MenuTree, leafNodeIDs *[]string) {
	for _, node := range *tree {
		if node.Children == nil || len(*node.Children) == 0 {
			*leafNodeIDs = append(*leafNodeIDs, node.RecordID)
			continue
		}
		a.fillLeafNodeID(node.Children, leafNodeIDs)
	}
}

// ToLeafRecordIDs 转换为叶子节点记录ID列表
func (a Menus) ToLeafRecordIDs() []string {
	var leafNodeIDs []string
	tree := a.ToTrees().ToTree()
	a.fillLeafNodeID(&tree, &leafNodeIDs)
	return leafNodeIDs
}

// MenuTree 菜单树
type MenuTree struct {
	RecordID   string       `json:"record_id" swaggo:"false,记录ID"`
	Code       string       `json:"code" binding:"required" swaggo:"true,菜单编号"`
	Name       string       `json:"name" binding:"required" swaggo:"true,菜单名称"`
	Type       int          `json:"type" binding:"required,max=3,min=1" swaggo:"true,菜单类型(1：模块 2：功能 3：资源)"`
	Sequence   int          `json:"sequence" swaggo:"false,排序值"`
	Icon       string       `json:"icon" swaggo:"false,菜单图标"`
	Path       string       `json:"path" swaggo:"false,访问路径"`
	ParentID   string       `json:"parent_id" swaggo:"false,父级内码"`
	ParentPath string       `json:"parent_path" swaggo:"false,父级路径"`
	Children   *[]*MenuTree `json:"children,omitempty" swaggo:"false,子级树"`
}

// MenuTrees 菜单树列表
type MenuTrees []*MenuTree

// ToTree 转换为树形结构
func (a MenuTrees) ToTree() []*MenuTree {
	mi := make(map[string]*MenuTree)
	for _, item := range a {
		mi[item.RecordID] = item
	}

	var list []*MenuTree
	for _, item := range a {
		if item.ParentID == "" {
			list = append(list, item)
			continue
		}
		if pitem, ok := mi[item.ParentID]; ok {
			if pitem.Children == nil {
				var children []*MenuTree
				children = append(children, item)
				pitem.Children = &children
				continue
			}
			*pitem.Children = append(*pitem.Children, item)
		}
	}
	return list
}
