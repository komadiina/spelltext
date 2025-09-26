package functions

import "github.com/komadiina/spelltext/client/types"

func GetUserByUsername(username string) types.SpelltextUser {
	return types.SpelltextUser{Username: username}
}
