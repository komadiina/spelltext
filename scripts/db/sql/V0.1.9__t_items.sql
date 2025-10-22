create table items (
    id serial primary key,
    prefix varchar(64) not null,
    suffix varchar(64) not null,
    item_template_id int not null,
    health int not null,
    power int not null,
    strength int not null,
    spellpower int not null,
    bonus_damage int not null,
    bonus_armor int not null,
    foreign key (item_template_id) references item_templates (id)
);

create index idx_items_item_template_id on items (item_template_id);

-- alter table items replica identity full;