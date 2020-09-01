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
## Run this container together with a MongoDB container (and Swagger-UI)

Before you start, make sure you have the correct hostname defined in `swagger/doc.go` :
```bash
//     Schemes: http
//     BasePath: /
//     Version: 1.0.0
//     Host: 192.168.0.109:8080
//
//     Contact: Otavio Krambeck <rofatto@gmail.com> https://otaviokr.github.io
//
```

The script `startup.sh` will generate the JSON file and then run the `docker-compose.yml` to start all services.

```bash
./startup.sh
```

## Features

These are the functionalities available

### Insert (/insert/\<db\>/\<collection\>)

You must send via POST a JSON object that represents a document. That document will be inserted in **collection**, in **database**.

It returns the ObjectID of the newly inserted document.

### Find (/find/\<db\>/\<collection\>)

If you provide a JSON object with filters, it will return all documents from **collection** in **database** that match the filter provided. If no filter is given, all documents in **collection** will be returned.

The method must be POST, and it returns the documents as a JSON object.

### Update (/update/\<db\>/\<collection\>)

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

## Collections expected

The reasoning behind having a definitive definition of the collections are explained on my tech-log: <here>.

TL;DR: the goal here is to get experience with tools and frameworks, and having the collections well-defined is very helpful to integrate Swagger.

### Quote

| Field           | Description                                  | Type   |
| --------------- | -------------------------------------------- | ------ |
| Publications    | How many times has this entry been published | int    |
| LastPublished   | Timestamp (epoch) when it was last published | int64  |
| OriginalQuote   | The quote, in the original language          | string |
| TranslatedQuote | The quote, translated                        | string |
| Author          | Name of the author                           | string |
