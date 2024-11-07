CREATE TABLE `users` (
  `id` varchar(255) PRIMARY KEY,
  `name` varchar(100) NOT NULL,
  `email` varchar(100) UNIQUE NOT NULL,
  `password` varchar(155) NOT NULL,
  `photo` varchar(255),
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `posts` (
  `id` varchar(255) PRIMARY KEY,
  `content` text,
  `user_id` varchar(255),
  `image` varchar(255),
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `comments` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `user_id` varchar(255),
  `post_id` varchar(255),
  `comment` varchar(255) NOT NULL,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `votes` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `user_id` varchar(255) UNIQUE,
  `post_id` varchar(255) UNIQUE,
  `vote` tinyint NOT NULL,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE `posts` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `comments` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `comments` ADD FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`);

ALTER TABLE `votes` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `votes` ADD FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`);
