package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss/tree"
	"golang.org/x/mod/modfile"
)

type ModNode struct {
	ID      string     `json:"id"` // Mod Path@Version
	Require []*ModNode `json:"require,omitempty"`
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

func buildDependencyTree(mf *modfile.File, modDepEntries map[string][]string) *ModNode {
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

		addDependencyChildren(n, modDepEntries)
	}
	return root
}

func addDependencyChildren(parent *ModNode, modDepEntries map[string][]string) {
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
		addDependencyChildren(child, modDepEntries)
	}
}

func extractModuleDependencies() (map[string][]string, error) {
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
