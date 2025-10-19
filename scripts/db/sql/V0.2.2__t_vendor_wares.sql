create table vendor_wares (
  vendor_id int not null,
  item_type_id int not null,
  
  foreign key (vendor_id) references vendors (id),
  foreign key (item_type_id) references item_types (id),
  primary key (vendor_id, item_type_id)
);

-- alter table vendor_wares replica identity full;