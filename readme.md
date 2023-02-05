# go-url-shortener-example
Example of url shortener.

Based on this: https://github.com/mahadevans87/short-url 
Rebuilt into my own patterns and coding style for my own education.
At the moment does not work, but I'll fix it soon.

## Running

### Running

Use `go run .` from the folder that contains `main.go`.  

I used such commands to run redis and debug this way of running:  
- `docker network create some-network`
- `docker run --network some-network -p 6379:6379 --name some-redis -d redis`
- `docker run -it --network some-network --rm redis redis-cli -h some-redis`