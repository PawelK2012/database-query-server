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

INSERT INTO public.usersTest (id,firt_name,last_name, email, password, is_admin, created_at, updated_at) VALUES (1, 'TestUser1', 'Test surname', 'test@email.com', 'pass1', false, '2004-10-19 10:23:54', '2004-10-19 10:23:54'), (2, 'TestUser2', 'Surname 2', 'test2@email.com', 'passx', false, '2004-10-19 10:23:54', '2004-10-19 10:23:54');