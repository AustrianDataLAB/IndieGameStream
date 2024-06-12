## TL;DR
Change MYSQL password in .env.deployment
`` docker compose build; `` `` docker compose up -d``
The api will be exposed to port 8080, access it with `localhost:8080`. 


## Environment variables
The docker image will use the following environment variables:

| Key                                                | Default Value | Options                |
|----------------------------------------------------|---------------|------------------------|
| PORT                                               | "8080"        |                        |
| GIN_MODE                                           | "release"     | "release", "debug"     |
| MYSQL_HOST                                         | "mysql"       |                        |
| MYSQL_PORT                                         | "3306"        |                        |
| MYSQL_DATABASE                                     | "api"         |                        |
| MYSQL_ROOT_USER                                    | "root"        |                        |
| <span style="color:red">MYSQL_ROOT_PASSWORD</span> | <span style="color:red">"changeme"</span>    |                        |
If you use the docker image directly (without our provided docker-compose), you must specify them.