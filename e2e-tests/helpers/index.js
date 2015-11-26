/** @module capacious-e2e-helper */

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
