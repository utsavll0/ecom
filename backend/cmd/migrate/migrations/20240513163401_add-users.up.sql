CREATE TABLE IF NOT EXISTS public.users
(
    "id"        SERIAL,
    "firstname" VARCHAR(255) NOT NULL,
    "lastname"  VARCHAR(255) NOT NULL,
    "email"     VARCHAR(255) NOT NULL,
    "password"  VARCHAR(255) NOT NULL,
    "createdAt" TIMESTAMP    NOT NULL DEFAULT current_timestamp,

    PRIMARY KEY (id),
    UNIQUE (email)

);