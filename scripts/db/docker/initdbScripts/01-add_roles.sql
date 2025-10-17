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

GRANT pg_read_all_data TO inventoryuser, charuser, gambauser, chatuser, storeuser;
GRANT pg_write_all_data TO inventoryuser, charuser, gambauser, chatuser, storeuser;

SELECT pg_reload_conf();