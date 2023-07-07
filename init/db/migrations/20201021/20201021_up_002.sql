/* 用户表新增字段 */
alter table ppm_org_user_organization
    add emp_no varchar(255) default '' not null comment '工号' after `status`;
alter table ppm_org_user_organization
    add weibo_ids varchar(1024) default ''not null comment '微博ID array ,分割' after emp_no;

/* 创建索引 */
create index index_ppm_org_user_emp_no on ppm_org_user_organization (emp_no);

-- 手机区号默认值
alter table ppm_org_user alter column mobile_region set default '';
