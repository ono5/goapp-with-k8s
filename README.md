# goapp-with-k8s

This app works on Kubernetes.

so, if you want to work this app, you needs to set up Kubernetes Environment.

## Step by Step to work this app

1. Clone repository

```
git clone https://github.com/ono5/goapp-with-k8s.git
```

2. Install Kubernetes and Docker

- [Docker Install](https://docs.docker.com/get-docker/)
- [Kubernetes Install](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

3. Register Docker Hub

To use Docker Image, Let's register [Docker Hub](https://hub.docker.com/).

And then, you create the following repository.

- <Docker Hub のアカウント名>/kub-demo-auths
- <Docker Hub のアカウント名>/kub-demo-users
- <Docker Hub のアカウント名>/kub-demo-tasks
- <Docker Hub のアカウント名>/kub-demo-frontend

4. Create Docker Image by Dockerfile

You also need to create docker image from each Dockerfile.

```
make build account=<DockerHubアカウント>
```

5. Push to Docker Hub

Push your perfect docker image to your repo on Docker Hub

```
make push account=<DockerHubアカウント>
```

6. Create Service & Deployment

Create Service and Deployment of Kubernetes

```
make create
```

7. Start frontend service

To start frontend service, you do the below command.

```
make service
```

8. Delete Service & Deployment

Delete Service and Deployment of Kubernetes

```
make service
```
