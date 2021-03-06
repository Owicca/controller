'use strict';

window.setupVideo = (id) => {
  let url = "/items/" + id + "/";

  if(videojs.getAllPlayers().length == 0) {
    videojs(document.getElementById("video"),
      {
        "sources": [{
          "src": url,
          "type": "video/mp4"
        }],
        "height": 768,
        "width": 480,
        "preload": "auto",
        "controls": true,
        "inactivityTimeout": 0,
      });
  } else {
    videojs.getAllPlayers()[0].src({
      "src": url,
      "type": "video/mp4",
    });
  }
}

const e = React.createElement;
const appContainer = document.getElementById("app");
let Controller = function() {
  let self = this;

  self.sendRequest = (url, method, data) => fetch(url, {mode: "cors", method: method,
    headers: {'Accept': 'application/json','Content-Type': 'application/json'}, body: data});

  self.getElements = () => self.sendRequest("/items/", "GET", null).then((res) => res.json());

  self.deleteElement = (id) => self.sendRequest("/items/"+id+"/", "DELETE", null).then((res) => res.json());

  return self;
}

var controller = new Controller();

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
    let deleteBut = this.props.canDelete ? e("button", {onClick: (e) => this.deleteElement(this.props.id)}, "Delete") : null;
    let name = e("span", null, this.props.name);
    let item = e("span", {id: this.props.id,
      onClick: (e) => window.setupVideo(this.props.id)
    }, name);

    return e("li", null, item, deleteBut);
  }
}

class List extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      elements: [],
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
        window.location.reload();
      } else {
        alert(js.error);
      }
    });
  }

  render() {
    let elements = this.state.elements.map((elem, idx) => e(ListElement, {key: idx, id: idx, name: elem.name,
      canDelete: true, deleteElement: this.deleteElement}));

    return e("ul", {id: "list"}, elements);
  }

}

class App extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    return e(List, {getElements: controller.getElements, deleteElement: controller.deleteElement});
  }
}

ReactDOM.render(e(App, {controller: controller}), appContainer);
