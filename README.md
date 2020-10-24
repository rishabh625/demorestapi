# fyndtest
Create a RESTful API for movies(something similar to IMDB).

# Info
Have used go mod as dependency manager

# Usage
    1. Clone this repo in your GOPATH <br/>
    2. Run make help to perform suitable operation <br/>
    3. Run Either by Binary or Docker <br/>

# Improvements To be Done
    For Current Implementation Each Search Request Query Goes to DB have to make LRU Cache and save all queries Result, This Will Reduce Load on Database  
    Or We can Create Cache Service(MicroService) or USe some external cache's or look for Search Engine such as Elastic Search

    Pagination To be Implemented, Limiting Data to be sent in one request, we can use scroll Approach too

    Filtering and Sorting so that less and relevant data is passed

    Proper Error handling and returning relevant error codes
    
# Scalabilty
    We can scale API in following ways
### Using any Cloud Service Providers
    1.  Use Dockerfile to create instance group <br/> 
    2.  Create Managed Instance Group and Add Load Balancer <br/>
    3.  Set AutoScaling based on RPS or CPU utilization <br/>

### Deploying on Kubernetes
    1. Create Deployment, Pod consisisting of docker containers <br/>
    2. Expose this Container by creating a Service of Type NodePort or LoadBalancer <br/>
    3. For LB we can go with cloud provider or metallb for bare metal servers <br/>
    4. Horizontal Pod AutoScaling <br/>

# Problems
    1. We might get different Latency for request from different region <br/>
    Soln : Distribute instance group among different geo-locations region or put it behind Geo-Location Based LB.<br/>
    2. When we implement Caching in memory we might overshoot memory or not properly manage it there can be more and more GC pauses <br/>
    Soln :- Either to use Cache's like Redis, or create own Cache as a service where GC is considered in development. <br/>
    3. Searching can be Time Consuming <br/>
    Soln:- If our application becomes Search intensive ,we should use combination of cache and powerfull search Engine like Elastic Search <br/>
    4. Data Transfer Issues <br/>
    Soln:- Since Pagination, Filtering is not implemented we will face issues as data increases by 5x, Implementing This feature will help us command over data transfer.<br/>
    However ther can be issues of network which get chokes based on data,Solution to this can be implementing compression for data transfer or using alternative to json over HTTP kile protobuff over gRPC which have built in compression feature and multi request over single connection <br/>
    5. Code Maintainability<br/>
    Soln: Create generalized object schema for Update,Create,Delete, currently that is not implemented <br/>

