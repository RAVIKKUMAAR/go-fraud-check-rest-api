# go-fraud-check-rest-api

## General Info
This project provides the RESTful interface for the fraud-check service that ultimately calls TensorFlow to perform the scoring for a purchase request

## Build

Edit the Dockerfile and change the FROM base image to match the platform you are building on.

FROM golang                  -- For x86 builds
FROM icr.io/ibmz/golang:1.15 -- for zCX builds

run 'docker build -t go-docker .' to build the docker image.

## Running

run 'docker run -d -p 8080:8080 go-docker'

Will listen on port 8080. Change the port if you have a conflict in your environment. 

## Backend REST API Server LoZ 
[Link](https://github.ibm.com/ai-on-z/go-fraud-check-rest-api/tree/master/loz)
