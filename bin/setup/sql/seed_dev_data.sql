-- load the extension that allows us to generate UUIDs
CREATE EXTENSION "uuid-ossp";

DO $$
DECLARE
    sid uuid;
    uid uuid;
BEGIN

    -- create a test event
    INSERT INTO events(event_id, name, description, location,
        start_time, end_time, respond_by, allowed_friends, created_at, updated_at)
    VALUES (
      'cd7bc650-2e71-11e5-a390-675459d99309', 'Picnic', 'Your normal picnic.', 'The Park',
      '2015-12-15 17:00:00.000000', '2015-12-15 22:00:00.000000', '2015-12-05 22:00:00.000000',
      2, '2015-07-11 22:36:31.024391', '2015-07-11 22:36:31.024391'
    );

    -- create menu options for the test event
    INSERT INTO menu_items(menu_item_id, fk_event_id, item_order, name, num_choices)
    VALUES (
      'F167EB18-864E-11E5-A016-6B70107C9BC3', 'cd7bc650-2e71-11e5-a390-675459d99309',
      1, 'Snacks', 1
    );

    INSERT INTO menu_items(menu_item_id, fk_event_id, item_order, name, num_choices)
    VALUES (
      'F1680616-864E-11E5-A016-63F8FBFFDC49', 'cd7bc650-2e71-11e5-a390-675459d99309',
      2, 'Sandwich', 1
    );

    INSERT INTO menu_items(menu_item_id, fk_event_id, item_order, name, num_choices)
    VALUES (
      'F1680AC6-864E-11E5-A016-CB0185CDAD5A', 'cd7bc650-2e71-11e5-a390-675459d99309',
      3, 'Dessert', 1
    );

    INSERT INTO menu_item_options(menu_item_option_id, fk_menu_item_id, name, description)
    VALUES (
      '3AB2D4F0-8658-11E5-9E1B-87E2A7E99275', 'F167EB18-864E-11E5-A016-6B70107C9BC3',
      'Cheese & Crackers', 'Your typical cheese and crackers snack.'
    );

    INSERT INTO menu_item_options(menu_item_option_id, fk_menu_item_id, name, description)
    VALUES (
      '3AB2E3E6-8658-11E5-9E1B-87685CA7BDDD', 'F167EB18-864E-11E5-A016-6B70107C9BC3',
      'Pretzels', 'See name.'
    );

    INSERT INTO menu_item_options(menu_item_option_id, fk_menu_item_id, name, description)
    VALUES (
      '3AB2E7B0-8658-11E5-9E1B-0B8BF81BC16C', 'F167EB18-864E-11E5-A016-6B70107C9BC3',
      'Graham Crackers', 'A cracker made of graham.'
    );

    INSERT INTO menu_item_options(menu_item_option_id, fk_menu_item_id, name, description)
    VALUES (
      '3AB2EB0C-8658-11E5-9E1B-A75C88531CA7', 'F1680616-864E-11E5-A016-63F8FBFFDC49',
      'BLT', 'Bacon, lettuce, and tomato. A classic.'
    );

    INSERT INTO menu_item_options(menu_item_option_id, fk_menu_item_id, name, description)
    VALUES (
      '3AB2EE68-8658-11E5-9E1B-4F74A992F1DF', 'F1680616-864E-11E5-A016-63F8FBFFDC49',
      'Grilled Cheese', 'You cannnot go wrong.'
    );

    INSERT INTO menu_item_options(menu_item_option_id, fk_menu_item_id, name, description)
    VALUES (
      '3AB2F624-8658-11E5-9E1B-4BE6473D4B3C', 'F1680AC6-864E-11E5-A016-CB0185CDAD5A',
      'Brownies', 'Moist and delicious.'
    );

    INSERT INTO menu_item_options(menu_item_option_id, fk_menu_item_id, name, description)
    VALUES (
      '3AB2FDB8-8658-11E5-9E1B-CF4A9AFB8DEF', 'F1680AC6-864E-11E5-A016-CB0185CDAD5A',
      'Chocolate Chip Cookies', 'Gooey and good.'
    );

    -- create invitee with friend
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

    -- set menu choices for the created invitee and friend
    INSERT INTO menu_choices(menu_choice_id, fk_guest_id, fk_menu_item_option_id, fk_menu_item_id)
    VALUES (
      'e8b849dc-9548-11e5-bea3-fbb30297c5f4', '24669e54-5ee2-11e5-a379-7b2796b289b2',
      '3AB2D4F0-8658-11E5-9E1B-87E2A7E99275', 'F167EB18-864E-11E5-A016-6B70107C9BC3'
    );

    -- INSERT INTO menu_choices(fk_guest_id, fk_menu_item_option_id, fk_menu_item_id)
    -- VALUES (
    --   '24669e54-5ee2-11e5-a379-7b2796b289b2', '3AB2D4F0-8658-11E5-9E1B-87E2A7E99275',
    --   'F167EB18-864E-11E5-A016-6B70107C9BC3'
    -- );

    INSERT INTO menu_choices(menu_choice_id, fk_guest_id, fk_menu_item_option_id, fk_menu_item_id)
    VALUES (
      'e8b85cce-9548-11e5-bea3-6b3c1ff816bb', '24669e54-5ee2-11e5-a379-7b2796b289b2',
      '3AB2EB0C-8658-11E5-9E1B-A75C88531CA7', 'F1680616-864E-11E5-A016-63F8FBFFDC49'
    );

    INSERT INTO menu_choices(menu_choice_id, fk_guest_id, fk_menu_item_option_id, fk_menu_item_id)
    VALUES (
      'e8b864e4-9548-11e5-bea3-e73560bb934e', '24669e54-5ee2-11e5-a379-7b2796b289b2',
      '3AB2F624-8658-11E5-9E1B-4BE6473D4B3C', 'F1680AC6-864E-11E5-A016-CB0185CDAD5A'
    );

    INSERT INTO menu_choices(menu_choice_id, fk_guest_id, fk_menu_item_option_id, fk_menu_item_id)
    VALUES (
      'e8b86a48-9548-11e5-bea3-83652079016b', '81e6d338-7917-11e5-8b8e-a37beb0fdab8',
      '3AB2E3E6-8658-11E5-9E1B-87685CA7BDDD', 'F167EB18-864E-11E5-A016-6B70107C9BC3'
    );

    -- INSERT INTO menu_choices(fk_guest_id, fk_menu_item_option_id, fk_menu_item_id)
    -- VALUES (
    --   '81e6d338-7917-11e5-8b8e-a37beb0fdab8', '3AB2E3E6-8658-11E5-9E1B-87685CA7BDDD',
    --   'F167EB18-864E-11E5-A016-6B70107C9BC3'
    -- );

    INSERT INTO menu_choices(menu_choice_id, fk_guest_id, fk_menu_item_option_id, fk_menu_item_id)
    VALUES (
      'e8b86f5c-9548-11e5-bea3-6f7c95e85662', '81e6d338-7917-11e5-8b8e-a37beb0fdab8',
      '3AB2EE68-8658-11E5-9E1B-4F74A992F1DF', 'F1680616-864E-11E5-A016-63F8FBFFDC49'
    );

    INSERT INTO menu_choices(menu_choice_id, fk_guest_id, fk_menu_item_option_id, fk_menu_item_id)
    VALUES (
      'e8b87448-9548-11e5-bea3-834551d829f5', '81e6d338-7917-11e5-8b8e-a37beb0fdab8',
      '3AB2FDB8-8658-11E5-9E1B-CF4A9AFB8DEF', 'F1680AC6-864E-11E5-A016-CB0185CDAD5A'
    );

    INSERT INTO menu_notes(fk_guest_id, note_body)
    VALUES (
      '24669e54-5ee2-11e5-a379-7b2796b289b2', 'Could I have some wine with the cheese and crackers?'
    );

END $$
