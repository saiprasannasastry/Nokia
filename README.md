# Implement Image Store Service using microservice architecture (programming language – Golang)

## DESCRIPTION
The design for API was taken using [jsonplaceholder.typicode.com/photos](jsonplaceholder.typicode.com/photos).
grpc and grpc gateway was used to solve the above problem
- Databse used postgres
- Message Queue used kafka
- container technology used docker
- The consumer container was just created to see catch the messages and show the logs

- To play around just run make build: This will build/ pull all the images and get the setup ready to run http calls

## TASKS
- Create/Delete Image Album using REST API
- Create/Delete Image in an Album using REST API
- Get an Image and all images in an Album using REST API (Persistence of your choice)
- Produce notification message to a topic whenever Create/Delete image using message broker like Kafka

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
  
  ## NOTE
  - Swaager UI is not working
  - for swagger.json check www folder
  - protocompile was used to generate swagger.json
