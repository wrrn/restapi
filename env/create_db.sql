DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS configuration CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;
CREATE TABLE users(
       id SERIAL PRIMARY KEY,
       username VARCHAR UNIQUE,
       password VARCHAR       
);


CREATE TABLE configuration(
       id SERIAL PRIMARY KEY,
       user_id INT,
       config_name VARCHAR,
       host_name VARCHAR,
       username VARCHAR,
       FOREIGN KEY (user_id) REFERENCES users(id)
);


create table sessions(
       session_id VARCHAR,
       user_id INT,
       FOREIGN KEY (user_id) REFERENCES users(id)
);

