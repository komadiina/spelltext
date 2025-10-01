insert into heroes
values (
    DEFAULT,
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

insert into heroes
values (
    DEFAULT,
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

insert into heroes
values (
    DEFAULT,
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

insert into stat_types
VALUES (DEFAULT, 'ST01', 'Strength');

insert into stat_types
VALUES (DEFAULT, 'ST02', 'Spellpower');

insert into stat_types
VALUES (DEFAULT, 'ST03', 'Health');

insert into stat_types
values (DEFAULT, 'ST04', 'Power');

insert into item_types
values (DEFAULT, 'IT01', 'Armor');

insert into item_types
values (DEFAULT, 'IT02', 'Consumable');

insert into item_types
values (DEFAULT, 'IT03', 'Weapon');

insert into item_types
values (DEFAULT, 'IT04', 'Vanity');

insert into item_types
values (DEFAULT, 'IT05', 'Wearable vanity');

insert into item_types
values (DEFAULT, 'IT06', 'Gold drop');

insert into item_types
values (DEFAULT, 'IT07', 'Token drop');

insert into equip_slots
values (DEFAULT, 'ES01', 'Head');

insert into equip_slots
values (DEFAULT, 'ES02', 'Neck');

insert into equip_slots
values (DEFAULT, 'ES03', 'Shoulders');

insert into equip_slots
values (DEFAULT, 'ES04', 'Gloves');

insert into equip_slots
values (DEFAULT, 'ES05', 'Ring');

insert into equip_slots
values (DEFAULT, 'ES06', 'Ring');

insert into equip_slots
values (DEFAULT, 'ES07', 'Chest');

insert into equip_slots
values (DEFAULT, 'ES08', 'Waist');

insert into equip_slots
values (DEFAULT, 'ES09', 'Backpocket');

insert into equip_slots
values (DEFAULT, 'ES10', 'Back');

insert into equip_slots
values (DEFAULT, 'ES11', 'Legs');

insert into equip_slots
values (DEFAULT, 'ES12', 'Shins');

insert into equip_slots
values (DEFAULT, 'ES13', 'Feet');

insert into item_templates
values (
    DEFAULT,
    'Peasant Helmet',
    1,
    1,
    0,
    1,
    1,
    'helmet',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Peasant Vest',
    1,
    1,
    0,
    1,
    1,
    'vest',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Wedding ring',
    1,
    1,
    0,
    1,
    1,
    'ring',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Health Potion',
    2,
    1,
    1,
    5,
    1,
    'hpot',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Power Potion',
    2,
    1,
    1,
    5,
    1,
    'ppot',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Strength Potion',
    2,
    2,
    1,
    5,
    1,
    'spot',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Spellpower Potion',
    2,
    2,
    1,
    5,
    1,
    'sppot',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Stoneblade',
    3,
    1,
    0,
    1,
    1,
    'weapon',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Warstaff',
    3,
    1,
    0,
    1,
    1,
    'weapon',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Sturdy shovel',
    3,
    1,
    0,
    1,
    1,
    'weapon',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Ygmirs head',
    4,
    2,
    0,
    1,
    1,
    'vanity',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Lucky? charm',
    5,
    1,
    0,
    1,
    1,
    'wvanity',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Honest gold bag',
    6,
    1,
    0,
    1,
    1,
    'gdrop',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Thief sack',
    6,
    2,
    0,
    1,
    1,
    'gdrop',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Bankrobber sack',
    6,
    3,
    0,
    1,
    1,
    'gdrop',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Gold bar',
    6,
    4,
    0,
    1,
    1,
    'gdrop',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Rugged token',
    7,
    1,
    0,
    1,
    1,
    'tdrop',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Polished token',
    7,
    2,
    0,
    1,
    1,
    'tdrop',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Headhunter token',
    7,
    3,
    0,
    1,
    1,
    'tdrop',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Worldslayer token',
    7,
    4,
    0,
    1,
    1,
    'tdrop',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Traveler Cloak',
    1,
    1,
    0,
    1,
    1,
    'cloak',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Cloak of Ember',
    1,
    2,
    0,
    1,
    1,
    'cloak',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Moonlit Shawl',
    1,
    3,
    0,
    1,
    1,
    'cloak',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Hunters Quiver',
    1,
    1,
    0,
    1,
    1,
    'backpocket',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Druidic Talisman',
    1,
    2,
    0,
    1,
    1,
    'trinket',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Mage Focus',
    1,
    2,
    0,
    1,
    1,
    'trinket',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Adventurer Necklace',
    1,
    1,
    0,
    1,
    1,
    'necklace',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Sage Locket',
    1,
    3,
    0,
    1,
    1,
    'necklace',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Lucky Rabbit Foot',
    5,
    1,
    0,
    1,
    1,
    'wvanity',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Iron Boots',
    1,
    1,
    0,
    1,
    1,
    'boots',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Plate Greaves',
    1,
    2,
    0,
    1,
    1,
    'shins',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Silken Gloves',
    1,
    2,
    0,
    1,
    1,
    'gloves',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Band of Fortitude',
    1,
    3,
    0,
    1,
    1,
    'ring',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Ring of Arcana',
    1,
    3,
    0,
    1,
    1,
    'ring',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Lesser Healing Potion',
    2,
    1,
    1,
    5,
    1,
    'hpot',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Greater Healing Potion',
    2,
    2,
    1,
    3,
    1,
    'hpot',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Greater Power Potion',
    2,
    2,
    1,
    3,
    1,
    'ppot',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Mana Crystal',
    2,
    2,
    1,
    5,
    1,
    'mcrystal',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Bandage',
    2,
    1,
    1,
    10,
    1,
    'bandage',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Rusted Dagger',
    3,
    1,
    0,
    1,
    1,
    'weapon',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Wicked Dirk',
    3,
    2,
    0,
    1,
    1,
    'weapon',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Ancient Longsword',
    3,
    3,
    0,
    1,
    1,
    'weapon',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Runed Staff',
    3,
    3,
    0,
    1,
    1,
    'weapon',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Shadow Cloak (Vanity)',
    4,
    1,
    0,
    1,
    1,
    'vanity',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Bronze Token',
    7,
    1,
    0,
    1,
    1,
    'tdrop',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    'Emerald Token',
    7,
    3,
    0,
    1,
    1,
    'tdrop',
    '{}'
  );

insert into armors
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

insert into armors
values (
    22,
    10,
    DEFAULT,
    'of Ember',
    8,
    10,
    0,
    0,
    10
  );

insert into armors
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

insert into armors
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

insert into armors
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

insert into armors
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

insert into armors
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

insert into armors
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

insert into armors
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

insert into armors
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

insert into armors
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

insert into armors
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

insert into armors
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

insert into armors
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

insert into armors
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

insert into armors
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

insert into armors
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

insert into armors
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

insert into armors
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

insert into armors
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

insert into armors
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

insert into weapons
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

insert into weapons
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

insert into weapons
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

insert into weapons
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

insert into consumables
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

insert into consumables
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

insert into consumables
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

insert into consumables
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

insert into consumables
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

insert into characters
values (DEFAULT, 1, 'Oggnjen', 2, 1, 1, 0, 25, 0, 25);

insert into characters
values (
    DEFAULT,
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

insert into characters
values (
    DEFAULT,
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

insert into characters
values (
    DEFAULT,
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

insert into characters
values (
    DEFAULT,
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

insert into characters
values (
    DEFAULT,
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

insert into characters
values (
    DEFAULT,
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

insert into item_templates
values (
    DEFAULT,
    'Arrow Bundle',
    6,
    1,
    1,
    20,
    1,
    'gdrop',
    '{}'
  );

insert into item_templates
values (
    DEFAULT,
    "Hunter's Token",
    7,
    2,
    0,
    1,
    1,
    'tdrop',
    '{}'
  );

insert into armors
values (48, 9, 'Bundle', null, 0, 0, 0, 0, 0);

insert into armors
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

insert into armors
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

insert into armors
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