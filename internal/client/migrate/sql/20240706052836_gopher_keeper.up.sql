create table storage
(
 key         TEXT not null
  constraint storage_pk
   primary key,
 description TEXT,
 created_at  integer DEFAULT CURRENT_TIMESTAMP,
 updated_at  integer DEFAULT CURRENT_TIMESTAMP,
 filename    text,
 blob        BLOB
);

PRAGMA case_sensitive_like=OFF;
