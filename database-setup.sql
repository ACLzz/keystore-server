create user :username WITH ENCRYPTED PASSWORD :'password';
create database :db with OWNER = :username;

\if :dev
    create database test_keystore with OWNER = :username;
\endif