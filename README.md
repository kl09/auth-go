Golang service for auth, based on Clean Architecture:

Register: 
```
curl -v -X POST http://localhost:8080/v1/register -d '{"email":"example@example.org","password":"12345"}' -H "content-type: application/json"
```

Auth:
```
curl -v -X POST http://localhost:8080/v1/auth -d '{"email":"example@example.org","password":"12345"}' -H "content-type: application/json"
```

Get by token:
```
curl -v -X GET http://localhost:8080/v1/users-by-token/2GdxFOD8YLyXmiI1-I2265SKo1SaQBq3AM1AQUZQcAHkty3yBS4-Yyi7HLtD4fAN4vuniK74sphFCBqQmkuE12Ucmv3dYxmwYFgCUoA7VkROMDzWUngrU7xcQG1pCLUw 
```

