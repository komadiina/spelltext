package utils

import (
	"context"
	"sync"

	sq "github.com/Masterminds/squirrel"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
	"github.com/komadiina/spelltext/server/character/server"
)

var once = sync.Once{}

func InitializeEquipmentSlots(s *server.CharacterService) {
	once.Do(func() {
		ctx := context.Background()

		// get equipment slots
		sql, _, err := sq.Select("*").From("equip_slots").ToSql()

		if err != nil {
			s.Logger.Error(err)
			return
		}

		rows, err := s.DbPool.Query(ctx, sql)
		if err != nil {
			s.Logger.Error(err)
			return
		}

		var equipSlots []*pbRepo.EquipSlot
		for rows.Next() {
			es := &pbRepo.EquipSlot{}

			err := rows.Scan(
				&es.Id,
				&es.Code,
				&es.Name,
			)

			if err != nil {
				s.Logger.Error(err)
				return
			}

			equipSlots = append(equipSlots, es)
		}

		var characterIds []uint64

		sql, _, err = sq.Select("character_id").From("characters").ToSql()
		if err != nil {
			s.Logger.Error(err)
			return
		}

		rows, err = s.DbPool.Query(ctx, sql)
		if err != nil {
			s.Logger.Error(err)
			return
		}

		for rows.Next() {
			var characterId uint64

			err := rows.Scan(&characterId)
			if err != nil {
				s.Logger.Error(err)
				return
			}

			characterIds = append(characterIds, characterId)
		}

		// foreach equipment slot -> insert where (character_id, equipment_slot_id) is missing
		for _, cid := range characterIds {
			for _, es := range equipSlots {
				sql, _, err := sq.Insert("character_equipments").Columns("character_id", "equip_slot_id").Values(cid, es.Id).ToSql()

				if err != nil {
					s.Logger.Error(err)
					return
				}

				_, err = s.DbPool.Exec(ctx, sql)
				if err != nil {
					s.Logger.Warn(err)
					continue
				}
			}
		}
	})
}
