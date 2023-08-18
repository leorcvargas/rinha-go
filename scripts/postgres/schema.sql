CREATE EXTENSION IF NOT EXISTS "unaccent";

CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE TABLE
    IF NOT EXISTS public.people (
        id uuid PRIMARY KEY NOT NULL,
        nickname varchar(32) NOT NULL,
        "name" varchar(100) NOT NULL,
        birthdate date NOT NULL,
        stack text NULL,
        trgm_q text NOT NULL,
        CONSTRAINT people_nickname_key UNIQUE (nickname)
    );

-- ALTER TABLE public.people

-- ADD

--     COLUMN trgm_q text GENERATED ALWAYS AS (

--         nickname || ' ' || "name" || ' ' || stack

--     ) STORED;


CREATE INDEX
    IF NOT EXISTS CONCURRENTLY idx_people_trigram ON public.people USING gist (trgm_q gist_trgm_ops);