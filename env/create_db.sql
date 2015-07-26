DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS configurations CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;
CREATE TABLE users(
       id SERIAL PRIMARY KEY,
       username VARCHAR UNIQUE,
       password VARCHAR       
);


CREATE TABLE configurations(
       id SERIAL PRIMARY KEY,
       config_name VARCHAR UNIQUE,
       host_name VARCHAR,
       port INT, 
       username VARCHAR
);


create table sessions(
       session_id VARCHAR,
       user_id INT,
       FOREIGN KEY (user_id) REFERENCES users(id)
);

