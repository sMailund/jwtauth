# jwt auth
Sample project to learn microservice authentication with JWTs in golang.

## about
This repository contains a test project for handling authentication between microservices through the use of JWTs.
In order to reduce traffic to the authentication service, 
authentication towards the services uses a signed JWT token from the auth service, 
which can be verified without further calls to the auth service in subsequent calls to other services.

## technical details
On startup, the auth service generates the RSA keypair 
(currently not persisted, as this project is not in actual use).
Users must POST to the `/login/` endpoint with username and password to obtain a signed JWT, 
which is stored in the cookies.
Additionally, there's a  `/create/` endpoint for creating new users.

Internal communication between services is over gRPC.
Currently, the only useful RPC is for other services to obtain the public key from the auth server, 
which is then used to verify incoming JWTs.

Each service is structured with the [onion architectural pattern](https://medium.com/@shivendraodean/software-architecture-the-onion-architecture-1b235bec1dec), 
to keep coupling low and allow for hot-swapping of external resources such as database connections.
