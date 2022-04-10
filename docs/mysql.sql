CREATE DATABASE hex_math;
USE hex_math;

DROP TABLE IF EXISTS questions;
CREATE TABLE `questions` (
  `id` int PRIMARY KEY,
  `Problem` text,
  `correctIndex` int
);

DROP TABLE IF EXISTS choices;
CREATE TABLE `choices` (
  `questionId` int,
  `choice` varchar(255),
  PRIMARY KEY (`questionId`, `choice`)
);

DROP TABLE IF EXISTS games;
CREATE TABLE `games` (
  `id` int,
  `playerName` int UNIQUE,
  `scenario` varchar(20) DEFAULT "NEW_QUESTION",
  `score` int,
  `countCorrect` int,
  `questionId` int,
  `questionTimeout` int DEFAULT 5
);

ALTER TABLE `choices` ADD FOREIGN KEY (`questionId`) REFERENCES `questions` (`id`);

ALTER TABLE `games` ADD FOREIGN KEY (`questionId`) REFERENCES `questions` (`iq`);

