CREATE DATABASE hex_math;
USE hex_math;

CREATE TABLE `questions` (
  `id` int PRIMARY KEY,
  `problem` text NOT NULL,
  `correct_index` int NOT NULL
);

CREATE TABLE `choices` (
  `question_id` int,
  `choice` varchar(255),
  PRIMARY KEY (`question_id`, `choice`)
);

CREATE TABLE `games` (
  `id` varchar(36)  PRIMARY KEY,
  `player_name` int NOT NULL,
  `scenario` varchar(20) NOT NULL DEFAULT "NEW_QUESTION",
  `score` int NOT NULL,
  `count_correct` int NOT NULL,
  `question_id` int NOT NULL,
  `question_timeout` int NOT NULL DEFAULT 5
);

ALTER TABLE `choices` ADD FOREIGN KEY (`question_id`) REFERENCES `questions` (`id`);

ALTER TABLE `games` ADD FOREIGN KEY (`question_id`) REFERENCES `questions` (`id`);

