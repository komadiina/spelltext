CREATE TABLE heroes (
  id SERIAL,
  name VARCHAR(64) NOT NULL,
  base_health INT NOT NULL,
  base_power INT NOT NULL,
  base_strength INT NOT NULL,
  base_spellpower INT NOT NULL,
  health_pl int not null default 0,
  power_pl int not null default 0,
  strength_pl int not null default 0,
  spellpower_pl int not null default 0,
  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_heroes_name ON heroes (name);