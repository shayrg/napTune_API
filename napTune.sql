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
-- Songs
create table songs(
	id int primary key auto_increment,
	name VARCHAR(20) not NULL,
	artistId int not NULL,
	length VARCHAR(20) not null,
	location VARCHAR(40),
	foreign key (artistId) references people(id)
);
-- Playlist
CREATE TABLE playlists (
	id int PRIMARY KEY auto_increment,
	name varchar(20) not null
);
-- Playlist Songs
CREATE TABLE playlistSongs (
	playlistId int,
	songId int,
	position int,
	FOREIGN KEY (playlistId) REFERENCES playlists(id),
	FOREIGN KEY (songId) REFERENCES songs(id),
	PRIMARY KEY (playlistId, songId, position)
)


--Default data
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
insert into people (
	firstName,
	lastName,
	roll
) values (
	"Song",
	"Writer",
	"artist"
);
insert into songs (
	name,
	artistId,
	length,
	location
) VALUES (
	"test song",
	last_insert_id(),
	"03:12",
	"test_song.mp3"
);
insert into playlists (name) VALUES ("test playlist");
INSERT into playlistSongs(playlistId, songId, position) VALUES (1,1, 0);
INSERT into playlistSongs(playlistId, songId, position) VALUES (1,1, 1);