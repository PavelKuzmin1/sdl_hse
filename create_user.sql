create database sdl;

create role "user" login password 'user';
grant usage on schema public to "user";