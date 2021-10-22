BEGIN;
ALTER TABLE collections ADD COLUMN active boolean  NOT NULL default true;

update collections set active=false where id > 85;

COMMIT;