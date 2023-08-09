# This image creates a standalone web server for demoing the prototype
FROM node:16@sha256:6cd6581a9ae814ebbc8077afd63ecbf4d37f9b59ac9d9304a1b5a9e743ac13e0

ARG PROJECT_ROOT="."

WORKDIR /usr/src/app
RUN echo '{"type": "module"}' > package.json
RUN npm install express@~4.18 http-proxy-middleware@~2.0
COPY  ${PROJECT_ROOT}/server.js ./
COPY ${PROJECT_ROOT}/build ./build

ENV SOURCEGRAPH_API_HOST=https://sourcegraph.com
EXPOSE 4173
CMD ["node", "server.js"]
