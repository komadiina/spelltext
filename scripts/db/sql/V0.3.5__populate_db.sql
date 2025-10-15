insert into users values (DEFAULT, 'oggnjen', 'changeme', 'oggnjen@mail.com');
insert into heroes values (DEFAULT, 'paladin', 100, 150, 50, 0, 10, 5, 10, 0);
insert into heroes values (DEFAULT, 'mage', 50, 200, 10, 40, 10, 5, 0, 20);
insert into characters values (DEFAULT, 1, 'Oggnjen', 1, 1, 0, 200, 10, 5, 5, 5, 5);
insert into characters values (DEFAULT, 1, 'Klotzna', 2, 1, 0, 200, 10, 5, 5, 5, 5);

insert into item_types values (DEFAULT, 'IT_ARMOR', 'Armor');
insert into item_types values (DEFAULT, 'IT_WEAPON', 'Weapon');
insert into item_types values (DEFAULT, 'IT_CONSUMABLE', 'Consumable');

insert into equip_slots values (DEFAULT, 'ES_HEAD', 'Head');
insert into equip_slots values (DEFAULT, 'ES_CHEST', 'Chest');
insert into equip_slots values (DEFAULT, 'ES_HANDS', 'Hands');
insert into equip_slots values (DEFAULT, 'ES_LEGS', 'Legs');
insert into equip_slots values (DEFAULT, 'ES_MH', 'Main-hand');
insert into equip_slots values (DEFAULT, 'ES_OH', 'Off-hand');

insert into item_templates values (DEFAULT, 'Bronze helmet', 1, 1, '', 5, 0, 0, '{}');
insert into item_templates values (DEFAULT, 'Bronze chest', 1, 2, '', 5, 1, 1, '{}');
insert into item_templates values (DEFAULT, 'Bronze armplates', 1, 3, '', 5, 0, 1, '{}');
insert into item_templates values (DEFAULT, 'Bronze legplates', 1, 4, '', 5, 1, 1, '{}');
insert into item_templates values (DEFAULT, 'Shovel', 2, 5, '', 10, 0, 1, '{}');
insert into item_templates values (DEFAULT, 'Staff', 2, 5, '', 10, 0, 1, '{}');
insert into item_templates values (DEFAULT, 'Plank', 2, 6, '', 7, 0, 1, '{}');

insert into items values (DEFAULT, 'Farmer''s', '', 1, 2, 2, 0, 0, 0, 10);
insert into items values (DEFAULT, 'Nearly-smelted', '', 2, 5, 5, 0, 0, 0, 50);
insert into items values (DEFAULT, 'Smelly', '', 3, 0, 10, 0, 0, 5, 25);
insert into items values (DEFAULT, 'Unkept', '', 4, 0, 0, 0, 0, 0, 100);
insert into items values (DEFAULT, '', '', 5, 0, 0, 25, 0, 10, 0);
insert into items values (DEFAULT, '', '', 6, 10, 1, 2, 20, 2, 10);
insert into items values (DEFAULT, '', '', 7, 50, 0, 0, 0, 0, 100);

insert into consumable_templates values (DEFAULT, 'Health potion', 1, 5);
insert into consumable_templates values (DEFAULT, 'Power potion', 1, 5);
insert into consumables values (DEFAULT, 1, 'Lesser', '', 1, 25, 0, 0, 0, 0);
insert into consumables values (DEFAULT, 1, 'Greater', '', 1, 50, 0, 0, 0, 0);
insert into consumables values (DEFAULT, 1, 'Major', '', 1, 75, 0, 0, 0, 0);
insert into consumables values (DEFAULT, 2, 'Lesser', '', 1, 0, 25, 0, 0, 0);
insert into consumables values (DEFAULT, 2, 'Greater', '', 1, 0, 50, 0, 0, 0);
insert into consumables values (DEFAULT, 2, 'Major', '', 1, 0, 75, 0, 0, 0);

-- insert into character_inventories values ();
-- insert into item_instances values ();

insert into vendors values (DEFAULT, 'Vendor1', 'Armor');
insert into vendors values (DEFAULT, 'Vendor2', 'Weapons');
insert into vendors values (DEFAULT, 'Vendor3', 'Consumables');
insert into vendors values (DEFAULT, 'Vendor4', 'Armor+Weapons');
insert into vendor_wares values (1, 1);
insert into vendor_wares values (2, 2);
insert into vendor_wares values (3, 3);
insert into vendor_wares values (4, 1);
insert into vendor_wares values (4, 2);