BEGIN;
ALTER TABLE collections DROP COLUMN type;
ALTER TABLE collections DROP COLUMN key;
update collections set filter_value='Daily Progress Digitized Microfilm', filter_name='FilterDigitalCollection' where id=1;
COMMIT;