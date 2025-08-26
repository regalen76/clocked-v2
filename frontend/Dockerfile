# syntax=docker/dockerfile:1
FROM node:20-alpine AS build
WORKDIR /app

# Respect Nuxt public runtime config at build time (default to /api)
ARG NUXT_PUBLIC_API_BASE=/api
ENV NUXT_PUBLIC_API_BASE=$NUXT_PUBLIC_API_BASE

COPY frontend/package.json ./
COPY frontend/ .

RUN npm install --no-audit --no-fund \
  && npm run build

FROM nginx:alpine
WORKDIR /usr/share/nginx/html

# Copy nginx config
COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf

# Copy built static output
COPY --from=build /app/.output/public/ /usr/share/nginx/html/

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]

