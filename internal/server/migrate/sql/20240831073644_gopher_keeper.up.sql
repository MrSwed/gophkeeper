CREATE EXTENSION if not exists moddatetime;
CREATE EXTENSION if not exists pgcrypto;

create table users
(
 id          uuid primary key,
 description text,
 email       text,
 packed_key  bytea,
 created_at  timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
 updated_at  timestamptz
);

create index users_packed_key_email_index
 on users (packed_key, email);

CREATE TRIGGER mdt_users
 BEFORE UPDATE
 ON users
 FOR EACH ROW
EXECUTE PROCEDURE moddatetime(updated_at);

create table storage
(
 key         varchar(255)                          not null
  primary key,
 user_id     uuid                                  not null
  constraint "storage_users.id_fk"
   references users,
 description text,
 filename    text,
 blob        bytea,
 created_at  timestamptz default CURRENT_TIMESTAMP not null,
 updated_at  timestamptz
);
create index storage_created_at_index
 on storage (created_at desc);
CREATE TRIGGER mdt_storage
 BEFORE UPDATE
 ON storage
 FOR EACH ROW
EXECUTE PROCEDURE moddatetime(updated_at);

create table clients
(
 token      bytea                                              not null
   primary key,
 user_id    uuid                     default gen_random_uuid() not null
  constraint clients_users_id_fk
   references users,
 meta       json,
 created_at timestamp with time zone default CURRENT_TIMESTAMP not null,
 expired_at timestamp with time zone
);

CREATE OR REPLACE FUNCTION hash_update_tg() RETURNS trigger AS
$$
BEGIN
 IF tg_op = 'INSERT' THEN
  NEW.token = digest(md5(random()::text), 'sha256');
  RETURN NEW;
 END IF;
END;
$$ LANGUAGE plpgsql;

