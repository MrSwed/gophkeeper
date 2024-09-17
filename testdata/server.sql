insert into users (id, email, password)
values ('be341e38-b8a9-4230-af77-fcc34c9f2e13', 'example@example.com', '$2a$10$Xlq4avgrTER5aAAJhL4HAu4WSEgGxx6vuzNoYcO3UflbLGzszMmY6');

insert into public.clients (token, user_id, meta, created_at, expired_at)
values  (E'\\x8CA0C5A18320FC2F264CFA95639EA27888727C6090D6F9CB0D6C5798A93FCB63', 'be341e38-b8a9-4230-af77-fcc34c9f2e13', '{"Addr":{"IP":"127.0.0.1","Port":59020,"Zone":""},"LocalAddr":{"IP":"127.0.0.1","Port":30063,"Zone":""},"AuthInfo":null}', '2024-09-17 18:38:50.799916 +00:00', '2025-09-17 18:38:50.680984 +00:00');