CREATE TABLE
    IF NOT EXISTS public.people (
        id uuid PRIMARY KEY NOT NULL,
        nickname varchar(32) NOT NULL,
        "name" varchar(100) NOT NULL,
        birthdate date NULL,
        stack text NULL,
        CONSTRAINT people_nickname_key UNIQUE (nickname)
    );