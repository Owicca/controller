'use strict';

import React from 'react';
import ReactDom from 'react-dom';

import Controller from './Controller';
import App from './App';


var controller = new Controller();

ReactDom.render(
  <App controller={controller} />,
  document.getElementById("app")
);
