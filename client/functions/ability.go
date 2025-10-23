package functions

import (
	"fmt"

	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/types"
	pbBuild "github.com/komadiina/spelltext/proto/build"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
)

func GetAbilities(c *types.SpelltextClient) (*[]*pbRepo.Ability, *[]*pbRepo.Ability, *[]*pbRepo.Ability) {
	req := &pbBuild.ListAbilitiesRequest{Character: c.Storage.SelectedCharacter}
	abilities, err := c.Clients.BuildClient.ListAbilities(*c.Context, req)
	if err != nil {
		c.Logger.Error(err)
		return nil, nil, nil
	}

	return &abilities.Available, &abilities.Locked, &abilities.Upgraded
}

func GetSpellDetails(ability *pbRepo.Ability) string {
	return fmt.Sprintf(
		`[%s][::b]%s[::-][""][white][""]%s%s%s%sbase damage: [%s]%d[""][white][""]%spower cost: [%s]%d[""][white][""]%scoefficients:%sstrength: [%s]%.2f%%[""][white][""] (%.2f%% per level)%sspellpower: [%s]%.2f[""][white][""] (%.2f%% per level)%scost: %d talent points[""]%s`,
		constants.TEXT_COLOR_GOLD,
		ability.Name, "\n",
		ability.Description, "\n", "\n",
		constants.TEXT_COLOR_DAMAGE, ability.BaseDamage, "\n",
		constants.TEXT_COLOR_POWER, ability.PowerCost,
		"\n++++++++++++++++++++++++++++++\n", "\n",
		constants.TEXT_COLOR_STRENGTH, ability.StrengthMultiplier*100.0, ability.StMultPerlevel*100.0, "\n",
		constants.TEXT_COLOR_SPELLPOWER, ability.SpellpowerMultiplier*100.0, ability.SpMultPerlevel*100.0, "\n",
		ability.TalentPointCost,
		"\n++++++++++++++++++++++++++++++",
	)
}

func GetSpellDetailsHelp() string {
	return fmt.Sprintf(
		`coefficients: how much [%s]strength[""][white][""] and [%s]spellpower[""][white][""] affect the output damage of the spell.
for example: base_damage=100,strength=100%% => 100 + (100 * 1.0) = 200 total damage.
talent points: how many [%s]talent points[""][white][""] are required to unlock the spell.`,
		constants.TEXT_COLOR_STRENGTH,
		constants.TEXT_COLOR_SPELLPOWER,
		constants.TEXT_COLOR_POWER,
	)
}

func UpgradeAbility(c *types.SpelltextClient, ability *pbRepo.Ability) error {
	req := &pbBuild.UpgradeAbilityRequest{CharacterId: c.Storage.SelectedCharacter.CharacterId, AbilityId: ability.Id}
	_, err := c.Clients.BuildClient.UpgradeAbility(*c.Context, req)
	return err
}
