-- load the extension that allows us to generate UUIDs
CREATE EXTENSION "uuid-ossp";

-- create the tables
CREATE TABLE events (
    event_id uuid DEFAULT uuid_generate_v1mc() PRIMARY KEY,
    name varchar(255) NOT NULL UNIQUE,
    description varchar(255) NOT NULL,
    start_time timestamp,
    end_time timestamp,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);

CREATE TABLE invitees (
	invitee_id uuid DEFAULT uuid_generate_v1mc() PRIMARY KEY,
	fk_event_id uuid REFERENCES events (event_id),
	first_name varchar(255) NOT NULL,
	last_name varchar(255) NOT NULL,
	email varchar(255) NOT NULL,
	date_attending boolean,
	date_first_name varchar(255),
	date_last_name varchar(255),
	created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);
