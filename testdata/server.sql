insert into users (id, email, password, created_at, updated_at, packed_key)
values
 ('be341e38-b8a9-4230-af77-fcc34c9f2e13', 'example1@example.com', '$2a$10$Xlq4avgrTER5aAAJhL4HAu4WSEgGxx6vuzNoYcO3UflbLGzszMmY6', '2024-09-17 12:00:00 +03:00', null, null),
 ('62911ad4-883f-4b5a-929d-c4a5766560f8', 'example2@example.com', '$2a$10$Xlq4avgrTER5aAAJhL4HAu4WSEgGxx6vuzNoYcO3UflbLGzszMmY6', '2024-09-17 12:00:00 +03:00', null, null),
 ('d581d082-4b74-4dcf-8db3-cbb6e9a2f996', 'example3@example.com', '$2a$10$Xlq4avgrTER5aAAJhL4HAu4WSEgGxx6vuzNoYcO3UflbLGzszMmY6', '2024-09-17 12:00:00 +03:00', '2024-09-17 12:50:00 +03:00', 'predefined packed data'),
 ('afa37fce-557f-4b65-afdf-11f54cebe07a', 'example4@example.com', '$2a$10$Xlq4avgrTER5aAAJhL4HAu4WSEgGxx6vuzNoYcO3UflbLGzszMmY6', '2024-09-17 12:00:00 +03:00', '2024-09-17 12:50:00 +03:00', 'some packed data');

insert into public.clients (token, user_id, created_at, expired_at)
values
 (E'\\x8CA0C5A18320FC2F264CFA95639EA27888727C6090D6F9CB0D6C5798A93FCB63', 'be341e38-b8a9-4230-af77-fcc34c9f2e13', '2024-09-17 18:38:50.799916 +00:00', '2025-09-17 18:38:50.680984 +00:00'),
 (E'\\x862AB376DF9DBD090F28F9DD9A2F5F1C9F88F05D27B63AE3942B5057C6BA2688', '62911ad4-883f-4b5a-929d-c4a5766560f8', '2024-09-17 22:33:42.264908 +03:00', null),
 (E'\\xC4B7F91016F52C039804D05E61C67A87A51BB8CD78FF04E51AB769ED8336D77E', 'd581d082-4b74-4dcf-8db3-cbb6e9a2f996', '2024-09-17 22:33:42.264908 +03:00', null),
 (E'\\x7210ABC35DC938383CE233297698D1B3B5CEA3AE1F0A75E69CBF48961B841EDB', 'afa37fce-557f-4b65-afdf-11f54cebe07a', '2024-09-22 21:09:11.430422 +03:00', null);


insert into public.storage (key, user_id, description, filename, blob, created_at, updated_at)
values
 ('some-exist-key', 'be341e38-b8a9-4230-af77-fcc34c9f2e13', null, null, 'some existed blob data', '2024-09-17 12:00:00 +03:00', '2024-09-17 12:50:00 +03:00'),
 ('some-exist-key1', 'be341e38-b8a9-4230-af77-fcc34c9f2e13', 'new description', null, 'some existed more new blob data', '2024-09-17 12:00:00 +03:00', '2024-09-17 12:50:00 +03:00'),
 ('some-exist-key2', 'be341e38-b8a9-4230-af77-fcc34c9f2e13', 'new description', null, 'some existed more new blob data', '2024-09-17 12:00:00 +03:00', '2024-09-17 12:50:00 +03:00'),

 ('some-exist-key', 'd581d082-4b74-4dcf-8db3-cbb6e9a2f996', null, null, 'some existed blob data', '2024-09-17 13:00:00 +03:00', null),
 ('some-exist-key1', 'd581d082-4b74-4dcf-8db3-cbb6e9a2f996', 'new description', null, 'some existed more new blob data1', '2024-09-17 12:01:00 +03:00', '2024-09-17 12:50:00 +03:00'),
 ('some-exist-key2', 'd581d082-4b74-4dcf-8db3-cbb6e9a2f996', 'new description2', null, 'some existed more new blob data2', '2024-09-17 12:02:00 +03:00', '2024-09-17 12:50:00 +03:00'),

 ('some-exist-key-for-delete', 'afa37fce-557f-4b65-afdf-11f54cebe07a', null, null, 'some existed blob data', '2024-09-17 12:00:00 +03:00', null),
 ('some-exist-key1-for-delete', 'afa37fce-557f-4b65-afdf-11f54cebe07a', 'new description', null, 'some existed more new blob data1', '2024-09-17 12:00:00 +03:00', '2024-09-17 12:50:00 +03:00'),
 ('some-exist-key2-for-delete', 'afa37fce-557f-4b65-afdf-11f54cebe07a', 'new description2', null, 'some existed more new blob data2', '2024-09-17 12:00:00 +03:00', '2024-09-17 12:50:00 +03:00')
