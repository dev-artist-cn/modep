package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type treeModel struct {
	root     *ModNode
	cursor   int
	expanded map[int]bool
}

func initialTreeModel(root *ModNode) treeModel {
	return treeModel{
		root:     root,
		expanded: make(map[int]bool),
	}
}

func (m treeModel) Init() tea.Cmd {
	return nil
}

func (m treeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.root.Require)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.expanded[m.cursor] = !m.expanded[m.cursor]
		}
	}
	return m, nil
}

func (m treeModel) View() string {
	s := "Dependency Tree:\n\n"
	for i, req := range m.root.Require {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}
		expanded := "+" // not expanded
		if m.expanded[i] {
			expanded = "-" // expanded!
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, expanded, req.ID)
		if m.expanded[i] {
			for _, child := range req.Require {
				s += fmt.Sprintf("    - %s\n", child.ID)
			}
		}
	}
	s += "\nPress q to quit.\n"
	return s
}
