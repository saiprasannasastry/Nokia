# Implement Image Store Service using microservice architecture (programming language – Golang)

## TASKS
- Create/Delete Image Album using REST API
- Create/Delete Image in an Album using REST API
- Get an Image and all images in an Album using REST API (Persistence of your choice)
- Produce notification message to a topic whenever Create/Delete image using message broker like Kafka

## PREREQUISITES
- git , to clone the repo
- Docker
- Internet connection( to pull kafka, zookeeper, golang, postgres)

## DESCRIPTION
The design for API was taken using [jsonplaceholder.typicode.com/photos](jsonplaceholder.typicode.com/photos).
grpc and grpc gateway was used to solve the above problem
- Databse used postgres
- Message Queue used kafka
- container technology used docker
- The consumer container was just created to see catch the messages and show the logs
- Try the supported calls in any order and you should see appropriate error/success. - The supported calls with examples are given below
- The main code is in cmd/server folder
- consumer logs shows panic a couple of times until it is able to connect to kafka. But this wont affect any calls you do because it would have connected to kafka container by then. It just takes a fraction of a second
- make sure port 8080, 5432, 8010,9092,5566 are not in use in your local system. 
  - if you wish to clean up the ports run `lsof -i tcp:$PORT` and then kill the port using PID `kill $PID` . Use this command if your bringing up the environment on mac
  - for ubunutu  run `netstat -anp|grep "$PORT"` then kill the port using PID `kill $PID` 

## To play around just clone this repo `git clone https://github.com/saiprasannasastry/Nokia.git` run `make build` : This will build/ pull all the images and get the setup ready to run http calls

  ## NOTE
  - except albumID everything has to be unique otherwise you will get "ALREADY_EXISTS" error. This makes sense because the photo id and title always has to be uniqe as well the URL
  - You can see error logs in server using `docker logs -f photo`
  - The container names are as follows
    - `PHOTO` main container
    - `postgres-db` DB container
    - `consumer` Kafka consumer
    - `kafka` kafka container
    - `zookerper` zooker container
   - logs of each container can be accessed using `docker logs -f $container name`
  - To bring down the container run `make clean`
    - This brings down the containers and cleans all the volumes and images
    - If you wish to keep the volumes and images for next trail , navigate to docker folder and run `docker-compose down`
  - Swaager UI is not working
  - for swagger.json check www folder
  - protocompile was used to generate swagger.json

## SUPORTED HTTP CALLS
- curl -X POST "http://localhost:8080/album" -d '{"id":2,"albumId":2,"title":"accusamus beatae ad facilis cum similique qui sunt",
"url":"https://via.placeholder.com/600/92c952","thumbNailUrl":"https://via.placeholder.com/150/92c952"}'
  - The above creates and stores the value in database

- Curl -X GET “”http://localhost:8080/getalbums”
  - The above calls gets all the albums from the database

- Curl -X GET "http://localhost:8080/getalbums/$(ALBUMID)"
  - The above call gets all the albums that corresponds to album id. Replace the $(ALBUMID) with the ID you want to retrive
  
- Curl -X GET "http://localhost:8080/getalbums/$(ALBUMID)/photo/$(PHOTOID)"
  - The above call gets a particular photo from the corresponding Album id. Replace the variables with appropriate values
     
- Curl -X PUT "http://localhost:8080/Updatealbum" -d '{"oldAlbumId":2,"newAlbumId":4,"oldTitle":"bleh", "newTitle":"sai1"}'
  - The above method was implemented keeping in mind if the the photo had to be moved from x to y
  
- curl -X DELETE http://localhost:8080/photo/($ID)
  - The photo you want to delete directly, the album id is not requred
  
  

