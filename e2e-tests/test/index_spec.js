import { should } from 'chai';
import supertest from 'supertest';

import {
  isStringValidUUID as validUUID,
  isDateLessThanASecondOld as validDate
} from '../helpers';

let api = supertest(`http://localhost:${process.env.PORT}/api/v1`);

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

    describe('menu items', () => {
      describe('with a valid event id', () => {
        it('should return a specific object', (done) => {
          api.get('/events/' + working_event_id + '/relationships/menu_items')
          .set('Accept', 'application/json')
          .expect(200)
          .expect('Content-Type', 'application/json')
          .expect(function(res) {
            // for now, until we figure out dates, block them out

            res.body[0].options[0].created_at = "FIXED_DATE";
            res.body[0].options[0].updated_at = "FIXED_DATE";
            res.body[0].options[1].created_at = "FIXED_DATE";
            res.body[0].options[1].updated_at = "FIXED_DATE";
            res.body[0].options[2].created_at = "FIXED_DATE";
            res.body[0].options[2].updated_at = "FIXED_DATE";
            res.body[0].created_at = "FIXED_DATE";
            res.body[0].updated_at = "FIXED_DATE";

            res.body[1].options[0].created_at = "FIXED_DATE";
            res.body[1].options[0].updated_at = "FIXED_DATE";
            res.body[1].options[1].created_at = "FIXED_DATE";
            res.body[1].options[1].updated_at = "FIXED_DATE";
            res.body[1].created_at = "FIXED_DATE";
            res.body[1].updated_at = "FIXED_DATE";

            res.body[2].options[0].created_at = "FIXED_DATE";
            res.body[2].options[0].updated_at = "FIXED_DATE";
            res.body[2].options[1].created_at = "FIXED_DATE";
            res.body[2].options[1].updated_at = "FIXED_DATE";
            res.body[2].created_at = "FIXED_DATE";
            res.body[2].updated_at = "FIXED_DATE";
          })
          .expect(
            [
              {
                menu_item_id: "f167eb18-864e-11e5-a016-6b70107c9bc3",
                item_order: 1,
                name: "Snacks",
                num_choices: 1,
                options: [
                  {
                    menu_item_option_id: "3ab2d4f0-8658-11e5-9e1b-87e2a7e99275",
                    name: "Cheese & Crackers",
                    description: "Your typical cheese and crackers snack.",
                    created_at: "FIXED_DATE",
                    updated_at: "FIXED_DATE"
                  },
                  {
                    menu_item_option_id: "3ab2e3e6-8658-11e5-9e1b-87685ca7bddd",
                    name: "Pretzels",
                    description: "See name.",
                    created_at: "FIXED_DATE",
                    updated_at: "FIXED_DATE"
                  },
                  {
                    menu_item_option_id: "3ab2e7b0-8658-11e5-9e1b-0b8bf81bc16c",
                    name: "Graham Crackers",
                    description: "A cracker made of graham.",
                    created_at: "FIXED_DATE",
                    updated_at: "FIXED_DATE"
                  }
                ],
                created_at: "FIXED_DATE",
                updated_at: "FIXED_DATE"
              },
              {
                menu_item_id: "f1680616-864e-11e5-a016-63f8fbffdc49",
                item_order: 2,
                name: "Sandwich",
                num_choices: 1,
                options: [
                  {
                    menu_item_option_id: "3ab2eb0c-8658-11e5-9e1b-a75c88531ca7",
                    name: "BLT",
                    description: "Bacon, lettuce, and tomato. A classic.",
                    created_at: "FIXED_DATE",
                    updated_at: "FIXED_DATE"
                  },
                  {
                    menu_item_option_id: "3ab2ee68-8658-11e5-9e1b-4f74a992f1df",
                    name: "Grilled Cheese",
                    description: "You cannnot go wrong.",
                    created_at: "FIXED_DATE",
                    updated_at: "FIXED_DATE"
                  }
                ],
                created_at: "FIXED_DATE",
                updated_at: "FIXED_DATE"
              },
              {
                menu_item_id: "f1680ac6-864e-11e5-a016-cb0185cdad5a",
                item_order: 3,
                name: "Dessert",
                num_choices: 1,
                options: [
                  {
                    menu_item_option_id: "3ab2f624-8658-11e5-9e1b-4be6473d4b3c",
                    name: "Brownies",
                    description: "Moist and delicious.",
                    created_at: "FIXED_DATE",
                    updated_at: "FIXED_DATE"
                  },
                  {
                    menu_item_option_id: "3ab2fdb8-8658-11e5-9e1b-cf4a9afb8def",
                    name: "Chocolate Chip Cookies",
                    description: "Gooey and good.",
                    created_at: "FIXED_DATE",
                    updated_at: "FIXED_DATE"
                  }
                ],
                created_at: "FIXED_DATE",
                updated_at: "FIXED_DATE"
              }
            ], done);
        });
      });

      describe('with an invalid event id', () => {
        it('should return a 404', (done) => {
          api.get('/events/cd7bc650-2e71-11e5-a390-675459d99308/relationships/menu_items')
          .set('Accept', 'application/json')
          .expect(404, done);
        });
      });
    });

    describe('seating request choices', () => {
      describe('with a valid event id', () => {
        it('should return a specific object', (done) => {
          api.get('/events/' + working_event_id + '/relationships/seating_request_choices')
          .set('Accept', 'application/json')
          .expect(200)
          .expect('Content-Type', 'application/json')
          .expect(
            [
              {
                "invitee_request_id": "EZXuzAu5FO9mw8UiBOqHakzvgJ1RMkOPoz4X27DpyvFwBM1i",
                "first_name": "Saxton",
                "last_name": "Hale"
              }
            ], done);
        });
      });

      describe('with an invalid event id', () => {
        it('should return a 404', (done) => {
          api.get('/events/cd7bc650-2e71-11e5-a390-675459d99308/relationships/seating_request_choices')
          .set('Accept', 'application/json')
          .expect(404, done);
        });
      });
    });
  });

  // TODO:
  // after('delete any events that were created during testing', () => {event_id_list.forEach((event_id) => {
  //   console.log(event_id);
  //   });
  // });

});
