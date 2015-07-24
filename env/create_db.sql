DROP TABLE IF EXISTS users;

CREATE TABLE users(
       id SERIAL PRIMARY KEY,
       username VARCHAR UNIQUE,
       password VARCHAR       
);

DROP TABLE IF EXISTS configuration;
CREATE TABLE configuration(
       id SERIAL PRIMARY KEY,
       user_id INT,
       config_name VARCHAR,
       host_name VARCHAR,
       username VARCHAR,
       FOREIGN KEY (user_id) REFERENCES users(id)
);
