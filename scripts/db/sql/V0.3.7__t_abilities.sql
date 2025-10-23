create table abilities (
    id serial primary key,
    name varchar(64) not null,
    description text not null,
    type varchar(16) not null, -- 'passive', 'defensive', 'offensive'
    talent_point_cost int not null default 1,

    power_cost int not null,
    base_damage int not null,

    -- scaling formula: base*(player_str*str_mult + player_sp*sp_mult)
    strength_multiplier float not null default 1.0,
    spellpower_multiplier float not null default 1.0,

    str_mult_pl float not null default 0.1, -- 0.1 = 10% increase in multi per level
    sp_mult_pl float not null default 0.1, -- 0.1 = 10% increase in multi per level
    min_level int not null default 1
);

create index idx_abilities_type on abilities (type);