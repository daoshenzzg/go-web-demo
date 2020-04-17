# ************************************************************
# Sequel Pro SQL dump
# Version 4529
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 5.7.28)
# Database: school
# Generation Time: 2020-04-17 07:08:57 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table student
# ------------------------------------------------------------

DROP TABLE IF EXISTS `student`;

CREATE TABLE `student` (
                           `id` bigint(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                           `stud_name` varchar(32) NOT NULL DEFAULT '' COMMENT '学生姓名',
                           `stud_age` int(11) NOT NULL DEFAULT '0' COMMENT '学生年龄',
                           `stud_sex` varchar(8) NOT NULL DEFAULT '' COMMENT '学生性别',
                           `create_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '创建时间',
                           `update_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '修改时间',
                           PRIMARY KEY (`id`),
                           UNIQUE KEY `uniq_stud_name` (`stud_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='学生表';

LOCK TABLES `student` WRITE;
/*!40000 ALTER TABLE `student` DISABLE KEYS */;

INSERT INTO `student` (`id`, `stud_name`, `stud_age`, `stud_sex`, `create_time`, `update_time`)
VALUES
(1,'张三',18,'男','2020-04-17 11:17:12','2020-04-17 11:17:12'),
(2,'李四',19,'男','2020-04-17 11:17:12','2020-04-17 11:17:12'),
(3,'如花',16,'女','2020-04-17 11:17:12','2020-04-17 11:17:12');

/*!40000 ALTER TABLE `student` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table teacher
# ------------------------------------------------------------

DROP TABLE IF EXISTS `teacher`;

CREATE TABLE `teacher` (
                           `id` bigint(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                           `teacher_name` varchar(32) NOT NULL DEFAULT '' COMMENT '老师姓名',
                           `create_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '创建时间',
                           `update_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '修改时间',
                           PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='老师表';

LOCK TABLES `teacher` WRITE;
/*!40000 ALTER TABLE `teacher` DISABLE KEYS */;

INSERT INTO `teacher` (`id`, `teacher_name`, `create_time`, `update_time`)
VALUES
(1,'苍老师','2020-04-17 11:17:12','2020-04-17 11:17:12');

/*!40000 ALTER TABLE `teacher` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
