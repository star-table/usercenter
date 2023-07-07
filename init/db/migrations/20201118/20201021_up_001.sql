
-- 登陆日志 增加user-agent字段
alter table `ppm_org_user_login_record`
    add `user_agent` varchar(255) default '' not null comment 'user-agent' after `source_channel`;
alter table `ppm_org_user_login_record`
    add `msg` varchar(255) default '' not null comment 'msg' after `user_agent`;


-- 添加最后修改密码时间
alter table `ppm_org_user`
    add `last_edit_pwd_time` datetime default '1970-01-01 00:00:00' not null comment '最后修改密码时间' after `login_fail_count`;