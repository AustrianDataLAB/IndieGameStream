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
IMAGE_LABEL = localenv

FRONTEND_NAME=frontend
FRONTEND_IMAGE=$(IMAGE_NAMESPACE)/$(FRONTEND_NAME):$(IMAGE_LABEL)
FRONTEND_DIRECTORY=frontend
FRONTEND_DOCKERFILE=frontend/Dockerfile

API_NAME=api
API_IMAGE=$(IMAGE_NAMESPACE)/$(API_NAME):$(IMAGE_LABEL)
API_DIRECTORY=api
API_DOCKERFILE=api/Dockerfile

OPERATOR_NAME=operator
OPERATOR_IMAGE=$(IMAGE_NAMESPACE)/$(OPERATOR_NAME):$(IMAGE_LABEL)
OPERATOR_DIRECTORY=operator
OPERATOR_DOCKERFILE=operator/Dockerfile

MYSQL_NAME=mysql
STUNNER_NAME=stunner

HELM_DIRECTORY=helm
HELM_FRONTEND=frontend
HELM_API=api
HELM_OPERATOR=operator
HELM_MYSQL=mysql
HELM_STUNNER=stunner

.PHONY: *

# Creates a single-node cluster and installs all components on it
all: build_images create_cluster load_images add_metallb install help

# Creates a cluster with three nodes and installs all components on it
multi_node: build_images create_multinode_cluster load_images add_metallb install help

# Removes the cluster and all local images
destroy: teardown_cluster delete_images

# Builds images
build_images: build_frontend build_api build_operator

# Loads images into the cluster
load_images: load_frontend load_api load_operator

# Deletes built images locally
delete_images: delete_frontend delete_api delete_operator

# Creates all our components
install: install_mysql install_stunner install_operator install_api install_frontend

# Cleans up the cluster, removing all components
uninstall: uninstall_frontend uninstall_api uninstall_mysql uninstall_stunner uninstall_operator

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
		--set-string appConfig.apiUrl=http://$$(kubectl get svc api -n api -o jsonpath='{.status.loadBalancer.ingress[0].ip}'):$$(kubectl get svc api -n api -o jsonpath='{.spec.ports[0].port}') \
		$(FRONTEND_NAME) $(ROOT_DIRECTORY)/$(HELM_DIRECTORY)/$(HELM_FRONTEND)

install_mysql:
	helm repo add mysql-operator https://mysql.github.io/mysql-operator/
	helm repo update
	helm install $(MYSQL_NAME)-operator mysql-operator/mysql-operator --version "2.1.3" --wait \
		--create-namespace --namespace=$(MYSQL_NAME)-operator
	helm install $(MYSQL_NAME) mysql-operator/mysql-innodbcluster --version "2.1.3" --wait \
        --create-namespace --namespace=$(MYSQL_NAME) \
		-f $(ROOT_DIRECTORY)/$(HELM_DIRECTORY)/$(HELM_MYSQL)/values.yaml \
		-f $(ROOT_DIRECTORY)/$(HELM_DIRECTORY)/$(HELM_MYSQL)/values-dev.yaml

install_api: wait_for_mysql
	helm install -f $(ROOT_DIRECTORY)/$(HELM_DIRECTORY)/$(HELM_API)/values.yaml \
		-f $(ROOT_DIRECTORY)/$(HELM_DIRECTORY)/$(HELM_API)/values-dev.yaml \
		$(API_NAME) $(ROOT_DIRECTORY)/$(HELM_DIRECTORY)/$(HELM_API)

install_operator:
	$(MAKE) -C $(ROOT_DIRECTORY)/$(OPERATOR_DIRECTORY) deploy IMG=$(OPERATOR_IMAGE)

install_stunner:
	helm repo add stunner https://l7mp.io/stunner
	helm repo update
	helm dependency build $(ROOT_DIRECTORY)/$(HELM_DIRECTORY)/$(HELM_STUNNER) --skip-refresh
	helm install $(STUNNER_NAME) $(ROOT_DIRECTORY)/$(HELM_DIRECTORY)/$(HELM_STUNNER) \
		--create-namespace --namespace=$(STUNNER_NAME)

uninstall_frontend:
	helm uninstall $(FRONTEND_NAME) --wait

uninstall_api:
	helm uninstall $(API_NAME) --wait

uninstall_operator:
	$(MAKE) -C $(ROOT_DIRECTORY)/$(OPERATOR_DIRECTORY) undeploy

uninstall_mysql:
	helm uninstall $(MYSQL_NAME) -n $(MYSQL_NAME) --wait
	helm uninstall $(MYSQL_NAME)-operator -n $(MYSQL_NAME)-operator --wait

uninstall_stunner:
	helm uninstall $(STUNNER_NAME) --namespace=$(STUNNER_NAME)

delete_frontend:
	- docker rmi $(FRONTEND_IMAGE)

delete_api:
	- docker rmi $(API_IMAGE)

delete_operator:
	- docker rmi $(OPERATOR_IMAGE)

redeploy_frontend: uninstall_frontend build_frontend load_frontend install_frontend

redeploy_api: uninstall_api build_api load_api install_api redeploy_frontend

redeploy_operator: uninstall_operator build_operator load_operator install_operator

wait_for_mysql:
	@while true; do \
	    POD_STATUS=$$(kubectl get pod -l app.kubernetes.io/component=router -n $(MYSQL_NAME) --no-headers -o custom-columns=":status.phase" 2>/dev/null); \
	    if [ "$$POD_STATUS" ]; then \
	        echo "MySQL router is created with status: $$POD_STATUS"; \
	        break; \
	    else \
	        echo "Waiting for MySQL router to be created..."; \
	        sleep 2; \
	    fi \
	done
	kubectl wait --for=condition=ready pod -l app.kubernetes.io/component=router -n $(MYSQL_NAME) --timeout=600s

help:
	$(info $(HEADER))
	@echo "----------------------------------------------------------------------"
	@echo "Type 'make destroy' to delete the cluster and locally build images"
	@echo "Type 'make teardown_cluster' to only delete the cluster"
	@echo "Type 'make redeploy_api/frontend/operator' to rebuild and redeploy the given component to the cluster"
	@echo "Type 'make help' to get this information"
	@echo "Other useful targets can be found in the Makefile"
	@echo "----------------------------------------------------------------------"
	@echo "The frontend is reachable under http://$$(kubectl get svc frontend -n frontend -o jsonpath='{.status.loadBalancer.ingress[0].ip}')"
	@echo "The API is reachable under http://$$(kubectl get svc api -n api -o jsonpath='{.status.loadBalancer.ingress[0].ip}'):$$(kubectl get svc api -n api -o jsonpath='{.spec.ports[0].port}')"
	@echo "----------------------------------------------------------------------"