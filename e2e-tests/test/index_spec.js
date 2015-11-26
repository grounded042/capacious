import { should } from 'chai';
import supertest from 'supertest';

import {
  isStringValidUUID as validUUID,
  isDateLessThanASecondOld as validDate
} from '../helpers';

let api = supertest('http://localhost:8000/api/v1');

describe('events', () => {
  let working_event_id = "";
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

  // TODO:
  // after('delete any events that were created during testing', () => {event_id_list.forEach((event_id) => {
  //   console.log(event_id);
  //   });
  // });

});
