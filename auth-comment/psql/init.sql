create table users( user_id SERIAL, username VARCHAR(255) UNIQUE NOT NULL, PRIMARY KEY(username));

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
create table posts(postid uuid DEFAULT uuid_generate_v4 (),postname VARCHAR(255) NOT NULL,
 createdby VARCHAR(255) NOT NULL, PRIMARY KEY (postid), UNIQUE (postname,createdby),
 CONSTRAINT fk_user FOREIGN KEY(createdby) REFERENCES users(username) on DELETE CASCADE);

CREATE EXTENSION ltree;
create TABLE comments(postid uuid DEFAULT uuid_generate_v4 (),
                comment TEXT,
                comment_reaction VARCHAR(255) NOT NULL,
                commented_user VARCHAR(255) NOT NULL,
                created_at TIMESTAMP,
                parent_path ltree UNIQUE NOT NULL,
                CONSTRAINT fk_post FOREIGN KEY(postid)
                REFERENCES posts(postid) on DELETE CASCADE);
CREATE INDEX section_parent_path_idx ON comments USING GIST (parent_path);