FROM node:17

# Set working directory 
WORKDIR /opt/catApi
# and copy files into that directory 
COPY . ./

RUN npm install -g npm@8 && npm install --save-dev

CMD npm start
