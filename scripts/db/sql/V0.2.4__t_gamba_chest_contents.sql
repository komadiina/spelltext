create table gamba_chest_contents (
    gamba_chest_id int not null,
    item_id int not null,
    
    foreign key (gamba_chest_id) references gamba_chests (id),
    foreign key (item_id) references items (id),
    primary key (gamba_chest_id, item_id)
);

alter table gamba_chest_contents replica identity full;