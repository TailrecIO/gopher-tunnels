
# Gopher Tunnels
A local tunnel built ~~with~~by Gopher.

If ngrok is blocked in your network, this one is an alternative that you can deploy 
to your own AWS account in a few steps.
This application runs on serverelss 100%. I mean all servers will be managed by Amazon 100%.
You may pay very little if there is no traffic coming to your service. (Amazon may charge for configuration).

This project uses [serverless framework](https://serverless.com/framework/docs/providers/aws/guide/). So, make sure that you're familiar with that a bit before get going. 

## Requirements
- [Go](https://golang.org/doc/install)
- [Go dep](https://github.com/golang/dep)
- [Node.js](https://nodejs.org/) 
- [Severless](https://serverless.com/framework/docs/providers/aws/guide/installation/)

### You may want to install this library to your $GOPATH
This is optional and these library can be installed locally to vendor directory.
- go get -u github.com/aws/aws-sdk-go
- go get -u github.com/pkg/errors
- go get -u github.com/rs/xid
- go get -u golang.org/x/crypto/nacl/box

## What are included in this project?

### Server (Lambda functions)
Binary path: `bin/server`
- **gopher_register** - the client program (gopher) has to call this function when it starts.  
- **gopher_webhook** - this function handle your webhook. The endpoint will be provided to you when you run the client.
- **gopher_respond** - this function handle the response returning from the client. 
- **gopher_cleanup** - work in progress 

### Client
Binary path: `bin/client`
- **gopher** - the client that you can distribute to your team or anyone you want.

### Tools
Binary path: `bin/tools`
- **config** - this tool will generate a client configuration file from application.yml.
- **echo_server** - this tool is really useful (not just for this service).   
  It is an http server that returns exactly what it receives like an echo service.
  The nice part of this tool is that it echo everything to the console. So, you can inspect 
  the HTTP request such as path, method, headers, and body. 

# How to configure the project

Please copy application.template.yml to application.yml
Edit that file to make sure that all attributes are correct.
Note that you don't know `baseApiEndpoint` until you deploy the API it to API Gateway.
This attribute can be left empty when you deploy your serverless because it's required by 
the client not the server. 

# How to build the project

Make sure that everything is installed to vendor directory  
`$ dep ensure`

Compile the project  
`$ make`

# How to deploy the project

## Deploying a server
You can deploy simply by running a simple command like the following. Thanks to serverless framework.    
`sls deploy` 

You can override some attributes such as stage and region via command line.   
`sls deploy --stage prod --region us-west-1` 

You can deploy individual function by providing the function name
`sls deploy function --function gopher_register`
 
## Packaging a client

Assume that your current directory is at the project root. Otherwise, you will get an error.  
If you want to create a client for the default stage (dev), you simply run `./dist.sh`.  
`dist.sh` will create the dist directory, generate a client configuration from your application.yml, and
copy all related binary files. You can just zip this directory and distribute it.

However, you can specify the stage name by providing the stage name as the first argument.
I'm a bash script noob and I don't care to improve it. So, contribution is very welcome :)

Note that `gopher` will read the configuration named `gopher.yml` if you don't provide any argument. 
Please see the usage section for more details.

### Usage
You can override `host`, `mode`, and `port` by providing the arguments for those attributes.
Otherwise, the values from a configuration file will be used.
```
$ ./gopher -h
Usage of ./gopher:
  -host string
    	target host
  -mode string
    	webhook mode: sync, async
  -port string
    	target port
```
If your configuration file is not `gopher.yml`, you have to provide the configuration 
file as the last argument as the following.
`$ ./gopher gopher.dev.yml`





 