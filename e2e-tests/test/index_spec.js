import { should } from 'chai';
import supertest from 'supertest';

import {
  isStringValidUUID as validUUID,
  isDateLessThanASecondOld as validDate
} from '../helpers';

let api = supertest('http://localhost:8000/api/v1');

describe('events', () => {
  let working_event_id = "cd7bc650-2e71-11e5-a390-675459d99309";
  let event_id_list = [];

  describe('creating', () => {
    describe('with valid data', () => {
      it('return a valid obj', (done) => {
        api.post('/events')
        .send({
          name: "Christmas Party",
          description: "A Christmas Party"
        })
        .set('Accept', 'application/json')
        .expect(201)
        .expect(function(res) {
          // check things that can't be checked by comparing objs

          // check the UUID
          let cur_event_id = res.body.event_id;
          if (!validUUID(cur_event_id)) {
            throw new Error("event_id is not a UUID")
          }
          event_id_list.push(cur_event_id);
          res.body.event_id = 'FIXED_ID';

          // make sure the created at and updated at dates are less than a second
          // old
          if (!validDate(res.body.created_at) && !validDate(new Date(res.body.updated_at))) {
            throw new Error("bad created_at or updated_at");
          }

          res.body.created_at = "FIXED_DATE";
          res.body.updated_at = "FIXED_DATE";
        })
        .expect({
          event_id: 'FIXED_ID',
          name: "Christmas Party",
          description: "A Christmas Party",
          location: "",
          start_time: "0001-01-01T00:00:00Z",
          end_time: "0001-01-01T00:00:00Z",
          respond_by: "0001-01-01T00:00:00Z",
          allowed_friends: 0,
          updated_at: "FIXED_DATE",
          created_at: "FIXED_DATE"
        })
        .expect('Content-Type', 'application/json', done);
      });
    });

    describe('with invalid data', () => {
      describe('name that already exists', () => {
        it('return a 500', (done) => {
          api.post('/events')
          .send({
            name: "Christmas Party",
            description: "A Christmas Party"
          })
          .set('Accept', 'application/json')
          .expect(500, done);
        });
      });
      // TODO: make sure it does not create the event
    });
  });

  describe('getting', () => {
    describe('with a valid id', () => {
      it('should return a specific object', (done) => {
        api.get('/events/' + working_event_id)
        .set('Accept', 'application/json')
        .expect(200)
        .expect('Content-Type', 'application/json')
        .expect({
          event_id: working_event_id,
          name: "Picnic",
          description: "Your normal picnic.",
          location: "The Park",
          start_time: "2015-12-15T17:00:00Z",
          end_time: "2015-12-15T22:00:00Z",
          respond_by: "2015-12-05T22:00:00Z",
          allowed_friends: 2,
          created_at: "2015-07-11T22:36:31.024391Z",
          updated_at: "2015-07-11T22:36:31.024391Z"
        }, done);
      });
    });

    describe('with an invalid id', () => {
      it('should return a 404', (done) => {
        api.get('/events/cd7bc650-2e71-11e5-a390-675459d99308')
        .set('Accept', 'application/json')
        .expect(404, done);
      });
    });

  });

  // TODO:
  // after('delete any events that were created during testing', () => {event_id_list.forEach((event_id) => {
  //   console.log(event_id);
  //   });
  // });

});
