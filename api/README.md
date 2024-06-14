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
| OAUTH_CLIENT                                       |         |              |
| AZURE_CLIENT_ID                                    |         |  |
| AZURE_TENANT_ID                                    |         |  |
| AZURE_STORAGE_ACCOUNT                              |         |  |
| AZURE_CLIENT_SECRET                                |         |  |
| AZURE_CONTAINER_NAME                               |         |  |
| AZURE_AKS_CLUSTER_NAME                             |         |  |
| AZURERM_SUBSCRIPTION_ID                            |         |  |
| AZURERM_RESOURCE_GROUP_NAME                        |         |  |


If you use the docker image directly (without our provided docker-compose), you must specify them.

The api will use a kubeconfig for talking to a Kubernetes API server. 
If --kubeconfig is set, will use the kubeconfig file at that location. 
Otherwise will assume running in cluster and use the cluster provided kubeconfig.
It also applies saner defaults for QPS and burst based on the Kubernetes controller manager defaults (20 QPS, 30 burst)

Config precedence:
* --kubeconfig flag pointing at a file
* KUBECONFIG environment variable pointing at a file
* In-cluster config if running in cluster
* $HOME/.kube/ config if exists.
