BEGIN;
UPDATE collections set active=true where id < 86;
COMMIT;