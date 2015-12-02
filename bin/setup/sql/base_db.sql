-- load the extension that allows us to generate UUIDs
CREATE EXTENSION "uuid-ossp";

-- create functions
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- create the tables
CREATE TABLE IF NOT EXISTS events (
  event_id uuid DEFAULT uuid_generate_v1mc() PRIMARY KEY,
  name varchar(255) NOT NULL UNIQUE,
  description varchar(255) NOT NULL,
  location varchar(255) NOT NULL,
  start_time timestamp,
  end_time timestamp,
  respond_by timestamp,
  allowed_friends int,
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp
);

DROP TRIGGER IF EXISTS update_event_updated_at_time ON events;
CREATE TRIGGER update_event_updated_at_time BEFORE UPDATE ON events FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

CREATE TABLE IF NOT EXISTS guests (
  guest_id uuid DEFAULT uuid_generate_v1mc() PRIMARY KEY,
  first_name varchar(255) NOT NULL,
  last_name varchar(255) NOT NULL,
  attending boolean DEFAULT false,
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp
);

DROP TRIGGER IF EXISTS update_guest_updated_at_time ON guests;
CREATE TRIGGER update_guest_updated_at_time BEFORE UPDATE ON guests FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();


CREATE TABLE IF NOT EXISTS invitees (
  invitee_id uuid DEFAULT uuid_generate_v1mc() PRIMARY KEY,
  fk_event_id uuid REFERENCES events (event_id),
  fk_guest_id uuid REFERENCES guests (guest_id),
  email varchar(255) NOT NULL UNIQUE,
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp
);

DROP TRIGGER IF EXISTS update_invitee_updated_at_time ON invitees;
CREATE TRIGGER update_invitee_updated_at_time BEFORE UPDATE ON invitees FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();


CREATE TABLE IF NOT EXISTS invitee_friends (
  invitee_friend_id uuid DEFAULT uuid_generate_v1mc() PRIMARY KEY,
  fk_invitee_id uuid REFERENCES invitees (invitee_id),
  fk_guest_id uuid REFERENCES guests (guest_id),
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp
);

DROP TRIGGER IF EXISTS update_invitee_friend_updated_at_time ON invitee_friends;
CREATE TRIGGER update_invitee_friend_updated_at_time BEFORE UPDATE ON invitee_friends FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();


CREATE TABLE IF NOT EXISTS menu_items (
  menu_item_id uuid DEFAULT uuid_generate_v1mc() PRIMARY KEY,
  fk_event_id uuid REFERENCES events (event_id),
  item_order int,
  name varchar(255),
  num_choices int,
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp,
  UNIQUE (fk_event_id, item_order)
);

DROP TRIGGER IF EXISTS update_menu_item_updated_at_time ON menu_items;
CREATE TRIGGER update_menu_item_updated_at_time BEFORE UPDATE ON menu_items FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();


CREATE TABLE IF NOT EXISTS menu_item_options (
  menu_item_option_id uuid DEFAULT uuid_generate_v1mc() PRIMARY KEY,
  fk_menu_item_id uuid REFERENCES menu_items (menu_item_id),
  name varchar(255),
  description varchar(255),
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp
);

DROP TRIGGER IF EXISTS update_menu_item_option_updated_at_time ON menu_item_options;
CREATE TRIGGER update_menu_item_option_updated_at_time BEFORE UPDATE ON menu_item_options FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();


CREATE TABLE IF NOT EXISTS menu_choices (
  menu_choice_id uuid DEFAULT uuid_generate_v1mc() PRIMARY KEY,
  fk_guest_id uuid REFERENCES guests (guest_id),
  fk_menu_item_id uuid REFERENCES menu_items (menu_item_id),
  fk_menu_item_option_id uuid REFERENCES menu_item_options (menu_item_option_id),
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp
);

DROP TRIGGER IF EXISTS update_menu_choice_updated_at_time ON menu_choices;
CREATE TRIGGER update_menu_choice_updated_at_time BEFORE UPDATE ON menu_choices FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();


CREATE TABLE IF NOT EXISTS menu_notes (
  menu_note_id uuid DEFAULT uuid_generate_v1mc() PRIMARY KEY,
  fk_guest_id uuid REFERENCES guests (guest_id),
  note_body varchar(255),
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp
);

DROP TRIGGER IF EXISTS update_menu_note_updated_at_time ON menu_notes;
CREATE TRIGGER update_menu_note_updated_at_time BEFORE UPDATE ON menu_notes FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();


CREATE TABLE IF NOT EXISTS invitee_seating_requests (
  invitee_seating_request_id uuid DEFAULT uuid_generate_v1mc() PRIMARY KEY,
  fk_invitee_id uuid REFERENCES invitees (invitee_id),
  fk_invitee_request_id uuid REFERENCES invitees (invitee_id),
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp
);

DROP TRIGGER IF EXISTS update_invitee_seating_requests_updated_at_time ON invitee_seating_requests;
CREATE TRIGGER update_invitee_seating_requests_updated_at_time BEFORE UPDATE ON invitee_seating_requests FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();


CREATE TABLE IF NOT EXISTS users (
    user_id uuid DEFAULT uuid_generate_v1mc() PRIMARY KEY,
    email varchar(255) NOT NULL UNIQUE,
    first_name varchar(255) NOT NULL,
    last_name varchar(255) NOT NULL,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);

DROP TRIGGER IF EXISTS update_users_updated_at_time ON users;
CREATE TRIGGER update_users_updated_at_time BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();


CREATE TABLE IF NOT EXISTS user_logins (
    user_login_id uuid DEFAULT uuid_generate_v1mc() PRIMARY KEY,
    fk_user_id uuid UNIQUE REFERENCES users (user_id),
    salt varchar(255) NOT NULL,
    password varchar(255) NOT NULL,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);

DROP TRIGGER IF EXISTS update_user_login_updated_at_time ON user_logins;
CREATE TRIGGER update_user_login_updated_at_time BEFORE UPDATE ON user_logins FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();


CREATE TABLE IF NOT EXISTS event_admins (
  event_admin_id uuid DEFAULT uuid_generate_v1mc() PRIMARY KEY,
  fk_user_id uuid REFERENCES users (user_id),
  fk_event_id uuid REFERENCES events (event_id),
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp,
  UNIQUE (fk_user_id, fk_event_id)
);

DROP TRIGGER IF EXISTS update_event_admin_updated_at_time ON event_admins;
CREATE TRIGGER update_event_admin_updated_at_time BEFORE UPDATE ON event_admins FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
