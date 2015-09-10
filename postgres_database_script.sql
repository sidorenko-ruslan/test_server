CREATE TABLE users
(
  id serial NOT NULL,
  first_name text,
  last_name text,
  patronymic text,
  CONSTRAINT users_pkey PRIMARY KEY (id)
)
WITH (
  OIDS=FALSE
);
ALTER TABLE users
  OWNER TO postgres;