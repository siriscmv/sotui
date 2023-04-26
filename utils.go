package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
)

func GetSelectedRow(m Model, currRow table.Row) (ResponseItem, bool) {
	for _, row := range m.response.Items {
		if currRow[0] == fmt.Sprintf("%d", row.Score) && currRow[2] == fmt.Sprintf("%d", row.ViewCount) {
			//TODO: Use ID for comparing rows
			return row, true
		}
	}
	return ResponseItem{}, false
}
