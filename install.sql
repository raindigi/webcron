CREATE TABLE `t_task` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT 'USER ID',
  `group_id` int(11) NOT NULL DEFAULT '0' COMMENT 'GROUP ID',
  `task_name` varchar(50) NOT NULL DEFAULT '' COMMENT 'TASK NAME',
  `task_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT 'TASK TYPE',
  `description` varchar(200) NOT NULL DEFAULT '' COMMENT 'DESCRIPTION',
  `cron_spec` varchar(100) NOT NULL DEFAULT '' COMMENT 'CRON SPEC',
  `concurrent` tinyint(4) NOT NULL DEFAULT '0' COMMENT 'CONCURRENT',
  `command` text NOT NULL COMMENT 'COMMAND',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT 'STATUS',
  `notify` tinyint(4) NOT NULL DEFAULT '0' COMMENT 'NOTIFY',
  `notify_email` text NOT NULL COMMENT 'NOTIFY EMAIL',
  `timeout` smallint(6) NOT NULL DEFAULT '0' COMMENT 'TIMEOUT',
  `execute_times` int(11) NOT NULL DEFAULT '0' COMMENT 'EXECUTE TIMES',
  `prev_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'PREVIOUS TIME',
  `create_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'CREATE TIME',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_group_id` (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `t_task_group` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT 'USERID',
  `group_name` varchar(50) NOT NULL DEFAULT '' COMMENT 'GROUP NAME',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT 'DESCRIPTION',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT 'CREATE TIME',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `t_task_log` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `task_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'TASK ID',
  `output` mediumtext NOT NULL COMMENT 'OUTPUT',
  `error` text NOT NULL COMMENT 'ERROR',
  `status` tinyint(4) NOT NULL COMMENT 'STATUS',
  `process_time` int(11) NOT NULL DEFAULT '0' COMMENT 'PROCESS TIME',
  `create_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'CREATE TIME',
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `t_user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_name` varchar(20) NOT NULL DEFAULT '' COMMENT 'USER NAME',
  `email` varchar(50) NOT NULL DEFAULT '' COMMENT 'EMAIL',
  `password` char(32) NOT NULL DEFAULT '' COMMENT 'PASSWORD',
  `salt` char(10) NOT NULL DEFAULT '' COMMENT 'SALT',
  `last_login` int(11) NOT NULL DEFAULT '0' COMMENT 'LAST LOGIN',
  `last_ip` char(15) NOT NULL DEFAULT '' COMMENT 'LAST IP',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT 'STATUS，0=ENABLED -1=DISABLED',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_name` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `t_user` (`id`, `user_name`, `email`, `password`, `salt`, `last_login`, `last_ip`, `status`)
VALUES (1,'admin','admin@example.com','7fef6171469e80d32c0559f88b377245','',0,'',0);
