FROM node:12

RUN npm install --unsafe-perm -g git+https://github.com/pakerfeldt/han-mqtt.git
VOLUME /config
CMD han-mqtt