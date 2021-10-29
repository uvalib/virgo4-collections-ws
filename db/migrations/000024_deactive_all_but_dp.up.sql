BEGIN;
UPDATE collections set active=false where id > 1;
COMMIT;