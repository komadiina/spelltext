create table consumable_templates (
    id SERIAL,
    name varchar(255) not null,
    stackable smallint not null default 0,
    gold_price int not null default 0,
    primary key (id)
)