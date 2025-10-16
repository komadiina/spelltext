package utils

import (
	"fmt"
	"strings"

	pbRepo "github.com/komadiina/spelltext/proto/repo"
	pb "github.com/komadiina/spelltext/proto/store"
)

func GetFullItemName(item *pbRepo.Item) string {
	return strings.Trim(fmt.Sprintf("%s %s %s", item.GetPrefix(), item.GetItemTemplate().GetName(), item.GetSuffix()), " ")
}

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
