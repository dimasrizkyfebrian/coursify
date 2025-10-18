# Stage 1: The "Builder" Stage
# Compiling frontend assets.
FROM node:20-alpine AS builder

# Setting the working directory inside the container
WORKDIR /app

# Copy the package.json and package-lock.json files first to take advantage of Docker layer caching
COPY package*.json ./

# Installing all dependencies
RUN npm install

# Copying the remaining application source code
COPY . .

# Building an application for production (generating static files in the /dist folder)
RUN npm run build

# ---

# Stage 2: The "Web Server" Stage
# Setting up an Nginx server to serve completed static files
FROM nginx:stable-alpine

# Copy ONLY the 'dist' folder that contains the build results from the 'builder' stage to the default Nginx directory that serves web content.
COPY --from=builder /app/dist /usr/share/nginx/html

# Copying the Nginx configuration file
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Inform Docker that the container will listen on port 80
EXPOSE 80

# The default command to run Nginx when the container starts
CMD ["nginx", "-g", "daemon off;"]