#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  DROP TEXT SEARCH DICTIONARY IF EXISTS dict_simple_ptbr CASCADE;
  CREATE TEXT SEARCH DICTIONARY public.dict_simple_ptbr (
      TEMPLATE = pg_catalog.simple,
      STOPWORDS = portuguese,
      Accept = false
  );

  DROP TEXT SEARCH DICTIONARY IF EXISTS dict_ispell_ptbr CASCADE;
  CREATE TEXT SEARCH DICTIONARY dict_ispell_ptbr (
    TEMPLATE = ispell,
    DictFile = brazilian,
    AffFile = brazilian,
    StopWords = portuguese
  );
  DROP TEXT SEARCH DICTIONARY IF EXISTS dict_snowball_ptbr CASCADE;
  CREATE TEXT SEARCH DICTIONARY dict_snowball_ptbr (
    TEMPLATE = snowball,
    Language = portuguese,
    StopWords = portuguese
  );

  CREATE EXTENSION IF NOT EXISTS "unaccent";
  CREATE EXTENSION IF NOT EXISTS "pg_trgm";
  CREATE EXTENSION IF NOT EXISTS "fuzzystrmatch";

  DROP TEXT SEARCH CONFIGURATION IF EXISTS people_terms CASCADE;
  CREATE TEXT SEARCH CONFIGURATION people_terms (COPY=pg_catalog.portuguese);
  ALTER TEXT SEARCH CONFIGURATION people_terms ALTER MAPPING FOR asciiword, asciihword, hword_asciipart, word, hword, hword_part
    WITH unaccent, portuguese_stem, dict_simple_ptbr, dict_ispell_ptbr, dict_snowball_ptbr;

  CREATE TABLE IF NOT EXISTS public.people (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    nickname varchar(32) NOT NULL,
    "name" varchar(100) NOT NULL,
    birthdate date NOT NULL,
    stack _text NULL,
    terms tsvector,
    CONSTRAINT people_pkey PRIMARY KEY (id)
  );
  CREATE INDEX idx_people_terms ON people USING GIN(terms);

  CREATE OR REPLACE FUNCTION update_people_terms() RETURNS trigger AS 
  \$\$
  BEGIN
    new.terms := 
      setweight(to_tsvector('public.people_terms', unaccent(coalesce(new.name, ''))), 'A') ||
      setweight(to_tsvector('public.people_terms', unaccent(coalesce(new.nickname, ''))), 'B') ||
      setweight(to_tsvector('public.people_terms', unaccent(coalesce(array_to_string(new.stack, ','), ''))), 'C');
    RETURN new;
  END
  \$\$ LANGUAGE plpgsql;
  DROP TRIGGER IF EXISTS tg_update_people_terms ON people;
  CREATE TRIGGER tg_update_people_terms BEFORE INSERT OR UPDATE ON people FOR EACH ROW EXECUTE PROCEDURE update_people_terms();

  CREATE OR REPLACE FUNCTION people_terms_tsquery(word text) RETURNS tsquery AS 
  \$\$
  BEGIN
    RETURN plainto_tsquery('public.people_terms', unaccent(trim(word)));
  END
  \$\$ LANGUAGE plpgsql;
EOSQL
