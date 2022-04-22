CREATE DATABASE hex_math;
USE hex_math;

CREATE TABLE `questions` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `problem` text NOT NULL,
  `correct_index` int NOT NULL,
  `answers` json NOT NULL
);

CREATE TABLE `games` (
  `id` varchar(36)  PRIMARY KEY,
  `player_name` TEXT NOT NULL,
  `scenario` varchar(20) NOT NULL,
  `score` int NOT NULL,
  `count_correct` int NOT NULL,
  `question_id` int NOT NULL,
  `question_timeout` int NOT NULL
);

ALTER TABLE `games` ADD FOREIGN KEY (`question_id`) REFERENCES `questions` (`id`);

INSERT INTO questions (problem, correct_index, answers) VALUES
('1 + 1', 2, '[1, 2, 3]'),
('1 + 2', 3, '[1, 2, 3]'),
('2 + 1', 3, '[1, 2, 3]'),
('3 - 2', 1, '[1, 2, 3]'),
('2 - 1', 1, '[1, 2, 3]'),
('2 + 2 - 1', 3, '[1, 2, 3]'),
('1 + 1 + 1', 3, '[1, 2, 3]'),
('3 - 1', 2, '[1, 2, 3]'),
('2 + 1 - 2', 1, '[1, 2, 3]'),
('1 + 1 - 1', 1, '[1, 2, 3]');