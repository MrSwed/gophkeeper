create table storage
(
 key         TEXT not null
  constraint storage_pk
   primary key,
 description TEXT     DEFAULT '' not null,
 created_at  datetime DEFAULT (datetime('now', 'localtime')),
 updated_at  datetime,
 filename    text,
 blob        BLOB
);

PRAGMA case_sensitive_like=OFF;
