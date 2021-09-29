truncate TABLE collection_features CASCADE;
truncate TABLE collections CASCADE;

  INSERT into collections (id, key, type, description, filter_name, filter_value, start_date, end_date)
  VALUES (1,
      'dpdm',
      'Collection',
      'The Daily Progress is the Charlottesville, VA, area newspaper, published daily from 1892 to the present. Issues from 1892 through 1964 have been digitized from the Library''s set of microfilm and are available for viewing online. This digital edition has been reviewed for scanning quality and we believe the digital images are as clear as they can be, given the variable condition of the originating microfilm. The microfilm is available for viewing in Alderman Library.',
      'FilterDigitalCollection',
      'Daily Progress Digitized Microfilm',
      '1893-03-16',
      '1964-12-31');

INSERT into collection_features (collection_id, feature_id) values (1, 1), (1, 2), (1, 4), (1, 4);