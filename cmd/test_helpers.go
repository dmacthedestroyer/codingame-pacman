package main

import "strings"

func BuildGameMap(cellString string) (gm GameMap) {
	rows := strings.Split(cellString, "\n")
	gm.height = len(rows) - 1
	for i, r := range rows {
		if i > 0 {
			gm.width = len(r)
			for _, c := range r {
				gm.cells = append(gm.cells, Cell{c})
			}
		}
	}

	return
}
