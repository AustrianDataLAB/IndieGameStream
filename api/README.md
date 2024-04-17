## Build the docker container
``docker build . -t api:latest ``
## Run the docker container
``docker run -p 8080:8080 -i api:latest``\
The api will be exposed to port 8080, access it with `localhost:8080`.