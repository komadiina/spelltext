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

  foreign key (hero_id) references heroes (id),
  foreign key (user_id) references users (id),
  primary key(character_id)
);


create unique index idx_characters_id on characters (character_id);
create index idx_characters_user_id on characters (user_id);