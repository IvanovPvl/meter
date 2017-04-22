# meter
Meter is a distributed disk space monitoring system. It consists of an agent and a server.

## Agent
It is installed on the servers for which the monitoring will be performed.

#### Build
```
go build run_agent.go
```

#### Run
```
go run run_agent -port=3000
```

Flags:
* port

Information from agent is available by request:
```
curl <agent_host>:<agent_port>
```

## Server
Aggregates data from all agents.

#### Build
```
go build run_server.go
```

#### Run
Create config.yml file and specify hosts on which agents are installed and port for server.
```
go run run_server
```
Information from server is available by /api/df route.