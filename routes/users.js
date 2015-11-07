'use strict';

const express = require('express');
let router = express.Router();

/* GET users listing. */
router.get('/', (req, res, next) => {
  res.render('index', { title: 'OAuth' });
});

module.exports = router;
