define HEADER
 ___           _ _       ____                      ____  _                            
|_ _|_ __   __| (_) ___ / ___| __ _ _ __ ___   ___/ ___|| |_ _ __ ___  __ _ _ __ ___  
 | || '_ \ / _` | |/ _ \ |  _ / _` | '_ ` _ \ / _ \___ \| __| '__/ _ \/ _` | '_ ` _ \ 
 | || | | | (_| | |  __/ |_| | (_| | | | | | |  __/___) | |_| | |  __/ (_| | | | | | |
|___|_| |_|\__,_|_|\___|\____|\__,_|_| |_| |_|\___|____/ \__|_|  \___|\__,_|_| |_| |_|
endef

ROOT_DIRECTORY = ../..
CLUSTER_NAME = indiegamestream
IMAGE_NAMESPACE = indiegamestream

FRONTEND_NAME=frontend
FRONTEND_IMAGE=$(IMAGE_NAMESPACE)/$(FRONTEND_NAME)
FRONTEND_DIRECTORY=frontend
FRONTEND_DOCKERFILE=frontend/Dockerfile

API_NAME=api
API_IMAGE=$(IMAGE_NAMESPACE)/$(API_NAME)
API_DIRECTORY=api
API_DOCKERFILE=api/Dockerfile

OPERATOR_NAME=operator
OPERATOR_IMAGE=$(IMAGE_NAMESPACE)/$(API_NAME)
OPERATOR_DIRECTORY=operator
OPERATOR_DOCKERFILE=operator/Dockerfile

MYSQL_NAME=mysql

HELM_DIRECTORY=helm
HELM_FRONTEND=frontend
HELM_API=api
HELM_OPERATOR=operator
HELM_MYSQL=mysql

.PHONY: *

all: build_images create_cluster load_images add_metallb install output
	
multi_node: build_images create_multinode_cluster load_images add_metallb install output

destroy: teardown_cluster delete_images

build_images: build_frontend build_api build_operator

load_images: load_frontend load_api load_operator

delete_images: delete_frontend delete_api delete_operator

install: install_frontend install_mysql install_api install_operator

uninstall: uninstall_frontend uninstall_api uninstall_mysql uninstall_operator

create_cluster:
	kind create cluster --name $(CLUSTER_NAME)

create_multinode_cluster:
	kind create cluster --name $(CLUSTER_NAME) --config kind-multinode-config.yaml

add_metallb:
	./add_metallb.sh

teardown_cluster:
	kind delete cluster --name $(CLUSTER_NAME)

build_frontend:  
	docker build $(ROOT_DIRECTORY)/$(FRONTEND_DIRECTORY)  -t $(FRONTEND_IMAGE) -f $(ROOT_DIRECTORY)/$(FRONTEND_DOCKERFILE)

build_api:
	docker build $(ROOT_DIRECTORY)/$(API_DIRECTORY) -t $(API_IMAGE) -f $(ROOT_DIRECTORY)/$(API_DOCKERFILE)

build_operator:
	docker build $(ROOT_DIRECTORY)/$(OPERATOR_DIRECTORY) -t $(OPERATOR_IMAGE) -f $(ROOT_DIRECTORY)/$(OPERATOR_DOCKERFILE)

load_frontend:
	kind load docker-image $(FRONTEND_IMAGE) --name $(CLUSTER_NAME)

load_api:
	kind load docker-image $(API_IMAGE) --name $(CLUSTER_NAME)

load_operator:
	kind load docker-image $(OPERATOR_IMAGE) --name $(CLUSTER_NAME)

install_frontend:
	helm install -f $(ROOT_DIRECTORY)/$(HELM_DIRECTORY)/$(HELM_FRONTEND)/values.yaml \
		-f $(ROOT_DIRECTORY)/$(HELM_DIRECTORY)/$(HELM_FRONTEND)/values-dev.yaml \
		$(FRONTEND_NAME) $(ROOT_DIRECTORY)/$(HELM_DIRECTORY)/$(HELM_FRONTEND)

install_mysql:
	helm install $(MYSQL_NAME) oci://registry-1.docker.io/bitnamicharts/mysql --version "11.1.0" --wait \
		--create-namespace --namespace=mysql -f $(ROOT_DIRECTORY)/$(HELM_DIRECTORY)/$(HELM_MYSQL)/values.yaml \
		-f $(ROOT_DIRECTORY)/$(HELM_DIRECTORY)/$(HELM_MYSQL)/values-dev.yaml

install_api:
	helm install -f $(ROOT_DIRECTORY)/$(HELM_DIRECTORY)/$(HELM_API)/values.yaml \
		-f $(ROOT_DIRECTORY)/$(HELM_DIRECTORY)/$(HELM_API)/values-dev.yaml \
		$(API_NAME) $(ROOT_DIRECTORY)/$(HELM_DIRECTORY)/$(HELM_API)

install_operator:
	helm install -f $(ROOT_DIRECTORY)/$(HELM_DIRECTORY)/$(HELM_OPERATOR)/values.yaml \
		-f $(ROOT_DIRECTORY)/$(HELM_DIRECTORY)/$(HELM_OPERATOR)/values-dev.yaml \
		$(OPERATOR_NAME) $(ROOT_DIRECTORY)/$(HELM_DIRECTORY)/$(HELM_OPERATOR)

uninstall_frontend:
	helm uninstall $(FRONTEND_NAME) --wait

uninstall_api:
	helm uninstall $(API_NAME) --wait

uninstall_operator:
	helm uninstall $(OPERATOR_NAME) --wait

uninstall_mysql:
	helm uninstall $(MYSQL_NAME) --wait

delete_frontend:
	- docker rmi $(FRONTEND_IMAGE)

delete_api:
	- docker rmi $(API_IMAGE)

delete_operator:
	- docker rmi $(OPERATOR_IMAGE)

redeploy_frontend: uninstall_frontend build_frontend load_frontend install_frontend

redeploy_api: uninstall_api build_api load_api install_api

redeploy_operator: uninstall_operator build_operator load_operator install_operator

output:
	$(info $(HEADER))
	@echo "Type 'make destroy' to delete the cluster and locally build images"
	@echo "Type 'make teardown_cluster' to only delete the cluster"
	@echo "Type 'make redeploy_api/frontend/operator' to rebuild and redeploy the given component to the cluster"