import json
import os

with open("data.json") as f:
    data = json.load(f)

npc_templates = data["templates"]
npcs = data["npcs"]

sqls = {
    "npc_templates": [],
    "npcs": [],
}

for template in npc_templates:
    sqls["npc_templates"].append(
        (
            f"INSERT INTO npc_templates (name, min_level, max_level, health_points, base_damage, base_xp_reward, drop_item_id, gold_reward) "
            f"VALUES ('{template["name"].replace("'", "''")}', {template["min_level"]}, {template["max_level"]}, {template['health_points']}, {template['base_damage']}, {template['base_xp_reward']}, {template['drop_item_id']}, {template['gold_reward']});"
        )
    )

for npc in npcs:
    sqls["npcs"].append(
        (
            f"INSERT INTO npcs (prefix, suffix, template_id, health_multiplier, damage_multiplier) "
            f"VALUES ('{npc['prefix'].replace("'", "''")}', '{npc['suffix'].replace("'", "''")}', {npc['template_id']}, {npc['health_multiplier']}, {npc['damage_multiplier']});"
        )
    )


with open("V0.5.0__populate_npc_templates.sql", "w") as f:
    f.write("\n".join(sqls["npc_templates"]))

with open("V0.5.1__populate_npcs.sql", "w") as f:
    f.write("\n".join(sqls["npcs"]))
