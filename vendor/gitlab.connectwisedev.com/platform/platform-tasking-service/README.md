<p align="center">
<img width=25% align="middle" src="docs/images/continuum-logo.svg">
<img width=25% align="middle" src="docs/images/Go-Logo_Blue.png">
</p>

# platform-tasking-service
This repo is used for developing DM Tasking functionality

### Repo contact persons:
- Lokesh Jain <lokesh.jain@connectwise.com>
- Irfan Shaikh <irfan.s@connectwise.com>
- Akshaya Kaushik <Akshaya.Kaushik@ConnectWise.com>

### REST API Endpoints are documented here
[REST API Endpoints](https://gitlab.connectwisedev.com/platform/platform-tasking-service/-/blob/master/api/swagger.yaml)

### SDD
[Design Document](https://continuum.atlassian.net/wiki/spaces/C2E/pages/312705070/Continuum+Task+Sequencing+Solution+2.0+SDD)

### Service registry page
[Services and URLs](https://continuum.atlassian.net/wiki/spaces/C2E/pages/221412064/Desktop+Management+2.0+Service+Registry)

### Ansible
- [CI](http://ci.corp.continuum.net:8080/job/dev_platform-tasking-service/) 
- [DT](http://dtansible.corp.continuum.net:8080/job/dt-ansible-microservices-deploy/) 
- [QA](http://qaansible.corp.continuum.net:8080/job/qa-ansible-microservices-deploy/) | [QA promotion](http://qaansible.corp.continuum.net:8080/job/qa-ansible-microservices-deploy-versioning/) 
- [Artifacts](http://artifact.corp.continuum.net:8081/artifactory/webapp/#/builds/dev_platform-tasking-service/)

### Performance testing
[Performance testing wiki](https://continuum.atlassian.net/wiki/spaces/C2E/pages/1963951651/Tasking+MS+Performance+Results)

### Kafka Topic contracts
[Kafka](https://gitlab.connectwisedev.com/platform/platform-api-model/-/tree/master/clients/model/Golang/resourceModel/tasking)

### Logs and Monitoring
- PROD [monitoring](https://rad43678.live.dynatrace.com/#processgroupdetails;id=PROCESS_GROUP-2534F6FCF5147FC1;gtf=l_72_HOURS;gf=all)

### Quality
[Pulse](http://pulse.corp.continuum.net/#/dashboard/5d89b74ead4136080e596abf) | [SonarQube](http://codescan.continuum.net/dashboard?id=platform-tasking-service)

### Service dependencies
- [Agent MS](https://gitlab.connectwisedev.com/platform/agent-service)
- [Asset MS](https://github.com/ContinuumLLC/platform-asset-service)
- [Scripting MS](https://gitlab.connectwisedev.com/platform/platform-scripting-service)
- [Dynamic Groups MS](https://gitlab.connectwisedev.com/platform/platform-dynamicgroup-service)
- [Entitlement MS](https://github.com/ContinuumLLC/platform-entitlement-service)
- [ITS webapi](https://github.com/ContinuumLLC/rmm-its-webapi)
- [GraphQL](https://github.com/ContinuumLLC/rmm-device-graphql-service)


### Used infrastructure components
- Cassandra
- Kafka
- Zookeeper
- In-memory [cache](https://github.com/bradfitz/gomemcache)
