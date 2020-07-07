# Description

This is the service that powers nickolinger.com/blog's upvotes. Currently, very 0.0.1.  

Uses go1.14

# Endpoints

POST `/users` - Returns a new user ID for front-end to store in localstorage.  
GET `/users/{userId}/posts/{name}` - Retrieves likes for a user for a particular post  
POST `/users/{userId}/likes/{id}` - Like a post  

# TODO'S
- [ ] Finish netlify webhook glue
- [ ] Tests
    - [ ] Write mock firestore service
    - [ ] Write integration tests for public / private key auth
    