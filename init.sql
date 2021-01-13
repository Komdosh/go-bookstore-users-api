create user gousr with password 'gousr';

create database gomicroservices
    with owner gousr;


-- switch to gomicroservices before run next lines

create schema users_db;

create table users_db.users
(
    id           bigserial    not null
        constraint users_pk
            primary key,
    first_name   varchar(255),
    last_name    varchar(255),
    email        varchar(255) not null,
    date_created timestamp,
    password     varchar(255)  not null,
    status       varchar(255) not null
);

alter table users_db.users
    owner to gousr;

create unique index users_id_uindex
    on users_db.users (id);

create unique index users_email_uindex
    on users_db.users (email);
