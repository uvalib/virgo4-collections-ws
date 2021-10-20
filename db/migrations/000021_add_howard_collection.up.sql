BEGIN;

  INSERT into collections (id, title, description, filter_name, item_label)
  VALUES (90,
      'James Murray Howard University of Virginia Historic Buildings and Grounds Collection, University of Virginia Library',
      'James Murray Howard (1948-2008) was the Architect for Historic Buildings and Grounds at the University of Virginia from 1982 to 2002. He supervised a comprehensive restoration program for the Academical Village, overseeing numerous restoration projects over the course of his career. This image collection is his personal archive of his researches into the history of the creation of Jefferson’s buildings, of techniques and processes used in their construction, decoration and restoration, and of his teaching career directing what he called "a practical working laboratory for University students."',
      'FilterDigitalCollection',
      'Image');

INSERT into collection_features (collection_id, feature_id) values (90, 1);


INSERT into images (collection_id, title, width, height, filename) values (90, 'University of Virginia', 450, 150, 'uva.png');

LOCK TABLE collections IN EXCLUSIVE MODE;
SELECT setval('collections_id_seq', COALESCE((SELECT MAX(id)+1 FROM collections), 1), false);

COMMIT;