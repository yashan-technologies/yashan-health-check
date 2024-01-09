package checkcontroller

import (
	"yhc/commons/yasdb"

	"git.yasdb.com/pandora/tview"
	"github.com/gdamore/tcell/v2"
)

const (
	_ok = "OK"
)

func newCheckedList(title string, border bool) *tview.CheckList {
	checkedList := tview.NewCheckList()
	checkedList.SetTitle(title)
	checkedList.SetBorder(border)
	checkedList.ShowSecondaryText(false)
	return checkedList
}

func newTable(title string, border, borders bool) *tview.Table {
	table := tview.NewTable()
	table.SetTitle(title)
	table.SetBorder(border)
	table.SetBorders(borders)
	return table
}

func newModal(app *tview.Application, index *tview.Flex, err error) *tview.Modal {
	modal := tview.NewModal()
	modal.SetBackgroundColor(tcell.ColorRed)
	modal.SetText(err.Error())
	modal.AddButtons([]string{_ok})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == _ok {
			app.SetRoot(index, true)
		}
	})
	return modal
}

func newTextView(text string, border bool, align int, color tcell.Color) *tview.TextView {
	header := tview.NewTextView()
	header.SetBorder(border)
	header.SetTextAlign(align)
	header.SetText(text)
	header.SetTextColor(color)
	return header
}

func newButton(title string, border bool) *tview.Button {
	button := tview.NewButton(title)
	button.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite).Bold(true))
	button.SetActivatedStyle(tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlue).Bold(true))
	button.SetDisabledStyle(tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite).Bold(false))
	button.SetBorder(border)
	return button
}

func newPages(views ...*PagePrimitive) *tview.Pages {
	page := tview.NewPages()
	for _, view := range views {
		if view == nil {
			continue
		}
		page.AddPage(view.Name, view.Primitive, true, view.Show)
	}
	return page
}

func newYasdbPage(yasdb *yasdb.YashanDB) *PagePrimitive {
	return &PagePrimitive{
		Name:      _yasdb,
		Primitive: dbInfoPage(yasdb),
		Show:      true,
	}
}

func newFlex(title string, border bool, direction int) *tview.Flex {
	f := tview.NewFlex()
	f.SetBorder(border)
	f.SetTitle(title)
	f.SetDirection(direction)
	return f
}
