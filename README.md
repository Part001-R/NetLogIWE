![](https://github.com/Part001-R/assets/blob/main/assets/netLogIWE.jpeg)

Pet project - collecting messages over the network and archiving them in a database. gRPCS is used.

Server recieve data in format:
```protobuf
message MessageRequest{
    string typeMessage = 1; // I, W, E, T(test)
    string nameProject = 2;
    string locationEvent = 3; 
    string bodyMessage = 4; 
}
``````

If the save is successful, it returns - Ok.
```protobuf
message MessageResponse{
    string status = 1;
}
``````

FaultForGRPC - a project that generates messages.

+ `v0.0.1` - Basic functionality.
+ `v0.0.2` - Fix. Working with the database through the interface.
+ `v0.0.3` - Added table overflow tracking and adding new ones.
+ `v0.0.4` - Added tests (Part-1) and fix.
+ `v0.0.5` - Added tests (Part-2 SUCCESS) and fix.
+ `v0.0.6` - Dockerfile, compose, workflows, fix.
