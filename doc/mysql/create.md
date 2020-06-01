## file_mate
文件详情
```sql
DROP TABLE IF EXISTS `file_mate`;
CREATE TABLE `file_mate` (
  `id` bigint(11) unsigned NOT NULL AUTO_INCREMENT,
  `filesha256` varchar(512) NOT NULL COMMENT '文件hash',
  `size` bigint(11) unsigned zerofill NOT NULL COMMENT '文件大小',
  `filename` varchar(128) NOT NULL DEFAULT '' COMMENT '文件名',
  `location` varchar(128) NOT NULL DEFAULT '' COMMENT '文件路径',
  `created_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=71 DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
```

## user_info
保存用户信息
```sql
DROP TABLE IF EXISTS `user_info`;
CREATE TABLE `user_info` (
  `id` bigint(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_name` varchar(32) NOT NULL COMMENT '用户名',
  `password` char(64) NOT NULL COMMENT '密码 ',
  `phone` varchar(11) NOT NULL COMMENT '手机号',
  `status` tinyint(1) NOT NULL DEFAULT '2' COMMENT '用户状态 1.启用 2.禁用 3.锁定',
  `created_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_phone` (`phone`) USING BTREE COMMENT '手机号唯一'
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1
```