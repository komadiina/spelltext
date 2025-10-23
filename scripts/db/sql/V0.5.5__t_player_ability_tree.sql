create table player_ability_tree (
    character_id int not null,
    ability_id int not null,
    level int not null default 1,
    
    foreign key (character_id) references characters (character_id) on delete cascade,
    foreign key (ability_id) references abilities (id),
    primary key (character_id, ability_id)
)