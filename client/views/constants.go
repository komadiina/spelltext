package views

import "github.com/rivo/tview"

// AppStorage accessors
var (
	SELECTED_VENDOR_ID   = "selectedVendorID"
	SELECTED_VENDOR_NAME = "selectedVendorName"
)

func STNewFlex() *tview.Flex {
	flex := tview.NewFlex().SetFullScreen(true)
	flex.SetBorderPadding(5, 5, 10, 10)
	return flex
}
