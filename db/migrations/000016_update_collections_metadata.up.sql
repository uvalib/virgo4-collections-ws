BEGIN;
UPDATE collections set
   description = 'The Daily Progress is the Charlottesville, VA, area newspaper, published daily from 1892 to the present. Issues from 1892 through 1964 have been digitized from the Library''s set of microfilm and are available for viewing online. This digital edition has been reviewed for scanning quality and we believe the digital images are as clear as they can be, given the variable condition of the originating microfilm. <a href=''https://search.lib.virginia.edu/sources/uva_library/items/u1870648''>The microfilm is available for request</a>.'
   where id = 1;

UPDATE collections set item_label='Item' where id > 1;
COMMIT;