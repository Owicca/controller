'use strict';
import React from 'react';
import ReactDom from 'react-dom';


class ListElement extends React.Component {
  constructor(props) {
    super(props);

    this.deleteElement = this.deleteElement.bind(this);
  }

  deleteElement(id) {
    let confirmation = confirm("Are you sure you want to delete: "+this.props.name+"?");

    if (this.props.canDelete && confirmation) {
      this.props.deleteElement(this.props.id);
    }
  }

  render() {
    let name = <span>{this.props.name}</span>;
    let href = "/items/"+this.props.href + "/";
    let deleteBut = this.props.canDelete ? <button onClick={(e) => this.deleteElement(href)}>{'Delete'}</button> : null;
    let link = <a id={this.props.href} href={href} target='_blank'>{name}</a>;

    return (<li>{link}{deleteBut}</li>);
  }
}

class List extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      elements: {},
    };

    this.deleteElement = this.deleteElement.bind(this);
  }

  componentDidMount() {
    this.props.getElements().then(js => {
      if (js.success == true) {
        this.setState({
          elements: js.data
        });
      } else {
        alert(js.error);
      }
    });
  }

  deleteElement(id) {
    this.props.deleteElement(id).then((js) => {
      if(js.success == true) {
        let elem = this.state.elements;
        elem.splice(id, 1);

        this.setState({
          elements: elem
        });
      } else {
        alert(js.error);
      }
    });
  }

  render() {
    let children = [];
    if(this.state.elements.hasOwnProperty("children")) {
      let obj = this.state.elements.children;
      for (let child in obj) {
        children.push(
          <ListElement
            key={child} href={obj[child].info.pseudoname}
            name={obj[child].info.name}
            canDelete={false}
            deleteElement={this.deleteElement} />
        );
      }
    }

    return (<ul id="list">{children}</ul>);
  }

}
export default List;
