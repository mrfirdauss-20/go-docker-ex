CREATE DATABASE hex_math;
USE hex_math;

CREATE TABLE `questions` (
  `id` int PRIMARY KEY,
  `problem` text,
  `correct_index` int
);

CREATE TABLE `choices` (
  `question_id` int,
  `choice` varchar(255),
  PRIMARY KEY (`question_id`, `choice`)
);

CREATE TABLE `games` (
  `id` varchar(36)  PRIMARY KEY,
  `player_name` int,
  `scenario` varchar(20) DEFAULT "NEW_QUESTION",
  `score` int,
  `count_correct` int,
  `question_id` int,
  `question_timeout` int DEFAULT 5
);

ALTER TABLE `choices` ADD FOREIGN KEY (`question_id`) REFERENCES `questions` (`id`);

ALTER TABLE `games` ADD FOREIGN KEY (`question_id`) REFERENCES `questions` (`id`);

