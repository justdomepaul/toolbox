# toolbox

### Golang toolbox list

- array (only for golang 1.18 upper)
- base58
- config
- database
- errorhandler
- firebase
- grpc
- interceptor
- jwt
- key
- restful
- services
- shorten
- shutdown
- spannertool (part only for golang 1.18 upper)
- stringtool
- timestamp
- zap


### GCP Simulator
##### [Cloud Storage](https://cloud.google.com/storage) : oittaa/gcp-storage-emulator
```yaml
storage:
    image: oittaa/gcp-storage-emulator
    environment:
      PORT: 9023
    ports:
      - "9023:9023
    command: [ "start", "--in-memory", "--default-bucket=staging.project.appspot.com" ]
```
##### [Cloud Spanner](https://cloud.google.com/spanner) : gcr.io/cloud-spanner-emulator/emulator
```yaml
spanner:
    image: gcr.io/cloud-spanner-emulator/emulator
    ports:
      - "9010:9010"
      - "9020:9020"
```
##### [Cloud Pub/Sub](https://cloud.google.com/pubsub) : justdomepaul/gcp-pubsub-simulator
```yaml
spanner:
    image: justdomepaul/gcp-pubsub-simulator
    ports:
      - "9000:9000"
```