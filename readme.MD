# This is a try to build distributaed application in GO 
  and it is based on the very good course : 
  [Building Distributed Applications with Go](https://app.pluralsight.com/library/courses/go-build-distributed-applications)


# About the application 

The application has 2 main component : sensors and discovery

##discovery has 2 main responsibility 
  - It will regularly check for new sensors and register them 
  - It will read the sensors data and print them to the console 
## Sensor 
  - Sensor will send reading data regularly 
  - Sensor will respond to discovery requests 


# To Access rabit admin page  : 

[rabit admin page](http://127.0.0.1:15672/#/)

**user** : rabbitmq
**password** : rabbitmq


# To run the application 

Make sure that golang is installed on your machine

* ```docker-compose up -d``` \
  right now, only rabbitmq is dockerized


* in the root directory of the project run 

  ```go run discovery/cmd/main.go```\
  and this will run the discovery component


* in the root directory of the project run in new terminals   

    
    ```go run sensor/cmd/main.go -name sensor1``` \
    ```go run sensor/cmd/main.go -name sensor2``` \
    ```go run sensor/cmd/main.go -name sensor3``` \
    ```go run sensor/cmd/main.go -name sensor4 -freq 5``` 


Those command will run and register 4 sensors to the system

At any time you can add / stop sensor , and the discovery component will take of it 


You may check the queue tab in the rabbitmq admin page to see what is going on

