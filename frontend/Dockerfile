# Stage 1: Compile and Build angular
FROM node:10-alpine as build

# Set working directory
WORKDIR /app

# Add the source code to app
COPY indiestream .

# Install all the dependencies
RUN npm install

# Generate the build of the application
RUN npm run build

# Stage 2: Serve app with nginx server
FROM nginx:alpine

# Copy the build output to replace the default nginx contents
COPY --from=build /app/dist/indiestream/browser /usr/share/nginx/html

# Start nginx 
CMD ["nginx", "-g", "daemon off;"]
