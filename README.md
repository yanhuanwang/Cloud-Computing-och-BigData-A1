# Deployment of User and Expense Services

This README provides instructions for deploying the User and Expense services using Kubernetes. The deployment consists of a PostgreSQL database, Expense Service, and User Service, all managed within a Kubernetes cluster.

## Prerequisites

- Kubernetes cluster
- `kubectl` configured to interact with your Kubernetes cluster
- Docker (to build images for User and Expense services)
- PostgreSQL database

## Overview

The deployment includes the following components:

1. **PostgreSQL Database**: A single PostgreSQL instance used by both the Expense and User services.
2. **Expense Service**: A Go-based REST API for managing expenses.
3. **User Service**: A Go-based REST API for managing user information.
4. **Job to Create User Database**: A Kubernetes job that creates the `userdb` database.

## Kubernetes Deployment Configuration

### Step 1: PostgreSQL Deployment

The PostgreSQL instance is deployed as a Kubernetes `Deployment` and `Service`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres-deployment
  labels:
    app: postgres
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:latest
        ports:
        - containerPort: 5432
        env:
        - name: POSTGRES_USER
          value: "postgres"
        - name: POSTGRES_PASSWORD
          value: "mysecretpassword"
        - name: POSTGRES_DB
          value: "expensedb"
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
      volumes:
      - name: postgres-storage
        persistentVolumeClaim:
          claimName: postgres-pvc
```

A `PersistentVolumeClaim` is also used to provide persistent storage for the PostgreSQL database.

### Step 2: Create Databases Job

A Kubernetes `Job` is used to create the `userdb` database:

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: create-databases-job
spec:
  template:
    spec:
      containers:
      - name: create-databases
        image: postgres:latest
        env:
        - name: PGPASSWORD
          value: "mysecretpassword"
        command: ["/bin/sh", "-c", "until pg_isready -h postgres -U postgres; do echo waiting for database; sleep 2; done; psql -h postgres -U postgres -c 'CREATE DATABASE userdb;'"]
      restartPolicy: OnFailure
```

### Step 3: Deploy Expense Service

The Expense Service is deployed as follows:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: expense-service-deployment
  labels:
    app: expense-service
spec:
  replicas: 1
  selector:
    matchLabels:
      app: expense-service
  template:
    metadata:
      labels:
        app: expense-service
    spec:
      containers:
      - name: expense-service
        image: expense-service:latest
        imagePullPolicy: Never
        ports:
        - containerPort: 8081
        readinessProbe:
          httpGet:
            path: /readiness
            port: 8081
          initialDelaySeconds: 15
          periodSeconds: 10
          timeoutSeconds: 5
```

The `Expense Service` is exposed using a `LoadBalancer` service:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: expense-service
spec:
  selector:
    app: expense-service
  ports:
    - protocol: TCP
      port: 8081
      targetPort: 8081
  type: LoadBalancer
```

### Step 4: Deploy User Service

The User Service is deployed similarly to the Expense Service:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service-deployment
  labels:
    app: user-service
spec:
  replicas: 1
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
      - name: user-service
        image: user-service:latest
        imagePullPolicy: Never
        ports:
        - containerPort: 8082
        readinessProbe:
          httpGet:
            path: /readiness
            port: 8082
          initialDelaySeconds: 15
          periodSeconds: 10
          timeoutSeconds: 5
```

The `User Service` is also exposed using a `LoadBalancer` service:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: user-service
spec:
  selector:
    app: user-service
  ports:
    - protocol: TCP
      port: 8082
      targetPort: 8082
  type: LoadBalancer
```

### Step 5: Persistent Volume Claim

A `PersistentVolumeClaim` is used to store the PostgreSQL data:

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
```

## Deployment Instructions

1. **Deploy PostgreSQL and Persistent Volume Claim**
   ```sh
   kubectl apply -f postgres-service-deployment.yaml
   kubectl apply -f postgres-pvc.yaml
   ```

2. **Run the Job to Create User Database**
   ```sh
   kubectl apply -f create-databases-job.yaml
   ```

3. **Deploy the Expense and User Services**
   ```sh
   kubectl apply -f expense-service-deployment.yaml
   kubectl apply -f user-service-deployment.yaml
   ```

4. **Verify Deployments and Services**

   Verify that all pods are running:
   ```sh
   kubectl get pods
   ```

   Verify that all services are running:
   ```sh
   kubectl get services
   ```
5. **Accessing Services**

   Expense Service: You can access the Expense Service from your web browser using the following URL:

   ```sh
   http://localhost:8081/static/dist/index.html
   ```
   This will open the Expense Service web interface.

   User Service: You can access the User Service from your web browser using the following URL:

   ```sh
   http://localhost:8082/static/dist/index.html
   ```
   This will open the User Service web interface.

   Expense Service: The Expense Service APIs can be accessed using tools like Postman or curl at
   ```sh
   http://localhost:8081.
   ```

   User Service: The User Service APIs can be accessed using tools like Postman or curl at
   ```sh
   http://localhost:8082.
   ```
## License

This project is licensed under the MIT License.

