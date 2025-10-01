create table if not exists item_template_stats (
  item_template_id int not null,
  stat_type_id int not null,
  value int not null,
  foreign key (item_template_id) references item_templates(id) on delete cascade,
  foreign key (stat_type_id) references stat_types(id),
  primary key (item_template_id, stat_type_id)
);