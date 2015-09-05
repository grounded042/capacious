-- load the extension that allows us to generate UUIDs
CREATE EXTENSION "uuid-ossp";

DO $$
DECLARE
    sid uuid;
    uid uuid;
BEGIN

    SELECT uuid_generate_v1mc() INTO uid;

    INSERT INTO events(event_id, name, description, start_time, 
        end_time, created_at, updated_at)
    VALUES (
        'cd7bc650-2e71-11e5-a390-675459d99309', 'Picnic', 'Your normal picnic.', '2015-12-15 17:00:00.000000',
        '2015-12-15 22:00:00.000000', '2015-07-11 22:36:31.024391', '2015-07-11 22:36:31.024391'
    );
    
END $$
