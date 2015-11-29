/** @module capacious-e2e-helper */

import jwt from 'jsonwebtoken';

/** The name of the module. */
export const name = 'capacious-e2e-helper';

/**
 * check and see if a string is a valid UUID.
 * @param  {string} checkMe - the string to check fo UUID validity.
 * @return {Boolean} true if checkMe is a valid UUID, false if checkMe is not a
 * valid UUID .
 */
export function isStringValidUUID(checkMe) {
  return /[0-9a-f]{8}-([0-9a-f]{4}-){3}[0-9a-f]{12}/.test(checkMe);
}

/**
 * check and see if a date is less than a second old
 * @param  {string|Date} date - the date to check
 * @return {Boolean} true if date is less than a second old, false if date is
 * not less than a second old
 */
export function isDateLessThanASecondOld(date) {
  let now = new Date();
  let before_now = new Date(now.getTime() - 1000);
  let compare_me = new Date(date);

  return (compare_me > before_now && compare_me < now);
}

/**
 * get a JWT that is valid both in structure and in content
 * @param  {string} secret - the secret to use in the JWT
 * @return {string} valid JWT with legit user id
 */
export function validJWT(secret) {
  return jwt.sign({ sub: "cd7bc650-2e71-11e5-a390-675459b99309" }, secret, {
    algorithm: "HS512",
    expiresIn: "2 days",
  });
}

/**
 * get a JWT that is valid in structure, but the user id inside is not valid
 * @param  {string} secret - the secret to use in the JWT
 * @return {string} valid JWT with an invalid user id
 */
export function validJWTWithInvalidUser(secret) {
  return jwt.sign({ sub: "81e6d338-7917-11e5-8b8e-a37beb0fdae8" }, secret, {
    algorithm: "HS512",
    expiresIn: "2 days",
  });
}

/**
 * given a uuid, validate that it is indeed a UUID and then pass back a string
 * to replace it. Throw an error if it's not a valid UUID
 * @param  {string} uuid - the UUID to validate
 * @return {string} the fixed string to assign to replaced UUIDs
 */
export function validateAndCleanUUID(uuid) {
  if (!isStringValidUUID(uuid)) {
    throw new Error("not a valid uuid")
  }

  return 'FIXED_ID';
}

/**
 * given menu item choices, check that the attributes are valid UUIDs and then
 * set them to 'FIXED_ID' to aid in testing dynamic ids
 * @param  {array} choices - the array of choices to validate and clean
 * @return {array} the cleaned and validated array of choices
 */
export function validateAndCleanMenuChoicesUUIDs(choices) {
  return choices.map((choice) => {
    choice.menu_choice_id = validateAndCleanUUID(choice.menu_choice_id);
    choice.menu_item_id = validateAndCleanUUID(choice.menu_item_id);
    choice.menu_item_option_id = validateAndCleanUUID(choice.menu_item_option_id);

    return choice;
  });
}
