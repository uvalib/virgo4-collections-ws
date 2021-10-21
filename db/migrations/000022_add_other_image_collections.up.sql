BEGIN;

  INSERT into collections (id,title, description, filter_name, item_label)
  VALUES (91,
      'Holsinger Studio (Charlottesville, Va.)',
      'The Holsinger Studio Collection constitutes a unique photographic record of life in Charlottesville and Albemarle County, Virginia, from before the turn of the century through World War I. The collection consists of approximately 9,000 dry-plate glass negatives and 500 celluloid negatives from the commercial studio of Rufus W. Holsinger. Approximately two-thirds of the collection are studio portraits, and among these are nearly 500 portraits of African-American citizens of Charlottesville and the surrounding area. Many of the portraits are unidentified, but some are of visiting celebrities and dignitaries. In 1915 while teaching art at the University of Virginia, Georgia O’Keeffe sat for a portrait with Mr. Holsinger.',
      'FilterSubject',
      'Image');

INSERT into collection_features (collection_id, feature_id) values (91, 1);

INSERT into collections (id,title, description, filter_name, item_label)
  VALUES (92,
      'Jackson Davis Collection of African American Photographs',
      'The Jackson Davis Collection consists of papers and photographs, given to Special Collections, in 1948 by Helen Mansfield Lynch and Ruth Elizabeth Langhorne, Davis’s daughters, and supplemented in 1999 by Sally Guy Browne and Helen Langhorne, Davis’s granddaughters.',
      'FilterDigitalCollection',
      'Image');

INSERT into collection_features (collection_id, feature_id) values (92, 1);

LOCK TABLE collections IN EXCLUSIVE MODE;
SELECT setval('collections_id_seq', COALESCE((SELECT MAX(id)+1 FROM collections), 1), false);

COMMIT;