'use strict';

let Controller = function() {
  let self = this;

  self.sendRequest = (url, method, data) => fetch(url, {mode: "cors", method: method,headers: {'Accept': 'application/json','Content-Type': 'application/json'}, body: data});

  self.getElements = () => self.sendRequest("/items/", "GET", null).then((res) => res.json());

  self.deleteElement = (id) => self.sendRequest("/items/"+id+"/", "DELETE", null).then((res) => res.json());

  return self;
}

export default Controller;
