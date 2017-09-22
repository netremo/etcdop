# etcdop
ETCD operations (dump/restore to yaml and json files)
Currently supports only YAML file format.

# HOW-TO

Dump entire ETCD server to file:
```
./etcdop -url http://address:port -out filename.yaml 
```
Restore from file to ETCD server:
```
./etcdop -url http://address:port -in filename.yaml 
```