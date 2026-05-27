#!/usr/bin/make -sf

.DEFAULT := install
.PHONY: all clean test install uninstall

all:
	$(MAKE) -C mrprober

docker clean test:
	$(MAKE) -C mrprober $@

# Management on K8S
install:
	kubectl apply -f ./k8s --recursive

uninstall:
	kubectl delete -f ./k8s --recursive
