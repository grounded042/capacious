# Capacious API

Capacious (having a lot of space inside; roomy.) has been made so that planning for a large event can be easy. If you have a lot of guests and you need to manage their responses, meal selection, etc. this tool can help you.

## Setup

1. Install dependencies: `go get`
2. Create environment file: `cp .env.example .env`
3. Run the script bellow to apply env vars.
4. Create database: `createdb capacious-go`
5. Add seed data: `make setup`

To apply environment variables:

```
while read l; do
  export $1
  done < .env
```

## Architecture
- Controllers - take care of supplying the web calls with data by consuming services. Here we can take care of marshaling the data we get from the services into JSON
- Services - respond to requests by getting the appropriate data from gateways and running validation and logic on them - this is the layer that unit tests will run on
- DAL - interface to whatever holds our data. We should be able to call functions on this layer and not care about what is db/data store technology is holding out data
- Entities - holds data models that other layers consume


## Running

1. Run: `go run main.go`
2. Navigate to [http://localhost:8000/api/v1/events](http://localhost:8000/api/v1/events) to see the magic.

## Contributing 
Just fork and make a pull request and I'll happily review it.

---
Copyright 2015 Jon Carl. All rights reserved.
