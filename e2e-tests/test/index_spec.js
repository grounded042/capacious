import { should } from 'chai';
import supertest from 'supertest';

import {
  isStringValidUUID as validUUID,
  isDateLessThanASecondOld as validDate,
  validJWT,
  validJWTWithInvalidUser,
  validateAndCleanMenuChoicesUUIDs,
  validateAndCleanUUID,
  validateAndCleanSeatingRequestUUIDs
} from '../helpers';

let api = supertest(`http://localhost:${process.env.PORT}/api/v1`);
let secret = String(process.env.GO_JWT_MIDDLEWARE_KEY);

describe('events', () => {
  let working_event_id = "cd7bc650-2e71-11e5-a390-675459d99309";
  let event_id_list = [];

  describe('creating', () => {
    describe('with valid data', () => {
      describe('with a valid JWT', () => {
        it('should return a valid obj and a 201', (done) => {
          api.post('/events')
          .send({
            name: "Christmas Party",
            description: "A Christmas Party"
          })
          .set('Accept', 'application/json')
          .set('Authorization', `Bearer ${validJWT(secret)}`)
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
          })
          .expect('Content-Type', 'application/json', done);
        });
      });

      describe('without a JWT', () => {
        it('should return a specific message and a 401', (done) => {
          api.post('/events')
          .send({
            name: "Christmas Party 2.0",
            description: "A Christmas Party"
          })
          .set('Accept', 'application/json')
          .expect('"You need a valid user id to create an event!"\n')
          .expect(401, done);
        });
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
          .set('Authorization', `Bearer ${validJWT(secret)}`)
          .expect(500, done);
        });
      });
      // TODO: make sure it does not create the event
    });
  });

  describe('getting one', () => {
    describe('with a valid event id', () => {
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
        }, done);
      });
    });

    describe('with an invalid event id', () => {
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
                    description: "Your typical cheese and crackers snack."
                  },
                  {
                    menu_item_option_id: "3ab2e3e6-8658-11e5-9e1b-87685ca7bddd",
                    name: "Pretzels",
                    description: "See name."
                  },
                  {
                    menu_item_option_id: "3ab2e7b0-8658-11e5-9e1b-0b8bf81bc16c",
                    name: "Graham Crackers",
                    description: "A cracker made of graham."
                  }
                ],
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
                    description: "Bacon, lettuce, and tomato. A classic."
                  },
                  {
                    menu_item_option_id: "3ab2ee68-8658-11e5-9e1b-4f74a992f1df",
                    name: "Grilled Cheese",
                    description: "You cannnot go wrong."
                  }
                ],
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
                    description: "Moist and delicious."
                  },
                  {
                    menu_item_option_id: "3ab2fdb8-8658-11e5-9e1b-cf4a9afb8def",
                    name: "Chocolate Chip Cookies",
                    description: "Gooey and good."
                  }
                ]
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
              },
              {
                "invitee_request_id": "EZXuzAu5FO9mw8UiBOqHakzvgJ1RMkOPoz4X27DpyvFwBMxi",
                "first_name": "Soldier",
                "last_name": ""
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

  describe('getting event invitees', () => {
    describe('with a valid event id', () => {
      describe('with a valid JWT', () => {
        it('should return a specific object', (done) => {
          api.get(`/events/${working_event_id}/relationships/invitees`)
          .set('Accept', 'application/json')
          .set('Authorization', `Bearer ${validJWT(secret)}`)
          .expect('Content-Type', 'application/json')
          .expect((res) => {
            res.body = res.body.map((item) => {
              // take care of menu choices
              item.self.menu_choices = validateAndCleanMenuChoicesUUIDs(item.self.menu_choices);

              item.friends = item.friends.map((friend) => {
                friend.invitee_friend_id = validateAndCleanUUID(friend.invitee_friend_id);
                friend.self.guest_id = validateAndCleanUUID(friend.self.guest_id);
                friend.self.menu_choices = validateAndCleanMenuChoicesUUIDs(friend.self.menu_choices);

                return friend;
              });

              item.seating_request = validateAndCleanSeatingRequestUUIDs(item.seating_request, true);

              return item;
            });
          })
          .expect([
            {
              "invitee_id": "fb3c11f8-7917-11e5-8b8e-b3a0b1b9b068",
              "email": "shale@mann.co",
              "self": {
                "guest_id": "24669e54-5ee2-11e5-a379-7b2796b289b2",
                "first_name": "Saxton",
                "last_name": "Hale",
                "attending": false,
                "menu_choices": [
                  {
                    "menu_choice_id": "FIXED_ID",
                    "menu_item_id": "FIXED_ID",
                    "menu_item_option_id": "FIXED_ID"
                  },
                  {
                    "menu_choice_id": "FIXED_ID",
                    "menu_item_id": "FIXED_ID",
                    "menu_item_option_id": "FIXED_ID"
                  },
                  {
                    "menu_choice_id": "FIXED_ID",
                    "menu_item_id": "FIXED_ID",
                    "menu_item_option_id": "FIXED_ID"
                  }
                ],
                "menu_note": "Could I have some wine with the cheese and crackers?"
              },
              "friends": [
                {
                  "invitee_friend_id": "FIXED_ID",
                  "self": {
                    "guest_id": "FIXED_ID",
                    "first_name": "Helen",
                    "last_name": "",
                    "attending": false,
                    "menu_choices": [
                      {
                        "menu_choice_id": "FIXED_ID",
                        "menu_item_id": "FIXED_ID",
                        "menu_item_option_id": "FIXED_ID"
                      },
                      {
                        "menu_choice_id": "FIXED_ID",
                        "menu_item_id": "FIXED_ID",
                        "menu_item_option_id": "FIXED_ID"
                      },
                      {
                        "menu_choice_id": "FIXED_ID",
                        "menu_item_id": "FIXED_ID",
                        "menu_item_option_id": "FIXED_ID"
                      }
                    ],
                    "menu_note": ""
                  }
                }
              ],
              "seating_request": [
                {
                  "first_name": "Soldier",
                  "last_name": "",
                  "invitee_request_id": "FIXED_ID",
                  "invitee_seating_request_id": "FIXED_ID",
                }
              ]
            },
            {
              "invitee_id": "fb3c11f8-7917-11e5-8b8e-b3a0b1b9b078",
              "email": "soldier@mann.co",
              "self": {
                "guest_id": "81e6d338-7917-11e5-8b8e-a37beb0fdae8",
                "first_name": "Soldier",
                "last_name": "",
                "attending": false,
                "menu_choices": [],
                "menu_note": ""
              },
              "friends": [],
              "seating_request": []
            }
          ])
          .expect(200, done);
        });

        describe('but user is not an admin of this event', () => {
          it('should return a specific error and a 403', (done) => {
            api.get(`/events/${working_event_id}/relationships/invitees`)
            .set('Accept', 'application/json')
            .set('Authorization', `Bearer ${validJWTWithInvalidUser(secret)}`)
            .expect('"You are not authorized to view the list of invitees for this event!"\n')
            .expect(403, done);
          });
        });
      });

      describe('without a JWT', () => {
        it('should return a specific error and a 401', (done) => {
          api.get(`/events/${working_event_id}/relationships/invitees`)
          .set('Accept', 'application/json')
          .expect('Content-Type', 'application/json')
          .expect('"You need a valid user id to get a list of invitees for an event!"\n')
          .expect(401, done);
        });
      });
    });
  });

  describe('getting all', () => {
    describe('with a valid JWT', () => {
      it('should return 200 and a list of events assigned to the user in the JWT', (done) => {
        api.get('/events')
        .set('Accept', 'application/json')
        .set('Authorization', `Bearer ${validJWT(secret)}`)
        .expect(function(res) {
          // check things that can't be checked by comparing objs

          // check the UUID
          let second_event_id = res.body[1].event_id;
          if (!validUUID(second_event_id)) {
            throw new Error("second_event_id is not a UUID")
          }
          res.body[1].event_id = 'FIXED_ID';
        })
        .expect([
          {
            event_id: "cd7bc650-2e71-11e5-a390-675459d99309",
            name: "Picnic",
            description: "Your normal picnic.",
            location: "The Park",
            start_time: "2015-12-15T17:00:00Z",
            end_time: "2015-12-15T22:00:00Z",
            respond_by: "2015-12-05T22:00:00Z",
            allowed_friends: 2
          },
          {
            event_id: "FIXED_ID",
            name: "Christmas Party",
            description: "A Christmas Party",
            location: "",
            start_time: "0001-01-01T00:00:00Z",
            end_time: "0001-01-01T00:00:00Z",
            respond_by: "0001-01-01T00:00:00Z",
            allowed_friends: 0
          }
        ])
        .expect(200, done);
      });

      describe('with a user id that does not have any events', () => {
        it('should return an empty array', (done) => {
          api.get('/events')
          .set('Accept', 'application/json')
          .set('Authorization', `Bearer ${validJWTWithInvalidUser(secret)}`)
          .expect([])
          .expect(200, done);
        });
      });
    });

    describe('without a JWT', () => {
      it('should return an error and a 401', (done) => {
        api.get('/events')
        .set('Accept', 'application/json')
        .expect('"You need a valid user id to get your list of events!"\n')
        .expect(401, done);
      });
    });
  });

  // TODO:
  // after('delete any events that were created during testing', () => {event_id_list.forEach((event_id) => {
  //   console.log(event_id);
  //   });
  // });

});

describe('invitees', () => {
  describe('getting', () => {
    describe('with a valid invitee id', () => {
      it('should return a valid object', (done) => {
        api.get('/invitees/fb3c11f8-7917-11e5-8b8e-b3a0b1b9b068')
        .set('Accept', 'application/json')
        .expect(200)
        .expect('Content-Type', 'application/json')
        .expect((res) => {
          res.body.seating_request = validateAndCleanSeatingRequestUUIDs(res.body.seating_request, false);
        })
        .expect(
          {
            invitee_id: "fb3c11f8-7917-11e5-8b8e-b3a0b1b9b068",
            email: "shale@mann.co",
            self: {
              guest_id: "24669e54-5ee2-11e5-a379-7b2796b289b2",
              first_name: "Saxton",
              last_name: "Hale",
              attending: false,
              menu_choices: [
                {
                  menu_choice_id: "e8b849dc-9548-11e5-bea3-fbb30297c5f4",
                  menu_item_id: "f167eb18-864e-11e5-a016-6b70107c9bc3",
                  menu_item_option_id: "3ab2d4f0-8658-11e5-9e1b-87e2a7e99275"
                },
                {
                  menu_choice_id: "e8b85cce-9548-11e5-bea3-6b3c1ff816bb",
                  menu_item_id: "f1680616-864e-11e5-a016-63f8fbffdc49",
                  menu_item_option_id: "3ab2eb0c-8658-11e5-9e1b-a75c88531ca7"
                },
                {
                  menu_choice_id: "e8b864e4-9548-11e5-bea3-e73560bb934e",
                  menu_item_id: "f1680ac6-864e-11e5-a016-cb0185cdad5a",
                  menu_item_option_id: "3ab2f624-8658-11e5-9e1b-4be6473d4b3c"
                }
              ],
              menu_note: "Could I have some wine with the cheese and crackers?"
            },
            friends: [
              {
                invitee_friend_id: "e6afb5b0-7b64-11e5-b861-1f0fc9657754",
                self: {
                  guest_id: "81e6d338-7917-11e5-8b8e-a37beb0fdab8",
                  first_name: "Helen",
                  last_name: "",
                  attending: false,
                  menu_choices: [
                    {
                      menu_choice_id: "e8b86a48-9548-11e5-bea3-83652079016b",
                      menu_item_id: "f167eb18-864e-11e5-a016-6b70107c9bc3",
                      menu_item_option_id: "3ab2e3e6-8658-11e5-9e1b-87685ca7bddd"
                    },
                    {
                      menu_choice_id: "e8b86f5c-9548-11e5-bea3-6f7c95e85662",
                      menu_item_id: "f1680616-864e-11e5-a016-63f8fbffdc49",
                      menu_item_option_id: "3ab2ee68-8658-11e5-9e1b-4f74a992f1df"
                    },
                    {
                      menu_choice_id: "e8b87448-9548-11e5-bea3-834551d829f5",
                      menu_item_id: "f1680ac6-864e-11e5-a016-cb0185cdad5a",
                      menu_item_option_id: "3ab2fdb8-8658-11e5-9e1b-cf4a9afb8def"
                    }
                  ],
                  menu_note: ""
                }
              }
            ],
            seating_request: [
              {
                "first_name": "Soldier",
                "last_name": "",
                "invitee_request_id": "FIXED_ID",
                "invitee_seating_request_id": "FIXED_ID",
              }
            ]
          },
        done);
      });
    });
  });

  describe('editting', () => {
    describe('with a valid invitee id', () => {
      it('should return a valid, updated object', (done) => {
        api.patch('/invitees/fb3c11f8-7917-11e5-8b8e-b3a0b1b9b068')
        .set('Accept', 'application/json')
        .send({
          email: "shale@mann.co",
          self: {
            guest_id: "24669e54-5ee2-11e5-a379-7b2796b289b2",
            first_name: "Saxton",
            last_name: "Hale",
            attending: true
          }
        })
        .expect(200)
        .expect('Content-Type', 'application/json')
        .expect(
          {
            invitee_id: "fb3c11f8-7917-11e5-8b8e-b3a0b1b9b068",
            email: "shale@mann.co",
            self: {
              guest_id: "24669e54-5ee2-11e5-a379-7b2796b289b2",
              first_name: "Saxton",
              last_name: "Hale",
              attending: true,
              menu_choices: null,
              menu_note: ""
            },
            friends: null,
            seating_request: null
          },
        done);
      });
    });
  });

  describe('creating invitee friend', () => {
    it('should return a valid, new object', (done) => {
      api.post('/invitees/fb3c11f8-7917-11e5-8b8e-b3a0b1b9b068/relationships/friends')
      .set('Accept', 'application/json')
      .send({
        self: {
          first_name: "Friend",
          last_name: "",
          attending: true
        }
      })
      .expect(201)
      .expect('Content-Type', 'application/json')
      .expect((res) => {
        // check the UUID
        let cur_id = res.body.invitee_friend_id;
        if (!validUUID(cur_id)) {
          throw new Error("invitee_friend_id is not a UUID")
        }
        res.body.invitee_friend_id = 'FIXED_ID';

        cur_id = res.body.self.guest_id;
        if (!validUUID(cur_id)) {
          throw new Error("self.guest_id is not a UUID")
        }
        res.body.self.guest_id = 'FIXED_ID';
      })
      .expect(
        {
          invitee_friend_id: "FIXED_ID",
          self: {
            guest_id: "FIXED_ID",
            first_name: "Friend",
            last_name: "",
            attending: true,
            menu_choices: null,
            menu_note: ""
          }
        },
      done);
    });
  });

  describe('editting invitee friend', () => {
    it('should return a valid, updated object', (done) => {
      api.patch('/invitees/fb3c11f8-7917-11e5-8b8e-b3a0b1b9b068/relationships/friends/e6afb5b0-7b64-11e5-b861-1f0fc9657754')
      .set('Accept', 'application/json')
      .send({
        self: {
          guest_id: "81e6d338-7917-11e5-8b8e-a37beb0fdab8",
          first_name: "Helen 2",
          last_name: "",
          attending: false
        }
      })
      .expect(200)
      .expect('Content-Type', 'application/json')
      .expect(
        {
          invitee_friend_id: "e6afb5b0-7b64-11e5-b861-1f0fc9657754",
          self: {
            guest_id: "81e6d338-7917-11e5-8b8e-a37beb0fdab8",
            first_name: "Helen 2",
            last_name: "",
            attending: false,
            menu_choices: null,
            menu_note: ""
          }
        },
      done);
    });
  });

  describe('setting menu choices', () => {
    it('should return an object with valid UUIDs', (done) => {
      api.post('/invitees/fb3c11f8-7917-11e5-8b8e-b3a0b1b9b068/relationships/menu_choices')
      .set('Accept', 'application/json')
      .send([
          {
            menu_item_option_id: "3ab2fdb8-8658-11e5-9e1b-cf4a9afb8def",
            menu_item_id: "f1680ac6-864e-11e5-a016-cb0185cdad5a"
          }
      ])
      .expect(200)
      .expect('Content-Type', 'application/json')
      .expect((res) => {
        // check the UUID
        let cur_id = res.body[0].menu_choice_id;
        if (!validUUID(cur_id)) {
          throw new Error("menu_choice_id is not a UUID")
        }
        res.body[0].menu_choice_id = 'FIXED_ID';
      })
      .expect([
          {
            menu_choice_id: "FIXED_ID",
            menu_item_id: "f1680ac6-864e-11e5-a016-cb0185cdad5a",
            menu_item_option_id: "3ab2fdb8-8658-11e5-9e1b-cf4a9afb8def"
          }
        ],
      done);
    });
  });

  describe('setting menu notes', () => {
    it('should return an object with valid UUIDs', (done) => {
      api.post('/invitees/fb3c11f8-7917-11e5-8b8e-b3a0b1b9b068/relationships/menu_note')
      .set('Accept', 'application/json')
      .send({
        note_body: "I like cheese."
      })
      .expect(200)
      .expect('Content-Type', 'application/json')
      .expect((res) => {
        // check the UUID
        let cur_id = res.body.menu_note_id;
        if (!validUUID(cur_id)) {
          throw new Error("menu_note_id is not a UUID")
        }
        res.body.menu_note_id = 'FIXED_ID';
      })
      .expect({
        menu_note_id: "FIXED_ID",
        note_body: "I like cheese."
      },
      done);
    });
  });

  describe('setting seating requests', () => {
    it('should return an object with valid UUIDs', (done) => {
      api.post('/invitees/fb3c11f8-7917-11e5-8b8e-b3a0b1b9b068/relationships/seating_requests')
      .set('Accept', 'application/json')
      .send([
        {
          invitee_request_id: "EZXuzAu5FO9mw8UiBOqHakzvgJ1RMkOPoz4X27DpyvFwBMxi",
        }
      ])
      .expect(200)
      .expect('Content-Type', 'application/json')
      .expect((res) => {
        // check the UUID
        let cur_id = res.body[0].invitee_seating_request_id;
        if (!validUUID(cur_id)) {
          throw new Error("invitee_seating_request_id is not a UUID")
        }
        res.body[0].invitee_seating_request_id = 'FIXED_ID';
      })
      .expect([
        {
          invitee_seating_request_id: "FIXED_ID",
          invitee_request_id: "EZXuzAu5FO9mw8UiBOqHakzvgJ1RMkOPoz4X27DpyvFwBMxi",
          first_name: "",
          last_name: ""
        }
      ],
      done);
    });
  });

  describe('setting invitee friend menu choices', () => {
    it('should return an object with valid UUIDs', (done) => {
      api.post('/invitees/fb3c11f8-7917-11e5-8b8e-b3a0b1b9b068/relationships/friends/e6afb5b0-7b64-11e5-b861-1f0fc9657754/relationships/menu_choices')
      .set('Accept', 'application/json')
      .send([
          {
            menu_item_option_id: "3ab2fdb8-8658-11e5-9e1b-cf4a9afb8def",
            menu_item_id: "f1680ac6-864e-11e5-a016-cb0185cdad5a"
          }
      ])
      .expect(200)
      .expect('Content-Type', 'application/json')
      .expect((res) => {
        // check the UUID
        let cur_id = res.body[0].menu_choice_id;
        if (!validUUID(cur_id)) {
          throw new Error("menu_choice_id is not a UUID")
        }
        res.body[0].menu_choice_id = 'FIXED_ID';
      })
      .expect([
          {
            menu_choice_id: "FIXED_ID",
            menu_item_id: "f1680ac6-864e-11e5-a016-cb0185cdad5a",
            menu_item_option_id: "3ab2fdb8-8658-11e5-9e1b-cf4a9afb8def"
          }
        ],
      done);
    });
  });

  describe('setting invitee friend menu notes', () => {
    it('should return an object with valid UUIDs', (done) => {
      api.post('/invitees/fb3c11f8-7917-11e5-8b8e-b3a0b1b9b068/relationships/friends/e6afb5b0-7b64-11e5-b861-1f0fc9657754/relationships/menu_note')
      .set('Accept', 'application/json')
      .send({
        note_body: "Gluten free please."
      })
      .expect(200)
      .expect('Content-Type', 'application/json')
      .expect((res) => {
        // check the UUID
        let cur_id = res.body.menu_note_id;
        if (!validUUID(cur_id)) {
          throw new Error("menu_note_id is not a UUID")
        }
        res.body.menu_note_id = 'FIXED_ID';
      })
      .expect({
        menu_note_id: "FIXED_ID",
        note_body: "Gluten free please."
      },
      done);
    });
  });
});
