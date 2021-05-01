# kubedock

Kubedock is an experimental implementation of the docker api that will orchestrate containers into a kubernetes cluster, rather than running containers locally. The main driver for this project is to be able improve running [testcontainers-java](https://www.testcontainers.org) enabled unit-tests in k8s, without the requirement of running docker-in-docker within resource heavy containers.

The current implementation is limited, but able to run containers that just expose ports, or copy resources towards the container. Containers that 'just' expose ports, require logging and copy resources to running containers will probably work. If the container mounts volumes, they are not supported yet.

Running this locally with a testcontainers enabled unit-test requires to run kubedock (`make run`). After that start the unit tests in another termninal with the below environment variables set, for example:

```bash
export TESTCONTAINERS_RYUK_DISABLED=true
export DOCKER_HOST=tcp://127.0.0.1:8080
mvn test
```

The default configuration for kubedock is to orchestrate in the default kubernetes namespace, this can be configured with the `NAMESPACE` environment variable (or via the -n argument). The service requires permissions to create Deployments in the namespace.