package constants

import "github.com/gdamore/tcell/v2"

const (
	PAGE_MAINMENU  = "pg-mainmenu"
	PAGE_ARMORY    = "pg-armory"
	PAGE_INVENTORY = "pg-inventory"
	PAGE_PROGRESS  = "pg-progress"
	PAGE_COMBAT    = "pg-combat"
	PAGE_STORE     = "pg-store"
	PAGE_VENDOR    = "pg-vendor"
	PAGE_GAMBA     = "pg-gamba"
	PAGE_CHAT      = "pg-chat"
	PAGE_LOGIN     = "pg-login"
)

const (
	ITEM_TYPE_UNSPECIFIED     = "IT00" // unused
	ITEM_TYPE_ARMOR           = "IT01"
	ITEM_TYPE_CONSUMABLE      = "IT02"
	ITEM_TYPE_WEAPON          = "IT03"
	ITEM_TYPE_VANITY          = "IT04"
	ITEM_TYPE_WEARABLE_VANITY = "IT05"
	ITEM_TYPE_GOLD_DROP       = "IT06"
	ITEM_TYPE_TOKEN_DROP      = "IT07"
	ITEM_TYPE_OTHER           = "IT99"
)

const (
	CURRENT_USER_ID      = "currentUserId"
	CURRENT_USER_NAME    = "currentUserName"
	SELECTED_CHARACTER   = "selectedCharacter"
	SELECTED_VENDOR_ID   = "selectedVendorID"
	SELECTED_VENDOR_NAME = "selectedVendorName"
)

var SHORTCUTS = map[tcell.Key]string{
	tcell.KeyCtrlA: PAGE_ARMORY,
	tcell.KeyCtrlI: PAGE_INVENTORY,
	tcell.KeyCtrlP: PAGE_PROGRESS,
	tcell.KeyCtrlC: PAGE_COMBAT,
	tcell.KeyCtrlS: PAGE_STORE,
	tcell.KeyCtrlG: PAGE_GAMBA,
	tcell.KeyCtrlY: PAGE_CHAT,
}
