drop table if exists users cascade;
drop table if exists auth cascade;
drop table if exists posts cascade;
drop table if exists comments;

create table users (
	id SERIAL PRIMARY KEY,
	email VARCHAR(255) UNIQUE NOT NULL,
	hashed_pwd VARCHAR(255) UNIQUE NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

create table auth (
	id SERIAL PRIMARY KEY,
	token VARCHAR(255) UNIQUE NOT NULL,
	user_id INTEGER REFERENCES users(id),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

create table posts (
	id SERIAL PRIMARY KEY,
	content TEXT NOT NULL,
	user_id INTEGER REFERENCES users(id),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

create table comments(
	id SERIAL PRIMARY KEY,
	content TEXT NOT NULL,
	post_id INTEGER REFERENCES posts(id),
	user_id INTEGER REFERENCES users(id),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
