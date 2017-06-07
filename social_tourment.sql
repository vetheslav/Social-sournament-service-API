# ************************************************************
# Sequel Pro SQL dump
# Версия 4541
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Адрес: 127.0.0.1 (MySQL 5.5.5-10.2.6-MariaDB)
# Схема: social_tourment
# Время создания: 2017-06-07 13:38:38 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Дамп таблицы players
# ------------------------------------------------------------

DROP TABLE IF EXISTS `players`;

CREATE TABLE `players` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `balance` decimal(10,2) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `players` WRITE;
/*!40000 ALTER TABLE `players` DISABLE KEYS */;

INSERT INTO `players` (`id`, `balance`)
VALUES
	(1,2010.00),
	(2,20.00),
	(3,30.00),
	(4,40.00);

/*!40000 ALTER TABLE `players` ENABLE KEYS */;
UNLOCK TABLES;


# Дамп таблицы tournament_player_backers
# ------------------------------------------------------------

DROP TABLE IF EXISTS `tournament_player_backers`;

CREATE TABLE `tournament_player_backers` (
  `tournament_id` int(11) DEFAULT NULL,
  `player_id` int(11) DEFAULT NULL,
  `backer_id` int(11) DEFAULT NULL,
  UNIQUE KEY `tor_play` (`tournament_id`,`player_id`,`backer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `tournament_player_backers` WRITE;
/*!40000 ALTER TABLE `tournament_player_backers` DISABLE KEYS */;

INSERT INTO `tournament_player_backers` (`tournament_id`, `player_id`, `backer_id`)
VALUES
	(1,2,3),
	(1,2,4);

/*!40000 ALTER TABLE `tournament_player_backers` ENABLE KEYS */;
UNLOCK TABLES;


# Дамп таблицы tournament_players
# ------------------------------------------------------------

DROP TABLE IF EXISTS `tournament_players`;

CREATE TABLE `tournament_players` (
  `tournament_id` int(11) NOT NULL,
  `player_id` int(11) NOT NULL,
  PRIMARY KEY (`tournament_id`,`player_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `tournament_players` WRITE;
/*!40000 ALTER TABLE `tournament_players` DISABLE KEYS */;

INSERT INTO `tournament_players` (`tournament_id`, `player_id`)
VALUES
	(1,1),
	(1,2);

/*!40000 ALTER TABLE `tournament_players` ENABLE KEYS */;
UNLOCK TABLES;


# Дамп таблицы tournaments
# ------------------------------------------------------------

DROP TABLE IF EXISTS `tournaments`;

CREATE TABLE `tournaments` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `deposit` decimal(15,2) DEFAULT NULL,
  `status` tinyint(4) DEFAULT 1,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `tournaments` WRITE;
/*!40000 ALTER TABLE `tournaments` DISABLE KEYS */;

INSERT INTO `tournaments` (`id`, `deposit`, `status`)
VALUES
	(1,1000.00,0),
	(2,1000.00,0),
	(3,1000.00,0),
	(4,1000.00,0),
	(5,1000.00,0),
	(6,1000.00,0),
	(7,1000.00,0),
	(8,1000.00,0),
	(9,1000.00,0);

/*!40000 ALTER TABLE `tournaments` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
