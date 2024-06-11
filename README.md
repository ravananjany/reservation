# reservation
#protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative     protos/blog/blog.proto


# get user and seat deatils based on requested section
curl --location --request GET 'http://localhost:8080/tickets/a' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id":"a9b7c24d-ef74-4417-a8ee-b7d243c6266f",
   "first_name":"priya2113",
   "last_name":"jany2113",
   "email_id":"abc21@gmail.com"
}'


# get receipt details 
curl --location --request GET 'http://localhost:8080/ticket/24099ed7-1be1-45d4-ab2e-81413b2c7f9f' \
--header 'Content-Type: application/json' \
--data-raw '{
    
   "first_name":"priya21",
   "last_name":"jany21",
   "email_id":"abc21@gmail.com"
}'

# purchase a ticket 
curl --location 'http://localhost:8080/ticket' \
--header 'Content-Type: application/json' \
--data-raw '{
    
   "first_name":"priya21",
   "last_name":"jany21",
   "email_id":"abc21@gmail.com"
}'

# Delete a user ticket
curl --location --request DELETE 'http://localhost:8080/ticket/19b57c2c-51b8-4fd6-85c1-95b53d1091a5' \
--header 'Content-Type: application/json' \
--data-raw '{
    
   "first_name":"priya21",
   "last_name":"jany21",
   "email_id":"abc21@gmail.com"
}'

# modify a user ticket 

curl --location --request PUT 'http://localhost:8080/ticket' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id": "310037b7-f384-4597-af1e-3bd1c421a138",
   "first_name":"priya215444444",
   "last_name":"jany21",
   "email_id":"abc21@gmail.com"
}'

# unit tests coverage for client is above 90 percent
# unit tests coverage for server is above 80 percent