Gateway for client

This system is based made out of 3 part, all written in Golang.
I used Golang because of its great concurency ability and because it can be compiled.

1) A configurable reverse proxy as a client gateway.
2) a server gateway.
3) a mock (and very dumb) service.

1) Client Gateway
This is built as a reverse proxy you can configure with a json file.

To add an endpoint, add a new object in the config file.
{
  "path":"/product/{id:[0-9]+}",
  "redirect":"http://localhost:8081/company/10/product/{id}",
  "timeout":10,
  "keepAlive":10,
  "methods":["GET"],
  "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOjEsInJvbGUiOiJzdXBlcl91c2VyIiwiZXhwIjoyNTIyNzcwODgyfQ.-wDTI9hDR-hprujVTu-ZJBwzdu7M4xF9BZNMRYhdHGg"
}

To get them working, the service has to be restarted, but a endpoint could be added to add new endpoints without restarting the server. This would add complexity but flexibility.
Each path has a proxy redirect to the server gateway. This is built in a preconceived notion that a load balancing mechanism is in front of it.
If the object has the parameter dumb_proxy set to true, only the scheme and host needs to be set in the redirect, the original requestURI would follow.
Each redirect can have multiple http methods set to it.
You can see that I have used an hardcoded jwt token. In this example, my expiration date is very high for testing purposes.
In the future we would have a way to login and automatically get one from the client service.
To secure the client even more, we could use client side SSL.

2) Server gateway
To be honest, I could have used the same code then for the client gateway and it would have worked.
But I wanted to give more flexibility so regular handlers are used here so many services could be called for each endpoint.
I could have made a hybrid version to get the best of both world but decided against for time constraint
I used Negroni and gorilla mux, two open sourced micro web framework, to help, in a package called negronimux.
I have built 2 handlers for product (GET AND POST), and many more could be built. These handler calls the mock service.
I have built a very basic Authentication middlewares using JWTToken, that validates the token and then calls a mock authenticate service which always return ok.

3) Mockservice
This is a dumb service just for testing purpose. Its had 3 endpoint, 2 for the product (POST and GET), and a validate auth which always return status 200.

All services could have additionnal logging. We could horizontaly scale this if a proper load balancing tool was added.

Unit testing. Unfortunately I didn't add any unit testing because of time constraint. This could be easily added after. I swear I can test code and would never have untested code in produciton.

------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

To start the project, in 3 different terminal from the root of the project (where the readme is) do:

./cmd/client/bin -port 8080 -config "conf.json"
./cmd/server/bin -port 8081 -jwtsecret "secret" -service "http://localhost:8082"
./cmd/mockservice/bin -port 8082

Flags can be modified but I suggest not to. More endpoint can easily be added to the config file, but no endpoints will be set for them in the server and mock service. They will return 404.
In the future, this will have its own docker files to make containers to easily deploy it.

To test the whole set up, hit "http://localhost:8080/product/10" with GET, you can try with many different ids and you will get a different result. this logic is in the mock service.

Hit "http://localhost:8080/product" with POST (no body needed, no logic was added here) and you should get a 201 status code. To show that it worked we send back the Body to the client. This does not affect the mock service.

You should see a hit on all services.
To be able to hit directly the server gateway (http://localhost:8081/company/10/product/5), you need to set up the Authorization header with a Bearer token set for now in the client config file.
There is no protection for the mock service but we could easily make it only accessible for the server gateway.
