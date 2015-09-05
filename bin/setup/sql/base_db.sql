-- load the extension that allows us to generate UUIDs
CREATE EXTENSION "uuid-ossp";

-- create the tables
CREATE TABLE events (
    event_id uuid DEFAULT uuid_generate_v1mc() PRIMARY KEY,
    name varchar(255) NOT NULL UNIQUE,
    description varchar(255) NOT NULL UNIQUE,
    start_time timestamp,
    end_time timestamp,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);