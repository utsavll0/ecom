CREATE TABLE IF NOT EXISTS public.order_items
(
    "id"        SERIAL,
    "orderId"   INT            NOT NULL,
    "productId" INT            NOT NULL,
    "quantity"  INT            NOT NULL,
    "price"     DECIMAL(10, 2) NOT NULL,

    primary key (id),
    FOREIGN KEY ("orderId") REFERENCES public.orders ("id"),
    FOREIGN KEY ("productId") REFERENCES public.products ("id")
);