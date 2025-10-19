insert into users values (DEFAULT, 'oggnjen', 'changeme', 'oggnjen@mail.com');

insert into heroes values (DEFAULT, 'Paladin', 100, 150, 50, 0, 10, 5, 10, 0);
insert into heroes values (DEFAULT, 'Mage', 50, 200, 10, 40, 10, 5, 0, 20);
insert into heroes values (DEFAULT, 'Warrior', 100, 100, 50, 0, 10, 5, 10, 0);
insert into heroes values (DEFAULT, 'Rogue', 50, 200, 10, 40, 10, 5, 0, 20);

insert into characters values (DEFAULT, 1, 'Oggnjen', 1, 1, 0, 9999999, 10, 5, 5, 5, 5);
insert into characters values (DEFAULT, 1, 'Klotzna', 2, 1, 0, 9999999, 10, 5, 5, 5, 5);
insert into characters values (DEFAULT, 1, 'Lanlan', 3, 1, 0, 9999999, 10, 5, 5, 5, 5);
insert into characters values (DEFAULT, 1, 'Jorunn', 4, 1, 0, 9999999, 10, 5, 5, 5, 5);

insert into item_types values (DEFAULT, 'IT_ARMOR', 'Armor');
insert into item_types values (DEFAULT, 'IT_WEAPON', 'Weapon');

insert into item_types values (DEFAULT, 'IT_TRINKET', 'Trinket'); -- todo
insert into item_types values (DEFAULT, 'IT_CONSUMABLE', 'Consumable'); -- todo
insert into item_types values (DEFAULT, 'IT_VANITY', 'Vanity'); -- todo

insert into equip_slots values (DEFAULT, 'ES_HEAD', 'Head'); -- id=1
insert into equip_slots values (DEFAULT, 'ES_NECK', 'Neck'); -- id=2
insert into equip_slots values (DEFAULT, 'ES_SHOULDERS', 'Shoulders'); -- id=3
insert into equip_slots values (DEFAULT, 'ES_BACK', 'Back'); -- id=4
insert into equip_slots values (DEFAULT, 'ES_CHEST', 'Chest'); -- id=5
insert into equip_slots values (DEFAULT, 'ES_WRIST', 'Wrist'); -- id=6
insert into equip_slots values (DEFAULT, 'ES_HANDS', 'Hands'); -- id=7
insert into equip_slots values (DEFAULT, 'ES_LEGS', 'Legs'); -- id=8
insert into equip_slots values (DEFAULT, 'ES_FEET', 'Feet'); -- id=9
insert into equip_slots values (DEFAULT, 'ES_FINGER_1', 'Ring'); -- id=10
insert into equip_slots values (DEFAULT, 'ES_FINGER_2', 'Ring'); -- id=11
insert into equip_slots values (DEFAULT, 'ES_MH', 'Main-hand'); -- id=12
insert into equip_slots values (DEFAULT, 'ES_OH', 'Off-hand'); -- id=13

insert into vendors values (DEFAULT, 'Taj''ma', 'Armor');
insert into vendors values (DEFAULT, 'Olut Ik''habad', 'Weapons');
insert into vendors values (DEFAULT, 'Potionist Ezra', 'Consumables');
insert into vendors values (DEFAULT, 'Villagesmith Cyril', 'Armor+Weapons');
insert into vendor_wares values (1, 1);
insert into vendor_wares values (2, 2);
insert into vendor_wares values (3, 3);
insert into vendor_wares values (4, 1);
insert into vendor_wares values (4, 2);