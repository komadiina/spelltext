create table gamba_chests (
    id serial primary key,
    name varchar(64) not null,
    description varchar(255) not null,
    price int not null
)