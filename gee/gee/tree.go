package gee

//路径树，支持“：xxx”、“*xxx”等动态路由

import (
	"log"
	"strings"
)

// gin 中把wild节点放最后
type node struct {
	children []*node
	fullPath string //完整路径，只在叶节点存在
	part     string //该节点的路径部分
	isWild   bool   //是否精确匹配
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.isWild == true || child.part == part {
			return child
		}
	}
	return nil
}

func (n *node) matchAllChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.isWild == true || child.part == part {
			nodes = append(nodes, child)
		}
	}
	return nodes

}

func (n *node) insert(fullPath string, parts []string, height int) {
	if n == nil {
		log.Fatal("n is a nullptr")
	}

	if len(parts) == height {
		n.fullPath = fullPath
		return
	}

	curPart := parts[height]

	child := n.matchChild(curPart)
	if child == nil {
		child = &node{part: curPart, isWild: (curPart[0] == ':' || curPart[0] == '*')}
		n.children = append(n.children, child)
	}
	child.insert(fullPath, parts, height+1)

}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		//非叶节点
		if n.fullPath != "" {
			return n
		}
		return nil
	}

	curPart := parts[height]
	children := n.matchAllChildren(curPart)

	for _, child := range children {
		if result := child.search(parts, height+1); result != nil {
			return result
		}
	}

	return nil

}
