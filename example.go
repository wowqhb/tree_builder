package main

import (
	"encoding/json"
	"fmt"

	"github.com/wowqhb/tree_builder/builder"
)

// ExampleData 用于定义示例数据结构，它有ID、父ID和名称字段
type ExampleData struct {
	Id       int
	ParentId int
	Name     string
}

// ExampleTree 用于定义示例树结构，包含树节点的所有必要信息
type ExampleTree struct {
	Id       int
	ParentId int
	Name     string
	Children []*ExampleTree
}

func main() {
	// 初始化示例数据数组，模拟从数据库或其他来源获取的数据
	inArr := []*ExampleData{
		{
			Id:       1,
			ParentId: 0,
			Name:     "01",
		},
		{
			Id:       2,
			ParentId: 1,
			Name:     "01-02",
		},
		{
			Id:       3,
			ParentId: 1,
			Name:     "01-03",
		},
		{
			Id:       4,
			ParentId: 2,
			Name:     "01-02-04",
		},
		{
			Id:       5,
			ParentId: 2,
			Name:     "01-02-05",
		},
		{
			Id:       6,
			ParentId: 5,
			Name:     "01-02-05-06",
		},
	}

	// rootFunc 是一个转换函数，用于将示例数据转换为树结构的根节点
	rootFunc := func(data *ExampleData) *ExampleTree {
		return &ExampleTree{
			Id:       data.Id,
			Name:     data.Name,
			ParentId: data.ParentId,
			Children: nil,
		}
	}

	// appendSubFunc 用于向父节点添加子节点
	appendSubFunc := func(parent *ExampleTree, sub ...*ExampleTree) {
		parent.Children = append(parent.Children, sub...)
	}

	// initDataFunc 用于初始化数据实体，将示例数据包装在DataEntity中
	initDataFunc := func(data *ExampleData) *builder.DataEntity[*ExampleData] {
		return &builder.DataEntity[*ExampleData]{
			Id:       int64(data.Id),
			ParentId: int64(data.ParentId),
			Data:     data,
		}
	}

	// 创建一个新的树构建器实例，使用示例数据和定义的函数
	treeBuilder := builder.NewTreeBuilder(inArr, rootFunc, appendSubFunc, initDataFunc)
	// 构建树结构，从ID为1的根节点开始
	retTree := treeBuilder.BuildTree(1)
	// 根据构建结果输出相应的信息
	if retTree == nil {
		fmt.Println("no root")
	} else {
		// 将树结构转换为JSON格式并打印
		jsonBytes, _ := json.Marshal(retTree)
		fmt.Println(string(jsonBytes))
	}
}
