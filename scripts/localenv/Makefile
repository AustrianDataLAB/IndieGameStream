ROOT_DIRECTORY = ../..
CLUSTER_NAME = indiegamestream

FRONTEND_NAME=frontend
FRONTEND_DIRECTORY=frontend
FRONTEND_DOCKERFILE=frontend/Dockerfile

API_NAME=api
API_DIRECTORY=api
API_DOCKERFILE=api/Dockerfile

HELM_DIRECTORY=helm
HELM_FRONTEND=frontend
HELM_API=api

.PHONY: *

all: build_images create_cluster load_images add_metallb install

multi_node: build_images create_multinode_cluster load_images add_metallb install

teardown: teardown_cluster delete_images

build_images: build_frontend build_api

load_images: load_frontend load_api

delete_images: delete_frontend delete_api

install: install_frontend install_api

uninstall: uninstall_frontend uninstall_api

create_cluster:
	kind create cluster --name $(CLUSTER_NAME)

create_multinode_cluster:
	kind create cluster --name $(CLUSTER_NAME) --config kind-multinode-config.yaml

add_metallb:
	./add_metallb.sh

teardown_cluster:
	kind delete cluster --name $(CLUSTER_NAME)

build_frontend:  
	docker build $(ROOT_DIRECTORY)/$(FRONTEND_DIRECTORY)  -t $(FRONTEND_NAME) -f $(ROOT_DIRECTORY)/$(FRONTEND_DOCKERFILE)

build_api:
	docker build $(ROOT_DIRECTORY)/$(API_DIRECTORY) -t $(API_NAME) -f $(ROOT_DIRECTORY)/$(API_DOCKERFILE)

load_frontend:
	kind load docker-image $(FRONTEND_NAME) --name $(CLUSTER_NAME)

load_api:
	kind load docker-image $(API_NAME) --name $(CLUSTER_NAME)

install_frontend:
	helm install -f $(ROOT_DIRECTORY)/$(HELM_DIRECTORY)/$(HELM_FRONTEND)/values-dev.yaml \
		$(FRONTEND_NAME) $(ROOT_DIRECTORY)/$(HELM_DIRECTORY)/$(HELM_FRONTEND)

install_api:
	helm install -f $(ROOT_DIRECTORY)/$(HELM_DIRECTORY)/$(HELM_FRONTEND)/values-dev.yaml \
		$(FRONTEND_NAME) $(ROOT_DIRECTORY)/$(HELM_DIRECTORY)/$(HELM_FRONTEND)

uninstall_frontend:
	helm uninstall $(FRONTEND_NAME)

uninstall_api:
	helm uninstall $(API_NAME)

delete_frontend:
	docker rmi $(FRONTEND_NAME)

delete_api:
	docker rmi $(API_NAME)