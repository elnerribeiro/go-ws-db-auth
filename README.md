# go-ws-db-auth

This is a Proof of Concept of using Golang in a somewhat standard commercial program.
Here we have some Rest web services with Authentication/Authorization using JWT, database connection with Postgres, a persistence layer that uses structs to abstract table models and Mustache to do query manipulation and some examples of sync/async calls.

AS THIS IS JUST A POC, NOT EVERYTHING IS DONE EXACTLY THE WAY IT SHOULD:

	1. Microservices: The system should have only one endpoint for each purpose. So, we should have an endpoint for Auth, other for users and other for insert tests. Each should run on a separate Go program and proxied by an Apache/Nginx server.
	2. Better modules: To maximize code reusability, when creating microservices, we should have each layer on its own modules, so Auth could import user, for example.
	3. Methods/Interfaces: Functions should be used only when there's no side effects. For calls with side effects, an interface should be created and its methods implemented.
	4. Tests and Documentation: Yeah... I know...
	5. Logging: Will be done soon (probably tomorrow)

Create database:

	docker pull postgres
	mkdir pgdata
	docker volume create -d local --name pgdata
	docker run -p5432:5432 -e POSTGRES_PASSWORD=root -v pgdata:/var/lib/postgresql/data postgres
	docker ps (get container id)
	docker exec -it <container_id> bash
	psql -Upostgres
	create database teste;
	\connect teste;
	create table usuario (id serial not null, email varchar(50) not null, role varchar(20) not null, password varchar(128) not null, primary key (id));
	create table ins_id (id serial not null, type varchar(20) not null, quantity int not null, status varchar(20) not null, tstampinit bigint, tstampend bigint, primary key (id));
	create table insert_batch(id serial not null, id_ins_id int not null, pos int not null, primary key(id), foreign key (id_ins_id) references ins_id(id));
	--PWD abc
	insert into usuario (email, role, password) values ('usuario@usuario.com', 'user', 'DDAF35A193617ABACC417349AE20413112E6FA4E89A97EA20A9EEEE64B55D39A2192992A274FC1A836BA3C23A3FEEBBD454D4423643CE80E2A9AC94FA54CA49F');
	--PWD 123
	insert into usuario (email, role, password) values ('admin@admin.com', 'admin', '3C9909AFEC25354D551DAE21590BB26E38D53F2173B8D3DC3EEE4C047E7AB1C1EB8B85103E3BE7BA613B31BB5C9C36214DC9F14A42FD7A2FDB84856BCA5C44C2');

Build and run:

	(Linux)
	go build -o main
	./main

	(Windows)
	go build -o main.exe
	.\main.exe

Default URL: 

	http://localhost:8000

Postman Collection:

	https://www.getpostman.com/collections/06704a4c68b44e63502e

Calls:

	/api/login (POST)
		Request:
			Headers:
				Content-Type: application/json
				Authorization: Basic SldUcGFzc3dvcmQxMjNA
			Body:
				{
					"email":"usuario@usuario.com","password":"DDAF35A193617ABACC417349AE20413112E6FA4E89A97EA20A9EEEE64B55D39A2192992A274FC1A836BA3C23A3FEEBBD454D4423643CE80E2A9AC94FA54CA49F"
				}
		Response:
			{
				"account": {
					"id": 1,
					"email": "usuario@usuario.com",
					"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjEsIlJvbGUiOiJ1c2VyIiwiZXhwIjoxNTg1MTI0ODIwfQ.d5GsE9WDWlbRxQRtfuAO-G0SFnYV1ZhAs-m5rb1t--E",
					"role": "user"
				},
				"message": "Logged In",
				"status": true
			}

	/api/validate (GET)
		Request:
			Headers:
				Content-Type: application/json
				Authorization: Bearer {{token}}
		Response:
			{
				"message": "success",
				"role": "user",
				"status": true,
				"userId": 1
			}
	
	/api/users (POST) - Only role admin
		Request:
			Headers:
				Content-Type: application/json
				Authorization: Bearer {{token}}
		Response:
			{
				"message": "Usuário sem permissão",
				"status": false
			}
			ou
			{
				"data": [
					{
						"id": 1,
						"email": "usuario@usuario.com",
						"role": "user"
					},
					{
						"id": 2,
						"email": "admin@admin.com",
						"role": "admin"
					}
				],
				"message": "success",
				"status": true
			}

	/api/user/{id} (GET) - Only role admin
		Request:
			Headers:
				Content-Type: application/json
				Authorization: Bearer {{token}}
		Response:
			{
				"data": {
					"id": 2,
					"email": "admin@admin.com",
					"role": "admin"
				},
				"message": "success",
				"status": true
			}

	/api/user (PUT) - Only role admin
		Request:
			Headers:
				Content-Type: application/json
				Authorization: Bearer {{token}}
			Body:
				{
					"email":"elner.ribeiro@gmail.comx",
					"id": 9, //optional, only used for updates
					"role": "admin",
					"password":"3C9909AFEC25354D551DAE21590BB26E38D53F2173B8D3DC3EEE4C047E7AB1C1EB8B85103E3BE7BA613B31BB5C9C36214DC9F14A42FD7A2FDB84856BCA5C44C2"
				}
		Response:
			{
				"data": {
					"id": 9,
					"email": "elner.ribeiro@gmail.comx",
					"role": "admin"
				},
				"message": "success",
				"status": true
			}

	/api/user/{id} (DELETE) - Only role admin
		Request:
			Headers:
				Content-Type: application/json
				Authorization: Bearer {{token}}
		Response:
			{
				"message": "success",
				"status": true
			}

	/api/insert (DELETE)
		Request:
			Headers:
				Content-Type: application/json
				Authorization: Bearer {{token}}
		Response:
			{
				"message": "success",
				"status": true
			}

	/api/insert/{id} (GET)
		Request:
			Headers:
				Content-Type: application/json
				Authorization: Bearer {{token}}
		Response:
			{
				"data": {
					"id": 6,
					"type": "sync",
					"quantity": 10,
					"status": "Finished",
					"list": [
						{
							"id": 120107,
							"id_ins_id": 6,
							"pos": 1
						},
						...
					]
				},
				"message": "success",
				"status": true
			}

	/api/insert/sync/{quantity} (PUT)
		Request:
			Headers:
				Content-Type: application/json
				Authorization: Bearer {{token}}
		Response:
			{
				"data": {
					"id": 6,
					"type": "sync",
					"quantity": 10,
					"status": "Finished",
					"tstampinit": 1585175306,
        			"tstampend": 1585175306,
					"list": [
						{
							"id": 120107,
							"id_ins_id": 6,
							"pos": 1
						},
						{
							"id": 120108,
							"id_ins_id": 6,
							"pos": 2
						},
						{
							"id": 120109,
							"id_ins_id": 6,
							"pos": 3
						},
						{
							"id": 120110,
							"id_ins_id": 6,
							"pos": 4
						},
						{
							"id": 120111,
							"id_ins_id": 6,
							"pos": 5
						},
						{
							"id": 120112,
							"id_ins_id": 6,
							"pos": 6
						},
						{
							"id": 120113,
							"id_ins_id": 6,
							"pos": 7
						},
						{
							"id": 120114,
							"id_ins_id": 6,
							"pos": 8
						},
						{
							"id": 120115,
							"id_ins_id": 6,
							"pos": 9
						}
					]
				},
				"message": "success",
				"status": true
			}

	/api/insert/async/{quantity} (PUT)
		Request:
			Headers:
				Content-Type: application/json
				Authorization: Bearer {{token}}
		Response:
			{
				"data": {
					"id": 5,
					"type": "async",
					"quantity": 20000,
					"status": "Running",
					"tstampinit": 1585174744
				},
				"message": "success",
				"status": true
			}
