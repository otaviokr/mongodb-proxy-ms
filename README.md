# docker-proxy-ms

## Overview

This application is a microsservice used to give to other services access to the MongoDB database.

## Run Only this container

```bash
# Mongo DB should be running
docker run -d -p 27017-27019:27017-27019 --name mongodb --network=mongonet mongo

# Build the image of this application
docker build -t docker-proxy-ms:1.0 .

# Start a container with this application (pay attention to use the same tag as the previous command for the image)
docker run -d -p 8080:8080 --env GIN_MODE=debug --env MONGODB_HOST=mongodb --name mongo-proxy --network=mongonet docker-proxy-ms:1.0
```
## Run this container together with a MongoDB container

```bash
# Use docker-compose
docker-compose up > output.log &

# To stop it.
docker-compose down
```

## Features

These are the functionalities available

### Insert (/insert/<db>/<collection>)

You must send via POST a JSON object that represents a document. That document will be inserted in **collection**, in **database**.

It returns the ObjectID of the newly inserted document.

### Find (/find/<db>/<collection>)

If you provide a JSON object with filters, it will return all documents from **collection** in **database** that match the filter provided. If no filter is given, all documents in **collection** will be returned.

The method must be POST, and it returns the documents as a JSON object.

### Update (/update/<db>/<collection>)

You must provide a filter to define the document(s) that will be updated (as a JSON object), as well as the fields with the new values (as a separated JSON object).

The method must be POST.

### Health

This is a simple GET request, with no parameters, that will return the available collections in MongoDB, if the database is up and running.

## Configuration

The application expects 4 parameters:

- **hostname**
- **port**
- **username** is optional, and may be defined only if the MongoDB instance requires credentials to login;
- **password** is optional, and it will be ignored if *username* is not defined.

## Connect a container with this app to another container with MongoDB

```bash
# We are using bridged networks, which is the default, so we don't need to explicitly define it.

# Create the network
docker network create --label <network_name> --attachable

# If the container is already running, you can attach it to the network without restarting it.
docker network connect <network_name> <container_name>

# If the container is stopped, or has not been yet created, use this command.
docker run -d --network=<network_name> <container_name>
```
