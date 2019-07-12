CREATE TABLE `student` (
  `id` bigint(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `stud_name` varchar(32) NOT NULL DEFAULT '' COMMENT '学生姓名',
  `stud_age` int(11) NOT NULL DEFAULT '0' COMMENT '学生年龄',
  `stud_sex` varchar(8) NOT NULL DEFAULT '' COMMENT '学生性别',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(11) NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_stud_name` (`stud_name`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4 COMMENT='学生表';

CREATE TABLE `teacher` (
  `id` bigint(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `teacher_name` varchar(32) NOT NULL DEFAULT '' COMMENT '老师姓名',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(11) NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT='老师表';