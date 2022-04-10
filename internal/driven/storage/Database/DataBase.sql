DROP DATABASE IF EXISTS hex_math;
CREATE DATABASE hex_math;
USE hex_math;

DROP TABLE IF EXISTS Questions;
CREATE TABLE `Questions` (
  `questionsId` int PRIMARY KEY,
  `Problem` varchar(255),
  `quantity` int DEFAULT 1
);

DROP TABLE IF EXISTS choices;
CREATE TABLE `Choices` (
  `QId` int,
  `choice` varchar(255),
  PRIMARY KEY (`QId`, `choice`)
);

DROP TABLE IF EXISTS Games;
CREATE TABLE `Games` (
  `GamesId` int,
  `PlayerName` int UNIQUE,
  `Scenario` varchar(255) DEFAULT "NEW_QUESTION",
  `score` int,
  `correct` int,
  `CurrentQuestion` int,
  `timeoute` int DEFAULT 5
);

ALTER TABLE `choices` ADD FOREIGN KEY (`QId`) REFERENCES `Questions` (`questionsId`);

ALTER TABLE `Games` ADD FOREIGN KEY (`CurrentQuestion`) REFERENCES `Questions` (`questionsId`);

