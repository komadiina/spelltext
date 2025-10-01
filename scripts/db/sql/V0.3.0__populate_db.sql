-- heroes
insert into `heroes`
values (
    null,
    'Druid',
    100,
    100,
    5,
    5,
    50,
    25,
    10,
    20
  );

-- id=1
insert into `heroes`
values (
    null,
    'Mage',
    50,
    200,
    5,
    25,
    25,
    50,
    5,
    30
  );

-- id=2
insert into `heroes`
values (
    null,
    'Warrior',
    150,
    50,
    25,
    5,
    50,
    20,
    25,
    5
  );

-- id=3
-- stat types
insert into `stat_types`
VALUES (null, 'ST01', 'Strength');

-- id=1
insert into `stat_types`
VALUES (null, 'ST02', 'Spellpower');

-- id=2
insert into `stat_types`
VALUES (null, 'ST03', 'Health');

-- id=3
insert into `stat_types`
values (null, 'ST04', 'Power');

-- id=4
-- item types
insert into `item_types`
values (null, 'IT01', 'Armor');

-- id=1
insert into `item_types`
values (null, 'IT02', 'Consumable');

-- id=2
insert into `item_types`
values (null, 'IT03', 'Weapon');

-- id=3
insert into `item_types`
values (null, 'IT04', 'Vanity');

-- id=4
insert into `item_types`
values (null, 'IT05', 'Wearable vanity');

-- id=5
insert into `item_types`
values (null, 'IT06', 'Gold drop');

-- id=6
insert into `item_types`
values (null, 'IT07', 'Token drop');

-- id=7
-- equip slots
insert into `equip_slots`
values (null, 'ES01', 'Head');

-- id=1
insert into `equip_slots`
values (null, 'ES02', 'Neck');

-- id=2
insert into `equip_slots`
values (null, 'ES03', 'Shoulders');

-- id=3
insert into `equip_slots`
values (null, 'ES04', 'Gloves');

-- id=4
insert into `equip_slots`
values (null, 'ES05', 'Ring');

-- id=5
insert into `equip_slots`
values (null, 'ES06', 'Ring');

-- id=6
insert into `equip_slots`
values (null, 'ES07', 'Chest');

-- id=7
insert into `equip_slots`
values (null, 'ES08', 'Waist');

-- id=8
insert into `equip_slots`
values (null, 'ES09', 'Backpocket');

-- id=9
insert into `equip_slots`
values (null, 'ES10', 'Back');

-- id=10
insert into `equip_slots`
values (null, 'ES11', 'Legs');

-- id=12
insert into `equip_slots`
values (null, 'ES12', 'Shins');

-- id=13
insert into `equip_slots`
values (null, 'ES13', 'Feet');

-- id=14
-- item templates
insert into `item_templates`
values (
    null,
    'Peasant Helmet',
    1,
    1,
    0,
    1,
    1,
    'helmet',
    '{}'
  );

-- id=1
insert into `item_templates`
values (
    null,
    'Peasant Vest',
    1,
    1,
    0,
    1,
    1,
    'vest',
    '{}'
  );

-- id=2
insert into `item_templates`
values (
    null,
    'Wedding ring',
    1,
    1,
    0,
    1,
    1,
    'ring',
    '{}'
  );

-- id=3
insert into `item_templates`
values (
    null,
    'Health Potion',
    2,
    1,
    1,
    5,
    1,
    'hpot',
    '{}'
  );

-- id=4
insert into `item_templates`
values (
    null,
    'Power Potion',
    2,
    1,
    1,
    5,
    1,
    'ppot',
    '{}'
  );

-- id=5
insert into `item_templates`
values (
    null,
    'Strength Potion',
    2,
    2,
    1,
    5,
    1,
    'spot',
    '{}'
  );

-- id=6
insert into `item_templates`
values (
    null,
    'Spellpower Potion',
    2,
    2,
    1,
    5,
    1,
    'sppot',
    '{}'
  );

-- id=7
insert into `item_templates`
values (
    null,
    'Stoneblade',
    3,
    1,
    0,
    1,
    1,
    'weapon',
    '{}'
  );

-- id=8
insert into `item_templates`
values (
    null,
    'Warstaff',
    3,
    1,
    0,
    1,
    1,
    'weapon',
    '{}'
  );

-- id=9
insert into `item_templates`
values (
    null,
    'Sturdy shovel',
    3,
    1,
    0,
    1,
    1,
    'weapon',
    '{}'
  );

-- id=10
insert into `item_templates`
values (
    null,
    'Ygmirs head',
    4,
    2,
    0,
    1,
    1,
    'vanity',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Lucky? charm',
    5,
    1,
    0,
    1,
    1,
    'wvanity',
    '{}'
  );

-- id=12
insert into `item_templates`
values (
    null,
    'Honest gold bag',
    6,
    1,
    0,
    1,
    1,
    'gdrop',
    '{}'
  );

-- id=13
insert into `item_templates`
values (
    null,
    'Thief sack',
    6,
    2,
    0,
    1,
    1,
    'gdrop',
    '{}'
  );

-- id=14
insert into `item_templates`
values (
    null,
    'Bankrobber sack',
    6,
    3,
    0,
    1,
    1,
    'gdrop',
    '{}'
  );

-- id=15
insert into `item_templates`
values (
    null,
    'Gold bar',
    6,
    4,
    0,
    1,
    1,
    'gdrop',
    '{}'
  );

-- id=16
insert into `item_templates`
values (
    null,
    'Rugged token',
    7,
    1,
    0,
    1,
    1,
    'tdrop',
    '{}'
  );

-- id=17
insert into `item_templates`
values (
    null,
    'Polished token',
    7,
    2,
    0,
    1,
    1,
    'tdrop',
    '{}'
  );

-- id=18
insert into `item_templates`
values (
    null,
    'Headhunter token',
    7,
    3,
    0,
    1,
    1,
    'tdrop',
    '{}'
  );

-- id=19
insert into `item_templates`
values (
    null,
    'Worldslayer token',
    7,
    4,
    0,
    1,
    1,
    'tdrop',
    '{}'
  );

-- id=20
insert into `item_templates`
values (
    null,
    'Traveler Cloak',
    1,
    1,
    0,
    1,
    1,
    'cloak',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Cloak of Ember',
    1,
    2,
    0,
    1,
    1,
    'cloak',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Moonlit Shawl',
    1,
    3,
    0,
    1,
    1,
    'cloak',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Hunters Quiver',
    1,
    1,
    0,
    1,
    1,
    'backpocket',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Druidic Talisman',
    1,
    2,
    0,
    1,
    1,
    'trinket',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Mage Focus',
    1,
    2,
    0,
    1,
    1,
    'trinket',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Adventurer Necklace',
    1,
    1,
    0,
    1,
    1,
    'necklace',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Sage Locket',
    1,
    3,
    0,
    1,
    1,
    'necklace',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Lucky Rabbit Foot',
    5,
    1,
    0,
    1,
    1,
    'wvanity',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Iron Boots',
    1,
    1,
    0,
    1,
    1,
    'boots',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Plate Greaves',
    1,
    2,
    0,
    1,
    1,
    'shins',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Silken Gloves',
    1,
    2,
    0,
    1,
    1,
    'gloves',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Band of Fortitude',
    1,
    3,
    0,
    1,
    1,
    'ring',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Ring of Arcana',
    1,
    3,
    0,
    1,
    1,
    'ring',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Lesser Healing Potion',
    2,
    1,
    1,
    5,
    1,
    'hpot',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Greater Healing Potion',
    2,
    2,
    1,
    3,
    1,
    'hpot',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Greater Power Potion',
    2,
    2,
    1,
    3,
    1,
    'ppot',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Mana Crystal',
    2,
    2,
    1,
    5,
    1,
    'mcrystal',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Bandage',
    2,
    1,
    1,
    10,
    1,
    'bandage',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Rusted Dagger',
    3,
    1,
    0,
    1,
    1,
    'weapon',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Wicked Dirk',
    3,
    2,
    0,
    1,
    1,
    'weapon',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Ancient Longsword',
    3,
    3,
    0,
    1,
    1,
    'weapon',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Runed Staff',
    3,
    3,
    0,
    1,
    1,
    'weapon',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Shadow Cloak (Vanity)',
    4,
    1,
    0,
    1,
    1,
    'vanity',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Bronze Token',
    7,
    1,
    0,
    1,
    1,
    'tdrop',
    '{}'
  );

insert into `item_templates`
values (
    null,
    'Emerald Token',
    7,
    3,
    0,
    1,
    1,
    'tdrop',
    '{}'
  );

-- armors
-- 							template		equip_slot	prefix	suffix			armor	health	power	strength	spellpower
-- cloaks / backs
insert into `armors`
values (
    21,
    10,
    'Travel-worn',
    'of Endurance',
    5,
    20,
    0,
    0,
    5
  );

insert into `armors`
values (
    22,
    10,
    null,
    'of Ember',
    8,
    10,
    0,
    0,
    10
  );

insert into `armors`
values (
    23,
    10,
    'Moonlit',
    'of the Night',
    6,
    15,
    0,
    0,
    12
  );

insert into `armors`
values (
    24,
    9,
    "Hunter\'s",
    null,
    0,
    0,
    5,
    3,
    0
  );

-- quiver in backpocket gives power/str
-- helmets
insert into `armors`
values (
    1,
    1,
    null,
    'of the Owl',
    10,
    5,
    5,
    0,
    20
  );

insert into `armors`
values (
    1,
    1,
    null,
    'of the Bear',
    50,
    10,
    10,
    5,
    5
  );

insert into `armors`
values (
    1,
    1,
    null,
    'of the Beast',
    100,
    50,
    10,
    20,
    0
  );

-- chests
insert into `armors`
values (
    2,
    7,
    'Armored',
    null,
    200,
    5,
    -25,
    -5,
    -5
  );

insert into `armors`
values (
    2,
    7,
    "Mage's",
    'of the Lost',
    -50,
    -20,
    50,
    0,
    50
  );

-- ring
insert into `armors`
values (
    3,
    6,
    "Farmer's",
    null,
    10,
    5,
    0,
    0,
    0
  );

insert into `armors`
values (
    3,
    6,
    "Conjurer's",
    null,
    -5,
    5,
    5,
    0,
    10
  );

-- neck / trinkets / rings
insert into `armors`
values (
    25,
    9,
    null,
    'of the Grove',
    0,
    25,
    0,
    5,
    10
  );

-- talisman in backpocket
insert into `armors`
values (
    26,
    2,
    'Arcane',
    'Focus',
    0,
    0,
    10,
    0,
    15
  );

-- mage's focus equipped in neck
insert into `armors`
values (
    27,
    2,
    null,
    'of the Path',
    2,
    10,
    5,
    2,
    0
  );

insert into `armors`
values (
    28,
    2,
    "Sage's",
    null,
    0,
    30,
    0,
    0,
    20
  );

insert into `armors`
values (
    29,
    9,
    'Lucky',
    null,
    0,
    5,
    0,
    0,
    0
  );

insert into `armors`
values (
    33,
    5,
    null,
    'of Fortitude',
    10,
    50,
    0,
    10,
    0
  );

insert into `armors`
values (
    34,
    6,
    'Runed',
    'of Arcana',
    0,
    10,
    0,
    0,
    25
  );

-- feet / shins / gloves
insert into `armors`
values (
    30,
    13,
    null,
    'of Treading',
    12,
    20,
    0,
    5,
    0
  );

-- Iron Boots -> feet slot
insert into `armors`
values (
    31,
    12,
    'Hardened',
    null,
    15,
    40,
    0,
    10,
    0
  );

-- Plate Greaves -> shins
insert into `armors`
values (
    32,
    4,
    'Silken',
    null,
    4,
    5,
    0,
    0,
    10
  );

-- gloves
-- weapons
insert into `weapons`
values (
    40,
    'Rusted',
    'Dagger',
    20,
    0,
    10,
    5,
    0
  );

insert into `weapons`
values (
    41,
    'Wicked',
    'Dirk',
    20,
    0,
    20,
    12,
    0
  );

insert into `weapons`
values (
    42,
    'Ancient',
    'Longsword',
    25,
    0,
    15,
    25,
    0
  );

insert into `weapons`
values (
    43,
    'Runed',
    'Staff',
    20,
    0,
    10,
    0,
    40
  );

-- consumables
-- ex. (Greater/Lesser) Healing Potion, (Powerful/Stunning) Power Potion, Slumpkin Strength Potion, Arcana Spellpower Potion
-- (simple stats stored on armors table when consumables give direct effects as templates)
insert into `consumables`
values (
    35,
    1,
    null,
    'Lesser',
    25,
    0,
    0,
    0
  );

insert into `consumables`
values (
    36,
    2,
    null,
    'Greater',
    75,
    0,
    0,
    0
  );

insert into `consumables`
values (
    37,
    2,
    null,
    'Potent',
    0,
    50,
    0,
    0
  );

insert into `consumables`
values (
    38,
    2,
    null,
    'Shimmering',
    0,
    30,
    0,
    30
  );

insert into `consumables`
values (
    39,
    1,
    null,
    'Cloth',
    10,
    0,
    0,
    0
  );

-- characters
insert into `characters`
values (null, 1, 'Oggnjen', 2, 1, 1, 0, 25, 0, 25);

insert into `characters`
values (
    null,
    1,
    'Thorin',
    15,
    3,
    500,
    50,
    20,
    30,
    0
  );

insert into `characters`
values (
    null,
    1,
    'Mira',
    20,
    1,
    750,
    80,
    30,
    10,
    20
  );

insert into `characters`
values (
    null,
    1,
    'Selene',
    6,
    2,
    1200,
    40,
    60,
    0,
    40
  );

insert into `characters`
values (
    null,
    1,
    'Borin',
    56,
    3,
    300,
    35,
    10,
    25,
    0
  );

insert into `characters`
values (
    null,
    1,
    'Lysa',
    15,
    1,
    950,
    70,
    25,
    5,
    30
  );

insert into `characters`
values (
    null,
    1,
    'Keth',
    500,
    2,
    2000,
    60,
    80,
    0,
    80
  );

-- stackable vendor goods and class-useful trinkets
insert into `item_templates`
values (
    null,
    'Arrow Bundle',
    6,
    1,
    1,
    20,
    1,
    'gdrop',
    '{}'
  );

insert into `item_templates`
values (
    null,
    "Hunter's Token",
    7,
    2,
    0,
    1,
    1,
    'tdrop',
    '{}'
  );

insert into `armors`
values (48, 9, 'Bundle', null, 0, 0, 0, 0, 0);

-- sample vendor variants for existing templates 
insert into `armors`
values (
    1,
    1,
    'Reinforced',
    'of the Owl',
    18,
    12,
    8,
    2,
    25
  );

insert into `armors`
values (
    2,
    7,
    'Woolen',
    'of Comfort',
    10,
    18,
    0,
    0,
    5
  );

insert into `armors`
values (
    3,
    6,
    'Gleaming',
    'of Clarity',
    2,
    8,
    0,
    0,
    12
  );