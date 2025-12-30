# Stage 1: Build
FROM node:20-alpine AS build
WORKDIR /app

# Define the argument that will be passed during build
ARG REACT_APP_API_URL
# Set it as an environment variable so CRA can see it
ENV REACT_APP_API_URL=$REACT_APP_API_URL

COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Stage 2: Serve with Nginx
FROM nginx:stable-alpine
COPY --from=build /app/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]