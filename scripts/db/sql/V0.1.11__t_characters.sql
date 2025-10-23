create table characters (
  character_id SERIAL,
  user_id int not null,
  character_name varchar(16) not null unique,
  hero_id int not null,
  
  level int not null,
  exp int not null,
  gold int not null,
  tokens int not null,
  points_health int not null,
  points_power int not null,
  points_strength int not null,
  points_spellpower int not null,
  unspent_points int not null default 1,

  foreign key (hero_id) references heroes (id),
  foreign key (user_id) references users (id) on delete cascade,
  primary key(character_id)
);


create index idx_characters_user_id on characters (user_id);
create index idx_characters_hero_id on characters (hero_id);

-- alter table characters replica identity full;