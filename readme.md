# go-url-shortener-example
Example of url shortener.

Based on this: https://github.com/mahadevans87/short-url 
Rebuilt into my own patterns and coding style for my own education.
At the moment works, but some things needed:
- another web-framework;
- graceful shutdown;
- custom logging;
- dockerfile etc etc;
- docker compose.

## Generate

API boilerplate code is generated using `oapi-codegen` tool from the `openapi.yaml` file.  
It's great tool that makes your actual API reflect the documentation.  

Get it there:  
`https://github.com/deepmap/oapi-codegen`  

And make sure that your `GOPATH/bin` path presents in `PATH` variable.  

Use this command to generate the `api.go` file:  
- `oapi-codegen -package=api --generate=types,chi-server api/openapi.yaml > internal/api/api.go`  

- `oapi-codegen -package=api --generate=types,chi-server api/petstore.yaml > internal/api/api2.go`
## Running

### Running

Use `go run .` from the folder that contains `main.go`.  

I used such commands to run redis and debug this way of running:  
- `docker network create some-network`
- `docker run --network some-network -p 6379:6379 --name some-redis -d redis`
- `docker run -it --network some-network --rm redis redis-cli -h some-redis`