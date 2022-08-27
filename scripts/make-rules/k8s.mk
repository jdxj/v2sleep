KUBE_CONFIG := --kubeconfig=$(DEPLOY)/kube_config.yaml
NAMESPACE := v2sleep

.PHONY: k8s.create.ns
k8s.create.ns:
	@kubectl $(KUBE_CONFIG) create namespace $(NAMESPACE)

.PHONY: k8s.create.secret
k8s.create.secret:
	@kubectl $(KUBE_CONFIG) -n $(NAMESPACE) create secret generic $(NAMESPACE)-conf --from-file=$(DEPLOY)/conf.yaml

.PHONY: k8s.create.deploy
k8s.create.deploy:
	@kubectl $(KUBE_CONFIG) -n $(NAMESPACE) create -f $(DEPLOY)/deployment.yaml

.PHONY: k8s.create.service
k8s.create.service:
	@kubectl $(KUBE_CONFIG) -n $(NAMESPACE) create -f $(DEPLOY)/service.yaml

.PHONY: k8s.set.image
k8s.set.image:
	@kubectl $(KUBE_CONFIG) -n $(NAMESPACE) set image deployment $(NAMESPACE)-dep j$(NAMESPACE)-c=$(DOCKER_TAG)

.PHONY: k8s.apply.secret
k8s.apply.secret:
	@kubectl $(KUBE_CONFIG) -n $(NAMESPACE) scale deployment $(NAMESPACE)-dep --replicas=0
	@kubectl $(KUBE_CONFIG) -n $(NAMESPACE) delete secrets $(NAMESPACE)-conf
	@kubectl $(KUBE_CONFIG) -n $(NAMESPACE) create secret generic $(NAMESPACE)-conf --from-file=$(DEPLOY)/conf.yaml
	@kubectl $(KUBE_CONFIG) -n $(NAMESPACE) scale deployment $(NAMESPACE)-dep --replicas=1

.PHONY: k8s.scale.%
k8s.scale.%:
	@kubectl $(KUBE_CONFIG) -n $(NAMESPACE) scale deployment $(NAMESPACE)-dep --replicas=$*
