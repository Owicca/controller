#!/bin/bash

cd frontend && \
  sudo docker run -ti -v `pwd`:/app -w /app node npm i && \
  /usr/local/bin/node node_modules/webpack/bin/webpack.js \
  --config webpack.config.js
