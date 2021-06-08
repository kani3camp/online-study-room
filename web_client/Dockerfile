FROM node:14.4.0-alpine
ENV TZ=Asia/Tokyo \
    LANG=C.UTF-8
WORKDIR /app
RUN apk update && apk add git
RUN apk add g++ make python3
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 3000