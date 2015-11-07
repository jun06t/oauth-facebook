'use strict';

let config = {};

config.oauth = {
  facebook: {
    clientID: 'your_client_id',
    clientSecret: 'your_client_secret',
    callbackURL: 'http://example.com:3000/auth/callback',
    enableProof: true,
    scope: ['email', 'user_friends', 'user_birthday', 'user_location']
  }
}

module.exports = config;
