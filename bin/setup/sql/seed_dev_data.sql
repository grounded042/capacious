-- load the extension that allows us to generate UUIDs
CREATE EXTENSION "uuid-ossp";

DO $$
DECLARE
    sid uuid;
    uid uuid;
BEGIN

    INSERT INTO events(event_id, name, description, location,
        start_time, end_time, created_at, updated_at)
    VALUES (
        'cd7bc650-2e71-11e5-a390-675459d99309', 'Picnic', 'Your normal picnic.', 'The Park', 
        '2015-12-15 17:00:00.000000', '2015-12-15 22:00:00.000000', '2015-07-11 22:36:31.024391', 
        '2015-07-11 22:36:31.024391'
    );

    INSERT INTO guests(guest_id, first_name, last_name, attending)
    VALUES (
    	'24669e54-5ee2-11e5-a379-7b2796b289b2', 'Saxton', 'Hale', FALSE
	);

    INSERT INTO guests(guest_id, first_name, last_name, attending)
    VALUES (
        '81e6d338-7917-11e5-8b8e-a37beb0fdab8', 'Helen', '', FALSE
    );


    INSERT INTO invitees(invitee_id, fk_event_id, fk_guest_id, email)
    VALUES (
        'fb3c11f8-7917-11e5-8b8e-b3a0b1b9b068', 'cd7bc650-2e71-11e5-a390-675459d99309',
        '24669e54-5ee2-11e5-a379-7b2796b289b2', 'shale@mann.co'
    );

    INSERT INTO invitee_friends(invitee_friend_id, fk_invitee_id, fk_guest_id)
    VALUES (
        'e6afb5b0-7b64-11e5-b861-1f0fc9657754', 'fb3c11f8-7917-11e5-8b8e-b3a0b1b9b068', 
        '81e6d338-7917-11e5-8b8e-a37beb0fdab8'
    );

END $$
