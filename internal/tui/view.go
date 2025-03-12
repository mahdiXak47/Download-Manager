package tui

import "fmt"

func (m Model) View() string {
	s := "Download Manager v0.1\n\n"

	switch m.Menu {
	case "add":
		s += "Add Download\n"
		s += "URL: " + m.InputURL + "\n"
		s += "Queue: " + m.InputQueue + "\n"
		s += "\nPress Enter to add, Esc to cancel"

	case "list":
		s += "Downloads:\n"
		for i, d := range m.Downloads {
			cursor := " "
			if i == m.Selected {
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
