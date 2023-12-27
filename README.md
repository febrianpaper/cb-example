# Circuit Breaker Example

This project is an example of how to use the Circuit Breaker pattern in Golang application.

# How to run
## Run Dummy Server
### Infobib Server
Open New terminal and type this commands:
```bash
$ cd infobib-server
$ go mod vendor
$ make run 
```

### Halosis Server
Open New terminal and type this commands:
```bash
$ cd halosis-server
$ go mod vendor
$ make run 
```

## Run Client
Before we can run client, we need to run Arango DB localy, and set the environment variable `ARANGO_USER`, `ARANGO_PASSWORD`, `ARANGO_HOST` and `ARANGO_DB_NAME` in app.env file. After that, we can type this commands:
```bash
$ cd client
$ go mod vendor
$ make run
```

## Send Message
### Using curl
``` bash
$ curl --location 'http://localhost:8888/v1/sms' \
--header 'Content-Type: application/json' \
--data '{
    "phone_number": "08123941231",
    "message": "Ini pesan sms",
    "template_name": "template-coba",
    "template_data": ["paper.id"]
}'
```