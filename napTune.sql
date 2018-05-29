drop database if exists napTune;
create database napTune;
use napTune;
-- People
create table people(
	id int primary key auto_increment,
	firstName varchar(20) not null,
	lastName varchar(20) not null,
	roll varchar(20) not null
);
-- Users
create table users(
	id int primary KEY,
	email varchar(20) not null,
	password varchar(20) not null,
	token varchar(50) not null DEFAULT "",
	expiration datetime,
	FOREIGN KEY (id) REFERENCES people(id)
);
insert into people (
	firstName,
	lastName,
	roll
) values (
	"root",
	"user",
	"admin"
);
insert into users(
	id,
	email,
	password,
	expiration
) VALUES (
	last_insert_id(),
	"root@localhost",
	"password",
	now()
);
-- Songs
-- create table songs (
-- 	id int primary key auto_increment
-- 	,name VARCHAR(20) not NULL
-- 	,artistId int not NULL
-- 	,length TIME not null
-- 	,location VARCHAR(20)
-- );
-- insert into songs (
-- 	name,
-- 	artistId,
-- 	length,
-- 	location
-- ) VALUES (
-- 		"test song"
-- 		,1
-- 		,00:03:00
-- 		,"/location/test_song.mp3"
-- )