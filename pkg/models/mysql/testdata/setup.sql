-- auto-generated definition
create table snippets
(
    id      int auto_increment primary key,
    title   varchar(100) not null,
    content text         not null,
    created datetime     not null,
    expires datetime     not null
);

create index idx_snippets_created
    on snippets (created);

CREATE TABLE users
(
    id              INTEGER      NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name            VARCHAR(255) NOT NULL,
    email           VARCHAR(255) NOT NULL,
    hashed_password CHAR(60)     NOT NULL,
    created         DATETIME     NOT NULL
);
ALTER TABLE users
    ADD CONSTRAINT users_uc_email UNIQUE (email);
INSERT INTO users (name, email, hashed_password, created)
VALUES ('Alice Jones',
        'alice@example.com',
        '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
        '2018-12-23 17:25:22');