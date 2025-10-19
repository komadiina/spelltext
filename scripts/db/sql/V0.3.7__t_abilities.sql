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
    spellpower_multiplier float not null default 1.0 
);

create index idx_abilities_id on abilities (id);
create index idx_abilities_type on abilities (type);

-- alter table abilities replica identity full;