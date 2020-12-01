ARG ARCH="amd64"
ARG OS="linux"

FROM node:10-slim as app
COPY frontend/ /app
WORKDIR /app
RUN npm install && npm run build

FROM golang:alpine as svc
WORKDIR /build
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend .
RUN go build -o plansharesvc .

FROM nginx:alpine
RUN apk add supervisor
COPY ./supervisord.conf /etc/
COPY nginx.conf /etc/nginx/nginx.conf
RUN mkdir /app /svc
COPY --from=app /app/dist /app
COPY --from=svc /build/plansharesvc /svc
ENTRYPOINT /usr/bin/supervisord -c /etc/supervisord.conf
