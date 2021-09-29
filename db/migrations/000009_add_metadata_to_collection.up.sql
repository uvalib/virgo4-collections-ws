BEGIN;
ALTER TABLE collections ADD COLUMN item_label VARCHAR (80) default '';
UPDATE collections set item_label='Issue' where id=1;
UPDATE collections set item_label='Book' where id>1;
COMMIT;