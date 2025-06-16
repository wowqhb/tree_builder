package builder

type DataEntity[T any] struct {
	Id       int64 `json:"id,omitempty"`
	ParentId int64 `json:"parent_id,omitempty"`
	Data     T     `json:"data,omitempty"`
}

type TreeBuiler[T any, TR any] struct {
	dataEntities  []*DataEntity[T]
	rootFunc      func(data T) *TR
	appendSubFunc func(parent *TR, sub ...*TR)
}

// NewTreeBuilder 创建并返回一个新的 TreeBuilder 实例。
// 它负责使用提供的数据实体和函数来构建一个树形结构。
// 参数:
//   - dataEntities: 一个数据实体的切片，用于构建树形结构。
//   - rootFunc: 一个函数，用于确定数据实体的根节点。
//   - appendSubFunc: 一个函数，用于将子节点附加到父节点。
//   - initDataFunc: 一个函数，用于初始化数据实体。
//
// 返回值:
//   - *TreeBuiler[T, TR]: 返回一个指向新创建的 TreeBuilder 实例的指针。
func NewTreeBuilder[T any, TR any](
	dataEntities []T,
	rootFunc func(data T) *TR,
	appendSubFunc func(parent *TR, sub ...*TR),
	initDataFunc func(data T) *DataEntity[T],
) *TreeBuiler[T, TR] {
	// 创建一个新的 TreeBuilder 实例。
	ret := &TreeBuiler[T, TR]{
		dataEntities:  nil,
		rootFunc:      rootFunc,
		appendSubFunc: appendSubFunc,
	}

	// 遍历提供的数据实体，使用 initDataFunc 初始化每个数据实体，并添加到 TreeBuilder 实例中。
	for _, x := range dataEntities {
		ret.dataEntities = append(ret.dataEntities, initDataFunc(x))
	}

	// 返回创建好的 TreeBuilder 实例。
	return ret
}

// BuildTree 构建树结构。
// 该方法根据给定的根节点ID，在提供的数据实体中搜索并构建一个树结构。
// 如果找到了匹配的根节点，它会使用rootFunc函数来创建根节点对象，
// 并使用subFunc函数查找其子节点。如果存在子节点，它们将通过appendSubFunc函数附加到根节点上。
// 类型参数T代表数据实体的类型，TR代表树结构结果的类型。
func (self *TreeBuiler[T, TR]) BuildTree(rootId int64) *TR {
	// 检查是否有数据实体存在。
	if len(self.dataEntities) > 0 {
		// 遍历所有数据实体，寻找与rootId匹配的根节点。
		for _, x := range self.dataEntities {
			if x.Id == rootId {
				// 使用rootFunc函数根据数据创建根节点对象。
				root := self.rootFunc(x.Data)
				// 如果根节点为空，则直接返回。
				if root == nil {
					return root
				}
				// 使用subFunc函数查找根节点的所有子节点。
				subs := self.subFunc(x.Id)
				// 如果存在子节点，将它们附加到根节点上。
				if len(subs) > 0 {
					self.appendSubFunc(root, subs...)
				}
				// 返回构建完成的根节点（树结构）。
				return root
			}
		}
	}
	// 如果没有找到匹配的根节点或数据实体为空，则返回nil。
	return nil
}

// subFunc 是一个递归函数，用于根据父ID构建树形结构的子节点。
// 该函数属于TreeBuilder泛型结构体，用于处理数据实体并构建子节点列表。
// 参数:
//
//	parentId (int64): 当前处理的父节点ID。
//
// 返回值:
//
//	[]*TR: 返回一个指向TR类型切片的指针，表示构建的子节点列表。
func (self *TreeBuiler[T, TR]) subFunc(parentId int64) []*TR {
	// 初始化返回值切片。
	var ret []*TR
	// 检查是否有数据实体需要处理。
	if len(self.dataEntities) > 0 {
		// 遍历所有数据实体。
		for _, x := range self.dataEntities {
			// 判断当前数据实体是否属于当前父节点。
			if x.ParentId == parentId {
				// 使用当前数据实体的数据调用rootFunc函数，构建子节点。
				sub := self.rootFunc(x.Data)
				// 如果构建的子节点为空，则跳过当前循环迭代。
				if sub == nil {
					continue
				}
				// 递归调用subFunc函数，构建当前子节点的子节点。
				subs := self.subFunc(x.Id)
				// 如果存在子节点，则调用appendSubFunc函数将它们附加到当前子节点上。
				if len(subs) > 0 {
					self.appendSubFunc(sub, subs...)
				}
				// 将构建的子节点添加到返回值切片中。
				ret = append(ret, sub)
			}
		}
	}
	// 返回构建的子节点列表。
	return ret
}
