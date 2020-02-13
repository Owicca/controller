'use strict';

import React from 'react';
import ReactDom from 'react-dom';
import List from './List';


class App extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    let controller = this.props.controller;
    return (<List getElements={controller.getElements} deleteElement={controller.deleteElement} />);
  }
}

export default App;
