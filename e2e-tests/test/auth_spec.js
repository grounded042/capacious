import { should } from 'chai';
import supertest from 'supertest';
import jwt from 'jsonwebtoken';

let api = supertest(`http://localhost:${process.env.PORT}/api/v1`);
let secret = String(process.env.GO_JWT_MIDDLEWARE_KEY);

describe('auth', () => {
  describe('valid JWTs', () => {
    it('should return 200 on valid JWTs', (done) => {
      let token = jwt.sign({ sub: "user_id" }, secret, {
        algorithm: "HS512",
        expiresIn: "2 days",
      });

      api.get('/events/cd7bc650-2e71-11e5-a390-675459d99309')
      .set('Authorization', `Bearer ${token}`)
      .expect(200, done);
    });
  });

  describe('invalid JWTs', () => {
    it('should return 500 on bad JWT secret', (done) => {
      let token = jwt.sign({ sub: "user_id" }, "this_is_not_the_right_secret", {
        algorithm: "HS512",
        expiresIn: "2 days",
      });

      api.get('/events/cd7bc650-2e71-11e5-a390-675459d99309')
      .set('Authorization', `Bearer ${token}`)
      .expect(500, done);
    });

    it('should return 400 on invalid header', (done) => {
      api.get('/events/cd7bc650-2e71-11e5-a390-675459d99309')
      .set('Authorization', 'Bearer')
      .expect(400, done);
    });

    it('should return 401 on an expired JWT', (done) => {
      let token = jwt.sign({ sub: "user_id", exp: 1437265807 }, secret, {
        algorithm: "HS512",
      });

      api.get('/events/cd7bc650-2e71-11e5-a390-675459d99309')
      .set('Authorization', `Bearer ${token}`)
      .expect(401, done);
    });
  });

  describe('login', () => {
    describe('with valid creds', () => {
      it('should return a valid JWT', (done) => {
        api.post('/token')
        .send({
          email: "1498@aperturescience.com",
          password: "GLaDOS",
        })
        .set('Accept', 'application/json')
        .expect((res) => {
          let token = res.body.token;

          jwt.verify(token, secret, {
            algorithms: ["HS512"]
          });
        })
        .expect(200, done);
      });
    });

    describe('with invalid creds', () => {
      it('should return no token and 401', (done) => {
        api.post('/token')
        .send({
          email: "1498@aperturescience.com",
          password: "GLaDOs",
        })
        .set('Accept', 'application/json')
        .expect((res) => {
          if (res.body.token !== undefined) {
            throw new Error("token is not undefined!")
          };
        })
        .expect(401, done);
      });
    });
  });
});
