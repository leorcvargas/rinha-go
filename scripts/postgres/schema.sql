CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE TABLE
    IF NOT EXISTS public.people (
        id uuid NOT NULL,
        nickname varchar(32) PRIMARY KEY,
        "name" varchar(100) NOT NULL,
        birthdate date NOT NULL,
        stack text NULL,
        search text NOT NULL
    );

CREATE INDEX
    CONCURRENTLY IF NOT EXISTS idx_people_trigram ON public.people USING GIST (search gist_trgm_ops);