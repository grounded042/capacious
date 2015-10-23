-- load the extension that allows us to generate UUIDs
CREATE EXTENSION "uuid-ossp";

-- create the tables
CREATE TABLE events (
    event_id uuid DEFAULT uuid_generate_v1mc() PRIMARY KEY,
    name varchar(255) NOT NULL UNIQUE,
    description varchar(255) NOT NULL,
    start_time timestamp,
    end_time timestamp,
    allowed_guests int,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);

CREATE TABLE guests (
	guest_id uuid DEFAULT uuid_generate_v1mc() PRIMARY KEY,
	first_name varchar(255) NOT NULL,
	last_name varchar(255) NOT NULL,
	attending boolean DEFAULT false,
	created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);

CREATE TABLE invitees (
	invitee_id uuid DEFAULT uuid_generate_v1mc() PRIMARY KEY,
	fk_event_id uuid REFERENCES events (event_id),
	fk_guest_id uuid REFERENCES guests (guest_id),
	email varchar(255) NOT NULL UNIQUE,
	created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);

CREATE TABLE invitee_guests (
	invitee_guest_id uuid DEFAULT uuid_generate_v1mc() PRIMARY KEY,
	fk_invitee_id uuid REFERENCES invitees (invitee_id),
	fk_guest_id uuid REFERENCES guests (guest_id),
	created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);


