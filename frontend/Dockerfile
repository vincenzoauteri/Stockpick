# Stage 1: Build the React application
FROM node:20-alpine as build

WORKDIR /app

COPY package.json package-lock.json ./frontend/
COPY frontend/tsconfig.json ./frontend/
COPY frontend/vite.config.ts ./frontend/
COPY frontend/src ./frontend/src
COPY frontend/public ./frontend/public

WORKDIR /app/frontend

RUN npm install
RUN npm run build

# Stage 2: Serve the application with Nginx
FROM nginx:stable-alpine

COPY --from=build /app/frontend/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
