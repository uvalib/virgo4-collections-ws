BEGIN;
DROP TABLE IF EXISTS fund_codes;
ALTER TABLE collections DROP COLUMN filter_value;
UPDATE collections set filter_name='FilterBookplate' where filter_name='FilterFundCode';
COMMIT;