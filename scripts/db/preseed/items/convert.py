import json
import random

with open("data.json") as f:
    data = json.load(f)

slot_ids = {
    "head": 1,
    "neck": 2,
    "shoulder": 3,
    "back": 4,
    "chest": 5,
    "wrist": 6,
    "hands": 7,
    "legs": 8,
    "feet": 9,
    "finger_1": 10,
    "finger_2": 11,
    "mh": 12,
    "oh": 13,
}

armor_slots = {
    "head",
    "neck",
    "shoulder",
    "back",
    "chest",
    "wrist",
    "hands",
    "legs",
    "feet",
    "finger_1",
    "finger_2",
}
weapon_slots = {"mh", "oh"}

item_sql = []
chest_sql = []

item_template_id = 1
for slot, items in data.items():
    item_sql.append(f"-- {slot}")
    equip_slot_id = slot_ids[slot]
    item_type_id = 1 if slot in armor_slots else 2

    for it in items:
        name = it["item_name"].replace("'", "''")
        gold_price = it["gold"]
        hp = it["health_points"]
        power = it["power_points"]
        strp = it["strength_points"]
        sp = it["spellpower_points"]
        bonus_damage = random.randint(0, 10)
        bonus_armor = random.randint(-5, 20)

        sql_template = (
            f"INSERT INTO item_templates (name, item_type_id, equip_slot_id, gold_price) "
            f"VALUES ('{name}', {item_type_id}, {equip_slot_id}, {gold_price}); -- {item_template_id}"
        )

        sql_item = (
            f"INSERT INTO items (prefix, suffix, item_template_id, health, power, strength, spellpower, bonus_damage, bonus_armor) "
            f"VALUES ('','',{item_template_id},{hp},{power},{strp},{sp},{bonus_damage},{bonus_armor});"
        )

        item_sql.append(sql_template)
        item_sql.append(sql_item)
        item_template_id += 1

# preseed gamba_chest_contents table
chest_sql.append(
    "INSERT INTO gamba_chests (name, description, price) VALUES ('poor man''s', 'broke?', 0);"
)

chest_sql.append(
    "INSERT INTO gamba_chests (name, description, price) VALUES ('it''s okay', 'keep grinding', 25);"
)

chest_sql.append(
    "INSERT INTO gamba_chests (name, description, price) VALUES ('legend', 'legends only club', 69);"
)

chest_sql.append(
    "INSERT INTO gamba_chests (name, description, price) VALUES ('minecraft steve', 'diamondzzzzz!!', 420);"
)

gamba_chests_amount = 4
for i in range(gamba_chests_amount):
    for j in range(1, item_template_id):
        if random.randint(0, 100) < 20:
            continue

        chest_sql.append(
            f"INSERT INTO gamba_chest_contents (gamba_chest_id, item_id) "
            f"VALUES ({i+1}, {j});"
        )

with open("V0.4.1__populate_items.sql", "w") as f:
    f.write("\n".join(item_sql))

with open("V0.4.2__populate_gamba_chests.sql", "w") as f:
    f.write("\n".join(chest_sql))
