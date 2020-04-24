drop table if exists posts cascade;
drop table if exists comments;

create table posts (
	id SERIAL PRIMARY KEY,
	content TEXT NOT NULL,
	author VARCHAR(255) NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

create table comments(
	id SERIAL PRIMARY KEY,
	content TEXT NOT NULL,
	author VARCHAR(255) NOT NULL,
	post_id INTEGER REFERENCES posts(id),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)
