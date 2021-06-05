package ui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

type Results struct {
	*tview.TextView
}

func newResults() *Results {
	view := tview.NewTextView().SetDynamicColors(true)
	fmt.Fprint(view,
		strings.Repeat("\n", 99)+
			"[green::b]秋天的后半夜，月亮下去了，太阳还没有出，只剩下一片乌蓝的天；除了夜游的东西，什么都睡着。华老栓忽然坐起身，擦着火柴，点上遍身油腻的灯盏，茶馆的两间屋子里，便弥满了青白的光。")
	view.ScrollToEnd()
	view.SetBorder(true).SetTitle("flac: learn 中文")

	return &Results{view}
}
