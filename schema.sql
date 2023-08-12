-- Tabela de pessoas

CREATE TABLE
    IF NOT EXISTS public.people (
        id uuid PRIMARY KEY NOT NULL,
        nickname varchar(32) NOT NULL,
        "name" varchar(100) NOT NULL,
        birthdate date NULL,
        stack text NULL,
        CONSTRAINT people_nickname_key UNIQUE (nickname)
    );

ALTER TABLE public.people
ADD
    COLUMN fts_q tsvector GENERATED ALWAYS AS (
        to_tsvector(
            'english',
            nickname || ' ' || "name" || ' ' || stack
        )
    ) STORED;

CREATE INDEX people_fts_q_idx ON public.people USING gin (fts_q);

-- -- Index para pesquisa de texto

-- CREATE INDEX

--     IF NOT EXISTS people_search_aggr ON public.people (search_aggr);