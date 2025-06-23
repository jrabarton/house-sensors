# House temperature sensors
This was writen during the June 2025 heatwave in the UK, the idea was to experiemnet with building a simple website that shows the temperature trends for each room in my house.

## Overview
I had some DHT11 and ESP01 that I wanted to use, so thought that the best way to do this would be to have each of them send the data they read each minute into a RabbitMQ instance. This data is consumed by a microservice that writes it into a database. An API then allows this data to be extracted and used. Finally a front end displays all this.

## Left to do
* wire up the ESP8266's 
  * with DHT11
  * a power regulator
  * cr2032 battery holder
* deploy code for ESP8266's 
* convert it to run on a raspberry pi with docker

## Getting an example working
Use docker compose to build all the infrastructure
Run data-seed to populate some random data
