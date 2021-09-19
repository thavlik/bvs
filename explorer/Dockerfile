FROM node:16.9.1-alpine AS builder
WORKDIR /app
COPY package.json .
COPY package-lock.json .
RUN npm i
COPY . .
RUN npm run build

FROM node:16.9.1-alpine
WORKDIR /app
COPY package.json .
COPY package-lock.json .
RUN npm i --production
COPY express.js .
COPY --from=builder /app/dist dist
CMD ["node", "./express.js"]
