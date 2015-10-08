var createCookie, deleteCookie, readCookie, startsWith;

startsWith = function(string, pattern) {
  return string.slice(0, pattern.length) === pattern;
};

createCookie = function(name, value, days) {
  var date, expires;
  expires = '';
  if (days) {
    date = new Date();
    date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
    expires = "; expires=" + (date.toGMTString());
  }
  return document.cookie = name + "=" + value + expires + "; path=/";
};

readCookie = function(name) {
  var cookie_fields, field, i, len, nameEQ;
  nameEQ = name + "=";
  cookie_fields = document.cookie.split(';');
  for (i = 0, len = cookie_fields.length; i < len; i++) {
    field = cookie_fields[i];
    while (startsWith(field, ' ')) {
      field = field.slice(1, +field.length + 1 || 9e9);
    }
    if (startsWith(field, nameEQ)) {
      return field.slice(nameEQ.length, +field.length + 1 || 9e9);
    }
  }
  return null;
};

deleteCookie = function(name) {
  return createCookie(name, '', -1);
};
