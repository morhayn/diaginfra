Directory for config files
Create config file - config.yaml

```
user: "username"
ssh_port: "22"
countlog: 300
list_urls:
  - "list url for testing"
logs:
  tomcat: "/var/log/tomcat/"
  elasticsearch: "/var/log/elasticsearch.log"
hosts:
  - name: "name server"
    ip: "ip server"
    list_ports:
      - "22"
      - "number port"
      - "number port"
    list_service:
      - "service"
      - "Tomcat:admin:pass:port"
      - "mongod"
      - "Mongo:admin:pass"
      - "elasticsearch"
      - "Elastic"
      - "demon1"
    wars:
      - "list war files"
```
user - user name for conect ssh to servers. Ssh key get from ~/.ssh/
ssh_port - port for connect to servers
countlog - count tail log file
logs - path to logs file
list_usrls - list urls http and https for test response status code
list_ports - list ports for checks
list_service -
  List service to check ssh commands
  Add new command -
   pkg/sshcmd.go add command
   pkg/handl.go add parse result command
```
  "tomcat" - systemctl is-active tomcat
  "elasticsearch" - sysytemctl is-active elasticasearch
  "default" - systemctl is-active default
  "Tomcat:admin:pass:port" - curl -u ADMIN:PASS http://127.0.0.1:PORT/manager/text/list
  "Jar:jarservice" - sudo systemctl is-active JARSERVICE
  "MondoDb:admin:pass" - mongo -u ADMIN -p "PASS"  --eval 'db.stats()'
  "Hazelcast:admin:pass" - curl --data "ADMIN&PASS" --silent "http://127.0.0.1:5701/hazelcast/rest/management/cluster/state"
  "Elastic" - curl -X GET http://127.0.0.1:9200/_cluster/health
  "Cassandra" - nodetool status
  "Postgress" - pg_lsclusters | awk 'FNR > 1 {print $4}'
  "Kafka" - export KAFAK_OPTS='-Djava.security.auth.login.config=/etc/kafka/kafka_jaas.conf'; /d01/kafka/bin/kafka-topics.sh --list --zookeeper localhost:2181
  "Rabbit" - rabbitmqctl status
  "Ceph" - ceph status | awk '/health/ {print $2}'
  "Docker" - docker ps --format '{"name":"{{.Names}}", "status":"{{.Status}}"}'
```