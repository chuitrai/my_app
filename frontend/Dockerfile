FROM node:18-alpine AS builder 

WORKDIR /app 

COPY package*.json ./

RUN npm install

COPY . . 

ARG REACT_APP_API_URL
ENV REACT_APP_API_URL=${REACT_APP_API_URL}

RUN npm run build 

FROM nginx:stable-alpine

COPY --from=builder /app/build /usr/share/nginx/html

EXPOSE 80
