CREATE TYPE order_status AS ENUM ('pending', 'completed', 'cancelled');

CREATE TABLE IF NOT EXISTS public.orders
(
    "id"        SERIAL,
    "userId"    INT            NOT NULL,
    "total"     DECIMAL(10, 2) NOT NULL,
    "status"    order_status   NOT NULL DEFAULT 'pending',
    "address"   TEXT           NOT NULL,
    "createdAt" TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,

    primary key (id),
    foreign key ("userId") REFERENCES public.users ("id")

);