package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/tree"
	"golang.org/x/mod/modfile"
)

func main() {
	// Define a flag for the -p argument
	printTree := flag.Bool("p", false, "Print the dependency tree")
	flag.Parse()

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
	modDepEntries, err := extractModuleDependencies()
	if err != nil {
		log.Printf("构建模块依赖关系失败: %v\n", err)
		os.Exit(1)
	}

	// 构建模块依赖树, 以当前模块的直接依赖为起点
	depTree := buildDependencyTree(mf, modDepEntries)
	if len(depTree.Require) == 0 {
		log.Println("未找到模块的直接依赖, 请查看是否 go.mod 中是否有直接依赖项并执行 go mod tidy")
		os.Exit(1)
	}

	if *printTree {
		renderRoot := tree.Root(depTree.ID)
		buildRenderTree(depTree, renderRoot)
		fmt.Println(renderRoot.String())
		return
	}

	p := tea.NewProgram(initialTreeModel(depTree))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
