'use strict';

const express = require('express');
const crypto = require('crypto');
const passport = require('passport');

let router = express.Router();

router.get('/', (req, res, next) => {
  if (!req.session.state) {
    var current_date = (new Date()).valueOf().toString();
    var random = Math.random().toString();
    var hash = crypto.createHash('sha1').update(current_date + random).digest('hex');

    req.session.state = hash
  }

  passport.authenticate('facebook', {
    state: req.session.state
  })(req, res, next);
});

router.get('/callback', (req, res, next) => {
  if (!req.session.state) {
    return res.status(400).send({err: 'no state parameter'});
  }

  // CSRF verification
  if (req.query.state !== req.session.state) {
    return res.status(400).send({err: 'invalid state parameter'});
  }

  passport.authenticate('facebook', {
    failureRedirect: '/',
    successRedirect: '/me'
  })(req, res, next);
});

module.exports = router;
