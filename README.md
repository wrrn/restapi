Configuration REST API
======================

## Installation
1. Install postgres
2. Add user ```tenable``` with password ```insecure```
3. Create database ```restapi```
3. Install the server (this will install binary in your $GOPATH/bin directory)
```go get -u github.com/warrenharper/restapi```
4. Run the sql script located in github.com/warrenharper/restapi/env/create_db.sql
```psql -U tenable -d restapi -a -f env/create_db.sql -h localhost```

## Running
Assuming $GOPATH/bin is in your path you can just run ```restapi```. This will start the server on port 8080.

## Running the tests
The tests assume that you have a database with the name ```testapi``` the has an identical schema to that of the database ```rest api```

## Vagrant
If you are familiar with vagrant you can cd into the root directory and run ```vagrant up``` and all of the enviroment will be setup. You will need to cross compile the binary if you are not running vagrant on a linux machine.

``` bash
GOOS=linux GOARCH=amd64 go build
```

#API
## Authentication
### Log in

```
POST /login
```

__Input__

| parameter| Description |
|-----------|------------|
|"username"| __Required__: The username as a string |
|"password"| __Required__: The password as a string |

__Example__
``` js
{
	"name": "john_doe",
	"password": "password"
}
```

__Response__

| Status | Body |
| ---- | ---- |
| 200 | "Authorized"|
| 401 | "Unauthorized" |

### Log out

``` bash
POST /logout
```
__Input__
This requires no input

__Response__

| Status | Body |
| ---- | ---- |
| 200 | "Success"|


#### __Note:__ If you are not authenticated you will receive a status code of 403 when you try to access anything


## Configuration
### List configurations

``` bash
GET /configurations/
```
__NOTE:__ Note the url ends in a forward slash ("/")

__Input__
This requires no input.

__Response__

| Status |
|:------:|
| 200    |

``` js
{
 "configurations": [
  {
   "id": 1,
   "name": "Config2",
   "hostname": "add.here",
   "port": 3384,
   "username": "warren"
  },
  {
   "id": 2,
   "name": "A",
   "hostname": "b.good",
   "port": 3,
   "username": "abernathy"
  }
 ]
}
```

### Get individual configuration
Get the configuration with the matching name

``` bash
GET /configurations/:name
```

__Input__
This requires no input.

__Response__

| Status |      Body     |            Description           |
|:------:| :-----------: | :------------------------------: |
| 200    | _See example_ | Found the configuration          |
| 404    |               | Could not find the configuration |

__Example__

``` bash
GET /configurations/Config2
```
``` js
{
 "configurations": [
  {
   "id": 6,
   "name": "Config2",
   "hostname": "add.here",
   "port": 3384,
   "username": "warren"
  },
 ]
}
```

### Add configuration
Add a configuration to the list of configurations

``` bash
POST /configurations/
```

__Input__

| parameter| Description | Type |
|-----------|------------| ---- |
|"name"| __Required__: The name of the configuation | string |
|"hostname"| __Required__: The hostname | string | 
|"port"| __Required__: The port  | int | 
|"username"| __Required__: The username for the configuration |string |

__Example__
``` js
 {
   "name": "Config2",
   "hostname": "add.here",
   "port": 3384,
   "username": "warren"
}
```

__Response__

| Status |      Body     |            Description           |
|:------:| :-----------: | :------------------------------: |
| 200    | _See example_ | Configuration was added          |
| 409    |  List of configurations that collide with the name of the addee | Name collision |


__Example__
``` js
{
 "configurations": [
  {
   "id": 6,
   "name": "Config2",
   "hostname": "add.here",
   "port": 3384,
   "username": "warren"
  },
 ]
}
```

### Delete an individual configuration
Delete the configuration with the matching name

``` bash
DELETE /configurations/:name
```

__Input__
This requires no input.

__Response__

| Status |      Body     |            Description           |
|:------:| :-----------: | :------------------------------: |
| 209    |   | Found the configuration          |
| 404    |               | Could not find the configuration |

__Example__

``` bash
DELETE /configurations/Config2
```


### Modify an individaul configuration

``` bash
PATCH /configurations/:name
```

__Input__

| parameter| Description | Type |
|-----------|------------| ---- |
|"name"| The name of the configuation | string |
|"hostname"| The hostname | string | 
|"port"| The port  | int | 
|"username"| The username for the configuration |string |

__Note:__ Any of the input fields that are ommited will remain the same



__Response__

| Status |      Body     |            Description           |
|:------:| :-----------: | :------------------------------: |
| 200    | _See example_ | Configuration was added          |
| 409    |  List of configurations that collide with the name of the modified | Name collision |


__Example__
``` bash
PATCH /configurations/Config2
```

```
{
 "name": "Config65"
}
```

_Response_
```js
{
 "configurations": [
  {
   "id": 6,
   "name": "Config65",
   "hostname": "add.here",
   "port": 3384,
   "username": "warren"
  },
 ]
}
```

## Sorting and Pagination
Sort your configurations and retrieve them by page

### Sorting
Add the ```sort``` parameter to the ```GET /configurations/``` request with one of the following values:

| Value | Description |
| :--:  | :---------: |
| name  | Sort configurations by name |
| hostname | Sort configurations by hostname |
| port | Sort configurations by port |
| username | Sort configurations by username |

__Example__

``` bash
GET /configurations/?sort=name
```

### Pagination
To paginate you must have both parameters ```page``` and ```per_page```.
```page``` is 0 based.
```per_page``` cannot be greater than 100.

__Example__

``` bash
GET /configurations/?page=0&per_page=50
```
