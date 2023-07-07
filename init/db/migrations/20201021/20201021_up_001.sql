/* Table: lc_org_position                                   */
begin;
create table if not exists lc_org_position
(
    id              bigint       not null default 0 comment 'ID（全局ID）',
    org_id          bigint       not null default 0 comment '组织ID',
    org_position_id bigint       not null default 0 comment '职级ID（局部的职级ID）1为主管 2为成员 其他为自定义',
    name            varchar(64)  not null default '' comment '名称',
    position_level  int          not null default 0 comment '级别',
    remark          varchar(255) not null default '' comment '说明',
    status          int          not null default 1 comment '状态 1启用 2禁用',
    creator         bigint       not null default 0,
    create_time     datetime     not null default CURRENT_TIMESTAMP,
    updator         bigint       not null default 0,
    update_time     datetime     not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    version         int          not null default 1,
    is_delete       tinyint      not null default 2,
    primary key (id)
) comment '职位信息';

create index index_org_position_org_id on lc_org_position (org_id);
create index index_org_position_org_position_id on lc_org_position (org_position_id);
create index index_org_position_name on lc_org_position (name);
commit;

/* 用户部门新增职位字段 */
begin;
-- 默认2为成员
alter table ppm_org_user_department
    add org_position_id bigint default 2 not null comment '职位ID' after is_leader;

create index index_user_department_org_position_id on ppm_org_user_department (org_position_id);
commit;

begin;
update ppm_org_user_department
set org_position_id = 2
where 1 = 1;
commit;