CREATE TABLE public.usersTest (
    id integer NOT NULL,
    firt_name character varying(255),
    last_name character varying(255),
    email character varying(255),
    password character varying(60),
    is_admin boolean DEFAULT false,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);