create table if not exists article
(
    id      serial,
    author  text,
    title   text,
    body    text,
    created timestamp not null
);

alter table article
    owner to postgres;

