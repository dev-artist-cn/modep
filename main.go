package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss/tree"
	"golang.org/x/mod/modfile"
)

func main() {
	// 解析go.mod文件
	modFilePath := "go.mod" // 这里假设go.mod在当前目录，可按需修改路径
	modFileContent, err := os.ReadFile(modFilePath)
	if err != nil {
		fmt.Printf("打开go.mod文件失败: %v\n", err)
		os.Exit(1)
	}
	mf, err := modfile.Parse("go.mod", []byte(modFileContent), nil)
	if err != nil {
		log.Printf("解析go.mod文件失败: %v\n", err)
		os.Exit(1)
	}

	// 通过执行 go mod graph 命令获取模块依赖关系, 返回一个map，key为模块名，value为依赖的模块名列表
	modDepEntries, err := parseModGraph()
	if err != nil {
		log.Printf("构建模块依赖关系失败: %v\n", err)
		os.Exit(1)
	}

	// 构建模块依赖树, 以当前模块的直接依赖为起点
	depTree := buildDepTree(mf, modDepEntries)
	if len(depTree.Require) == 0 {
		log.Println("未找到模块的直接依赖, 请查看是否 go.mod 中是否有直接依赖项并执行 go mod tidy")
		os.Exit(1)
	}
	printTree(depTree)
}

type ModNode struct {
	ID      string     `json:"id"` // Mod Path@Version
	Require []*ModNode `json:"require,omitempty"`
}

func printTree(root *ModNode) {
	renderRoot := tree.Root(root.ID)
	buildRenderTree(root, renderRoot)
	fmt.Println(renderRoot.String())
}

func buildRenderTree(root *ModNode, renderRoot *tree.Tree) {
	for _, req := range root.Require {
		if len(req.Require) == 0 {
			renderRoot.Child(req.ID)
			continue
		}

		ch := tree.Root(req.ID)
		renderRoot.Child(ch)
		buildRenderTree(req, ch)
	}
}

func buildDepTree(mf *modfile.File, modDepEntries map[string][]string) *ModNode {
	var root = &ModNode{
		ID:      mf.Module.Mod.String(),
		Require: make([]*ModNode, 0),
	}

	for _, req := range mf.Require {
		if req.Indirect {
			continue
		}

		var n = &ModNode{
			ID:      req.Mod.String(),
			Require: make([]*ModNode, 0),
		}
		root.Require = append(root.Require, n)

		appendDepChildren(n, modDepEntries)
	}
	return root
}

func appendDepChildren(parent *ModNode, modDepEntries map[string][]string) {
	if parent == nil {
		return
	}
	vs, ok := modDepEntries[parent.ID]
	if !ok {
		return
	}

	for _, v := range vs {
		child := &ModNode{
			ID:      v,
			Require: make([]*ModNode, 0),
		}
		parent.Require = append(parent.Require, child)
		appendDepChildren(child, modDepEntries)
	}
}

func parseModGraph() (map[string][]string, error) {
	cmd := exec.Command("go", "mod", "graph")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("执行命令失败 go mod graph: %v", err)
	}

	modDepEntries := make(map[string][]string)
	// parse the graph output
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, " ")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid line: %s", line)
		}
		k, v := parts[0], parts[1]
		if _, ok := modDepEntries[k]; !ok {
			modDepEntries[k] = make([]string, 0)
		}
		modDepEntries[k] = append(modDepEntries[k], v)
	}
	return modDepEntries, nil
}
