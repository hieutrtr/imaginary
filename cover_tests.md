# Imaginary Test cases
### Switching Object to Block
* Stop Ceph services (Monitor, Storage)
* Switch Imaginary using raw image from FS.
### Backup
* Tracking events of uploading from imaginary to Kafka with payload
```
{
  "pool":"POOL_NAME"
  "OID":"OBJECT_ID"
  "CREATED_TIME":TIMESTAMP
}
```
* Consume Kafka events and backup to Filesystem
* Cross-check by counting Objects
# Restore
* Consume Kafka events and restore from Filesystem to Ceph Object Storage
* Cross-check by counting Objects
