package tui

import "fmt"

func (m Model) View() string {
	s := "Download Manager v0.1\n\n"

	switch m.menu {
	case "add":
		s += "Add Download\n"
		s += "URL: " + m.inputURL + "\n"
		s += "Queue: " + m.inputQueue + "\n"
		s += "\nPress Enter to add, Esc to cancel"

	case "list":
		s += "Downloads:\n"
		for i, d := range m.downloads {
			cursor := " "
			if i == m.selected {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s [%s] %.1f%%\n", cursor, d.URL, d.Status, d.Progress)
		}

	default:
		s += "Commands:\n"
		s += "a: Add download\n"
		s += "l: List downloads\n"
		s += "p: Pause selected\n"
		s += "r: Resume selected\n"
		s += "q: Quit\n"
	}

	return s
}
