#Gaming website
This is a RESTful application-game in which each player holds certain amount of bonus points,
which could be deposited or taken. Player's balance could only be >= 0.
Players could join the tournament. When they join - they put a deposit from their account 
to the tournament's total prize.
Deposit could only be >= 0.

To start the server build & run the [server.go](./server/server.go) file. The server listens
on port :8080 

##Actions
|Command & URI         |Action                             |
|:---------------------|:----------------------------------|
|`POST` /user          |Registers a new user               |
|`GET` /user/{id}      |Gets a user's info                 |
|`DELETE` /user/{id}   |Removes a user                     |
|`POST` /user/{id}/take|Takes 300 points from users account|
|`POST` /user/{id}/fund|Adds 400 points from user's account|
---
`POST` /user  
**Request**  
  
{  
    "name" :  name,  
    "balance": 1000  
} 
   
**Response**  
  
{  
    "id": 1,  
    "name" :  name,  
    "balance": 700  
}  
---
`GET` /user/{id}  
**Response**  
  
{  
    "id": 1,  
    "name" :  name,  
    "balance": 700  
}  
---
`DELETE` /user/{id}  
**Response**    
{}  
---
`POST` /user/{id}/take  
**Request** 
   
{  
    "points" :  300  
}  

**Response**  
(initial user's balance is 1000)  
  
{  
    "id": 1,  
    "name" :  name,  
    "balance": 700  
}  
---
`POST` /user/{id}/fund  
**Request**  
  
{  
    "points" :  400  
}  
  
**Response**
(initial user's balance is 700) 
      
{  
    "id": 1,  
    "name" :  name,  
    "balance": 1100  
}  
---
