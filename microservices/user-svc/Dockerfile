FROM node

WORKDIR /user_svc

COPY . .

RUN yarn install

EXPOSE 3001

CMD [ "yarn", "start" ]