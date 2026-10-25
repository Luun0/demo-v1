# syntax=docker/dockerfile:1.7
FROM node:22-alpine AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci --ignore-scripts
COPY . .
RUN npm run build && npm run size:bundle

FROM nginxinc/nginx-unprivileged:1.27-alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY --from=deploy nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=3s --retries=5 CMD wget -qO- http://127.0.0.1:8080/ >/dev/null || exit 1
