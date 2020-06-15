drop table if exists users cascade;
drop table if exists posts cascade;
drop table if exists comments;
drop table if exists tags cascade;
drop table if exists post_tag;

create table users (
	id SERIAL PRIMARY KEY,
	email VARCHAR(255) UNIQUE NOT NULL,
	hashed_pwd VARCHAR(255) UNIQUE NOT NULL,
	nickname VARCHAR(255) NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

create table posts (
	id SERIAL PRIMARY KEY,
	title VARCHAR(255) NOT NULL,
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

create table tags(
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL
);

create table post_tag(
	id SERIAL PRIMARY KEY,
	post_id INTEGER REFERENCES posts(id),
	tag_id INTEGER REFERENCES tags(id)
);
