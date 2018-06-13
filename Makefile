.PHONY: deploy-operator clean-operator clean-cassandra build-app docker-app-build minikube-deploy minikube-start docker-clean

KUBE_INSTALL := $(shell command -v kubectl 2> /dev/null)
REPO=niravpatel/k8s-cassandra-workshop
KEY_SPACE=accountsapi
CASSANDRA_CLUSTER_URL=cassandra.default.svc.cluster.local
BUILD_PATH=build

check:
ifndef KUBE_INSTALL 
    $(error "kubectl is not available please install from https://kubernetes.io/")
endif

deploy-operator: check
	kubectl apply -f manifest/operator.yaml

provision-cassandra:
	kubectl create -f manifest/cluster.yaml

clean-operator: check
	kubectl delete -f manifest/operator.yaml

clean-cassandra:
	kubectl delete -f manifest/cluster.yaml

build-app:
	@rm -rf ./build
	@mkdir -p $(BUILD_PATH)
	GOOS=linux go build -o $(BUILD_PATH)/main . 

docker-app-build: check build-app
	@eval $$(minikube docker-env) ;\
	docker build -t $(REPO):latest .

docker-clean:
	docker rmi $(REPO)

minikube-deploy: 
	kubectl run web --env="KEY_SPACE=$(KEY_SPACE)" --env="CASSANDRA_CLUSTER_URL=$(CASSANDRA_CLUSTER_URL)" --image=$(REPO):latest --port=8080 --image-pull-policy=Never
	kubectl expose deployment web --target-port=8080 --type=NodePort

minikube-start: 
	minikube start --cpus 4 --memory 4096 --vm-driver hyperkit --kubernetes-version v1.9.4
