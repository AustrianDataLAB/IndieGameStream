# Stage 1: Compile and Build angular
FROM node:20-alpine as build
# Set working directory
WORKDIR /app
# Add the source code to app
COPY indiestream .
# Install all the dependencies
RUN npm run ci:install
# Generate the build of the application
RUN npm run build


# Stage 2a: Run frontend tests
FROM node:20-alpine as test
RUN apk add --no-cache chromium
# Set environment variable for Puppeteer to use Chromium installed via Alpine's package manager
ENV PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium-browser
WORKDIR /app
COPY indiestream .
RUN npm run ci:install
# Run tests using headless chromium
# CMD ["npm", "run", "test", "--", "--configuration=ci"]


# Stage 2b: Serve app with nginx server
FROM nginx:alpine
# Copy the build output to replace the default nginx contents
COPY --from=build /app/dist/indiestream/browser /usr/share/nginx/html
# Need to expose port 80
EXPOSE 80
# Start nginx 
CMD ["nginx", "-g", "daemon off;"]
