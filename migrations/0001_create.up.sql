create table users
(
    id          uuid not null
        constraint users_pk
            primary key,
    email       text not null,
    phone       text not null,
    portal_code varchar,
    admin       boolean
);

alter table users
    owner to http;

create table counterparty
(
    id      uuid    not null
        constraint counterparty_pk
            primary key,
    name    text    not null,
    inn     text    not null,
    kpp     text    not null,
    ogrn    text    not null,
    address text    not null,
    type    integer not null
);

alter table counterparty
    owner to http;