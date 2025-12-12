package site

import (
	"fmt"
	"testing"
)

func TestPager(t *testing.T) {
	var res [][]PagerItem = [][]PagerItem{
		{{Page: 1, Label: "|\u0026lt;", Active: false, Disabled: true}, {Page: -1, Label: "\u0026lt;-", Active: false, Disabled: true}, {Page: 1, Label: "1", Active: true, Disabled: false}, {Page: 2, Label: "2", Active: false, Disabled: false}, {Page: 3, Label: "-\u0026gt;", Active: false, Disabled: false}, {Page: 8, Label: "\u0026gt;|", Active: false, Disabled: false}},
		{{Page: 1, Label: "|\u0026lt;", Active: false, Disabled: true}, {Page: 0, Label: "\u0026lt;-", Active: false, Disabled: true}, {Page: 1, Label: "1", Active: false, Disabled: false}, {Page: 2, Label: "2", Active: true, Disabled: false}, {Page: 3, Label: "3", Active: false, Disabled: false}, {Page: 4, Label: "-\u0026gt;", Active: false, Disabled: false}, {Page: 8, Label: "\u0026gt;|", Active: false, Disabled: false}},
		{{Page: 1, Label: "|\u0026lt;", Active: false, Disabled: true}, {Page: 1, Label: "\u0026lt;-", Active: false, Disabled: false}, {Page: 2, Label: "2", Active: false, Disabled: false}, {Page: 3, Label: "3", Active: true, Disabled: false}, {Page: 4, Label: "4", Active: false, Disabled: false}, {Page: 5, Label: "-\u0026gt;", Active: false, Disabled: false}, {Page: 8, Label: "\u0026gt;|", Active: false, Disabled: false}},
		{{Page: 1, Label: "|\u0026lt;", Active: false, Disabled: false}, {Page: 2, Label: "\u0026lt;-", Active: false, Disabled: false}, {Page: 3, Label: "3", Active: false, Disabled: false}, {Page: 4, Label: "4", Active: true, Disabled: false}, {Page: 5, Label: "5", Active: false, Disabled: false}, {Page: 6, Label: "-\u0026gt;", Active: false, Disabled: false}, {Page: 8, Label: "\u0026gt;|", Active: false, Disabled: false}},
		{{Page: 1, Label: "|\u0026lt;", Active: false, Disabled: false}, {Page: 3, Label: "\u0026lt;-", Active: false, Disabled: false}, {Page: 4, Label: "4", Active: false, Disabled: false}, {Page: 5, Label: "5", Active: true, Disabled: false}, {Page: 6, Label: "6", Active: false, Disabled: false}, {Page: 7, Label: "-\u0026gt;", Active: false, Disabled: false}, {Page: 8, Label: "\u0026gt;|", Active: false, Disabled: false}},
		{{Page: 1, Label: "|\u0026lt;", Active: false, Disabled: false}, {Page: 4, Label: "\u0026lt;-", Active: false, Disabled: false}, {Page: 5, Label: "5", Active: false, Disabled: false}, {Page: 6, Label: "6", Active: true, Disabled: false}, {Page: 7, Label: "7", Active: false, Disabled: false}, {Page: 8, Label: "-\u0026gt;", Active: false, Disabled: false}, {Page: 8, Label: "\u0026gt;|", Active: false, Disabled: true}},
		{{Page: 1, Label: "|\u0026lt;", Active: false, Disabled: false}, {Page: 5, Label: "\u0026lt;-", Active: false, Disabled: false}, {Page: 6, Label: "6", Active: false, Disabled: false}, {Page: 7, Label: "7", Active: true, Disabled: false}, {Page: 8, Label: "8", Active: false, Disabled: false}, {Page: 9, Label: "-\u0026gt;", Active: false, Disabled: true}, {Page: 8, Label: "\u0026gt;|", Active: false, Disabled: true}},
		{{Page: 1, Label: "|\u0026lt;", Active: false, Disabled: false}, {Page: 6, Label: "\u0026lt;-", Active: false, Disabled: false}, {Page: 7, Label: "7", Active: false, Disabled: false}, {Page: 8, Label: "8", Active: true, Disabled: false}, {Page: 10, Label: "-\u0026gt;", Active: false, Disabled: true}, {Page: 8, Label: "\u0026gt;|", Active: false, Disabled: true}},
	}
	fmt.Printf("res: %v\n", res)

	for i := 1; i <= 8; i++ {
		idx := (i - 1)
		pg := NewPager(8, 2, i)
		pgv := pg.GetView()
		for j := 0; j < len(res[idx]); j++ {
			if pgv[j].Page != res[idx][j].Page ||
				pgv[j].Label != res[idx][j].Label ||
				pgv[j].Active != res[idx][j].Active ||
				pgv[j].Disabled != res[idx][j].Disabled {
				t.Errorf("View index %d has unexpected value: %v (expected %v)\n",
					j, pgv[j], res[idx][j])
			}
		}
	}
}
