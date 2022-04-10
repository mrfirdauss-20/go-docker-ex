CREATE TABLE `Questions` (
  `questionsId` int PRIMARY KEY,
  `Problem` varchar(255),
  `quantity` int DEFAULT 1
);

CREATE TABLE `choices` (
  `QId` int,
  `choice` varchar(255),
  PRIMARY KEY (`QId`, `choice`)
);

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

