CREATE USER inventoryuser WITH LOGIN PASSWORD 'changeme';
GRANT ALL PRIVILEGES ON SCHEMA public TO inventoryuser;

CREATE USER charuser WITH LOGIN PASSWORD 'changeme';
GRANT ALL PRIVILEGES ON SCHEMA public TO charuser;

CREATE USER gambauser WITH LOGIN PASSWORD 'changeme';
GRANT ALL PRIVILEGES ON SCHEMA public TO gambauser;

CREATE USER chatuser WITH LOGIN PASSWORD 'changeme';
GRANT ALL PRIVILEGES ON SCHEMA public TO chatuser;

CREATE USER storeuser WITH LOGIN PASSWORD 'changeme';
GRANT ALL PRIVILEGES ON SCHEMA public TO storeuser;

CREATE USER authuser WITH LOGIN PASSWORD 'changeme';
GRANT ALL PRIVILEGES ON SCHEMA public TO authuser;

CREATE USER combatuser WITH LOGIN PASSWORD 'changeme';
GRANT ALL PRIVILEGES ON SCHEMA public TO combatuser;

CREATE USER builduser WITH LOGIN PASSWORD 'changeme';
GRANT ALL PRIVILEGES ON SCHEMA public TO builduser;

GRANT pg_read_all_data TO inventoryuser, charuser, gambauser, chatuser, storeuser, authuser, combatuser, builduser;
GRANT pg_write_all_data TO inventoryuser, charuser, gambauser, chatuser, storeuser, authuser, combatuser, builduser;

SELECT pg_reload_conf();