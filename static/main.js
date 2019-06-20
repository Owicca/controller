'use strict';

const e = React.createElement;
const appContainer = document.getElementById("app");
let Controller = function() {
  let self = this;

  self.sendRequest = (url, method, data) => fetch(url, {mode: "cors", method: method,headers: {'Accept': 'application/json','Content-Type': 'application/json'}, body: data});

  self.getElements = () => self.sendRequest("/items/", "GET", null).then((res) => res.json());

  self.deleteElement = (id) => self.sendRequest("/items/"+id+"/", "DELETE", null).then((res) => res.json());

  return self;
}

var controller = new Controller();

class List extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    let elements = this.props.elements.map(elem => e(ListElement, {key: elem.id, id: elem.id, name: elem.name, canDelete: true}));

    return e("ul", {id: "list"}, elements);
  }

}

class ListElement extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    let deleteBut = this.props.canDelete ? e("button", {onClick: (e) => console.log("Delete #"+this.props.id)}, "Delete") : null;
    let name = e("span", null, this.props.name);
    let href = "/v/"+this.props.id+"/";
    let link = e("a", {id: this.props.id, href: href}, name);

    return e("li", null, link, deleteBut);
  }
}

class App extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    let elements = this.props.controller.getElements();

    return e(List, {elements: elements});
  }
}

ReactDOM.render(e(App, {controller: controller}), appContainer);
