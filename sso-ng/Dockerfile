### Stage: 1 ###
FROM node:14-alpine as Build
WORKDIR /app
COPY ./package.json ./package-lock.json /app/
RUN npm install
COPY . /app
#This will create dist/ and dist-server/ folders in docker
RUN npm run build:ssr

### Stage: 2 ###
FROM node:14-alpine
WORKDIR /app
COPY --from=Build /app/package.json /app
COPY --from=Build /app/dist/ /app/dist
COPY --from=Build /app/dist/shopping/server /app/dist-server
# RUN npm install

EXPOSE 4000

CMD ["npm", "run", "serve:ssr"]