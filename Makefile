.PHONY: build push create delete

build:
	docker build -t ${account}/kub-demo-auths ./auth-api
	docker build -t ${account}/kub-demo-users ./users-api
	docker build -t ${account}/kub-demo-tasks ./tasks-api
	docker build -t ${account}/kub-demo-frontend ./frontend

push:
	docker push ${account}/kub-demo-auths
	docker push ${account}/kub-demo-users
	docker push ${account}/kub-demo-tasks
	docker push ${account}/kub-demo-frontend

create:
	kubectl apply -f=kubernetes/auth-deployment.yml -f=kubernetes/auth-service.yml
	kubectl apply -f=kubernetes/users-deployment.yml -f=kubernetes/users-service.yml
	kubectl apply -f=kubernetes/tasks-deployment.yml -f=kubernetes/tasks-service.yml
	kubectl apply -f=kubernetes/frontend-deployment.yml -f=kubernetes/frontend-service.yml

delete:
	kubectl delete -f=kubernetes/auth-deployment.yml -f=kubernetes/auth-service.yml
	kubectl delete -f=kubernetes/users-deployment.yml -f=kubernetes/users-service.yml
	kubectl delete -f=kubernetes/tasks-deployment.yml -f=kubernetes/tasks-service.yml
	kubectl delete -f=kubernetes/frontend-deployment.yml -f=kubernetes/frontend-service.yml

service:
	minikube service frontend-service

