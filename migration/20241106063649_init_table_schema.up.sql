CREATE TABLE `users` (
  `id` char(36) PRIMARY KEY,
  `name` varchar(100) NOT NULL,
  `email` varchar(100) UNIQUE NOT NULL,
  `password` varchar(155) NOT NULL,
  `photo` varchar(255),
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `posts` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `content` text,
  `user_id` char(36),
  `image` varchar(255),
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `comments` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `user_id` char(36),
  `post_id` int,
  `comment` varchar(255) NOT NULL,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `votes` (
  `user_id` char(36),
  `post_id` int,
  `vote` tinyint NOT NULL,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`, `post_id`)
);

ALTER TABLE `posts` 
ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `votes` 
ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
ADD FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`);

ALTER TABLE `comments` 
ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
ADD FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`);