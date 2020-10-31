# demorestapi
Create a RESTful API for movies(something similar to IMDB) using Golang,Postgresql.

# Info
Have used go mod as dependency manager

# Usage
    1. Clone this repo in your GOPATH 
    2. Run make help to perform suitable operation 
    3. Run Either by Binary or Docker 

# Deployment
    make run-docker

# Features
    1. Add,Edit,Delete Movies By Admin
    2. Users can just view movies
    3. Session and Authentication of Users

# Improvements To be Done
    For Current Implementation Each Search Request Query Goes to DB have to make LRU Cache and save all queries Result, This Will Reduce Load on Database  
    Or We can Create Cache Service(MicroService) or USe some external cache's or look for Search Engine such as Elastic Search

    Pagination To be Implemented, Limiting Data to be sent in one request, we can use scroll Approach too

    Filtering and Sorting so that less and relevant data is passed

    Proper Error handling and returning relevant error codes
    
# Scalabilty
    We can scale API in following ways
### Using any Cloud Service Providers
    1.  Use Docker image and create instance group  
    2.  Create Managed Instance Group and Add Load Balancer 
    3.  Set AutoScaling based on RPS or CPU utilization 

### Deploying on Kubernetes
    1. Create Deployment, Pod consisisting of docker containers 
    2. Expose this Container by creating a Service of Type NodePort or LoadBalancer 
    3. For LB we can go with any cloud provider or metallb for bare metal servers 
    4. Enabling Horizontal Pod AutoScaling 

