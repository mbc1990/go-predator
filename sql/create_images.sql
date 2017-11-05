CREATE TABLE images(
  image_id serial primary key not null,
  filename varchar(1024) not null,
  original_url varchar(1024) not null,
  source varchar(1024) not null,
  source_id varchar(1024) not null,
  classified boolean not null,
  timestamp timestamp default current_timestamp
);
