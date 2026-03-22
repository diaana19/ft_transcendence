FROM node:20-alpine AS build-env
COPY . /app/
WORKDIR /app
RUN npm ci && \
  npm run build

FROM nginx:alpine
RUN rm -fr /usr/share/nginx/html/*
COPY --from=build-env --chown=nginx:nginx /app/build/client /usr/share/nginx/html
COPY --chown=nginx:nginx nginx.conf /etc/nginx/conf.d/default.conf
CMD ["nginx", "-g", "daemon off;"]
