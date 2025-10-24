create table hero_abilities (
    hero_id int not null,
    ability_id int not null,

    foreign key (hero_id) references heroes (id) on delete cascade,
    foreign key (ability_id) references abilities (id),
    primary key (hero_id, ability_id)
);