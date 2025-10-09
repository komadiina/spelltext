package utils

import (
	"fmt"
	pb "github.com/komadiina/spelltext/proto/store"
)

func BuyItem() {}

func GetItemName(item *pb.Item) string {
	var prefix string = ""
	var suffix string = ""

	if len(item.GetPrefix()) == 0 {
		prefix = ""
	} else {
		prefix = item.GetPrefix() + " "
	}

	if len(item.GetSuffix()) == 0 {
		suffix = ""
	} else {
		suffix = " " + item.GetSuffix()
	}

	return fmt.Sprintf("%s%s%s", prefix, item.GetName(), suffix)
}