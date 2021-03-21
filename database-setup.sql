set setts.usr to :username;
set setts.dev to :dev;
set setts.db to :db;

create database :db;
create user :username;

grant all privileges on database :db to :username;

\if :dev
    create database test_keystore;
    grant all privileges on database test_keystore to :username;
\endif