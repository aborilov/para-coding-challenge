### Coding Challenge

## Install

* clone repository
* run `docker-compose up`

that is all!

There are three endpoints:
* /api - just for testing(with authorization)
* /auth - account creation(without authorization)
* /login - authentication(without authorization)

Example:

* `curl http://127.0.0.1:8000/api`                       
  
  ```
  {"message":"Unauthorized"}
  ```
* `curl -d "email_address=aborilov@gmail.com&password=azsxdcfv" 127.0.0.1:8000/auth`

```
{
    "data":
        {
            "account_id":"3d3e193c-e475-4dc6-be6d-274b69efd490",
            "credentials": 
            {
              "jwt":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFib3JpbG92MUBnbWFpbC5jb20iLCJpc3MiOiJVTzdCNVBMOTRnczBrdGJLbEs4U0M4SllPYVFvaHhRQSJ9.Z2l8qyqqVDPqOVxK7lKW3OeWkAsZqtiKjEUtkgnsuGQ"
            }
        }
}
```
* `curl -d "email_address=aborilov@gmail.com&password=mama" 127.0.0.1:8000/login`

```
{
    "data":
    {
        "account_id":"3d3e193c-e475-4dc6-be6d-274b69efd490",
        "credentials":
        {
            "jwt":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFib3JpbG92MUBnbWFpbC5jb20iLCJpc3MiOiJVTzdCNVBMOTRnczBrdGJLbEs4U0M4SllPYVFvaHhRQSJ9.Z2l8qyqqVDPqOVxK7lKW3OeWkAsZqtiKjEUtkgnsuGQ"
        }
    }
}
```

* `curl 127.0.0.1:8000/api -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFib3JpbG92MUBnbWFpbC5jb20iLCJpc3MiOiJVTzdCNVBMOTRnczBrdGJLbEs4U0M4SllPYVFvaHhRQSJ9.Z2l8qyqqVDPqOVxK7lKW3OeWkAsZqtiKjEUtkgnsuGQ'`

```
Hello, World!
```
