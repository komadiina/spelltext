package hooks

import "github.com/komadiina/spelltext/client/types"

func CloseClients(c *types.SpelltextClient) {
	f := func(err error) {
		if err != nil {
			c.Logger.Error(err)
		}
	}

	if c.Clients != nil {
		f(c.Connections.Chat.Close())
		f(c.Connections.Store.Close())
		f(c.Connections.Inventory.Close())
		f(c.Connections.Character.Close())
		f(c.Connections.Gamba.Close())
		f(c.Connections.Auth.Close())
		f(c.Connections.Combat.Close())
		f(c.Connections.Build.Close())
	}
}
