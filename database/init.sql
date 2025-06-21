create table users (
    id serial primary key,
    name varchar(100) not null,
    birthday varchar(100) not null,
    school varchar(100) not null
);

copy users (name, birthday, school) 
from '/docker-entrypoint-initdb.d/Danhsach.csv'
with (format csv, header true);