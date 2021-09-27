CREATE TABLE IF NOT EXISTS collections (
   id serial PRIMARY KEY,
   key VARCHAR (80) UNIQUE NOT NULL,
   type VARCHAR (80) NOT NULL,
   description text,
   filter_name VARCHAR (80) NOT NULL,
   filter_value VARCHAR (80) UNIQUE NOT NULL,
   start_date VARCHAR (20),
   end_date VARCHAR (20)
);

CREATE TABLE IF NOT EXISTS features (
   id serial PRIMARY KEY,
   name varchar(80) UNIQUE NOT NULL
);

INSERT into features (name) values('search_within'), ('full_page_view'), ('calendar_navigation'), ('sequential_navigation');


CREATE TABLE IF NOT EXISTS collection_features (
   id serial PRIMARY KEY,
   collection_id integer NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
   feature_id integer NOT NULL REFERENCES features(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS images (
   id serial PRIMARY KEY,
   collection_id integer NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
   alt_text VARCHAR (255),
   title VARCHAR (255),
   width integer NOT NULL,
   height integer NOT NULL,
   filename VARCHAR (255) NOT NULL
);

