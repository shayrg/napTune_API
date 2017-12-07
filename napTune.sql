drop database if exists napTune;
create database napTune;
use napTune;
create table users(
	id int primary key auto_increment
	,firstName varchar(20) not null
	,lastName varchar(20) not null
	,email varchar(20) not null
	,password varchar(20) not null
	,token varchar(50)
	,expiration datetime
	,role varchar(20) not null
);
insert into users (
	firstName
	,lastName
	,email
	,password
	,role
  ,expiration
) values (
	"root"
	,"user"
	,"root@localhost"
	,"password"
	,"admin"
  ,NOW()
)