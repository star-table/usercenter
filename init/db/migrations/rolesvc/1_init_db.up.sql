/*==============================================================*/
/* Table: ppm_rol_role                                          */
/*==============================================================*/
create table if not exists ppm_rol_role
(
    id            bigint      not null,
    org_id        bigint      not null default 0,
    role_group_id bigint      not null default 0,
    name          varchar(64) not null default '',
    creator       bigint      not null default 0,
    create_time   datetime    not null default CURRENT_TIMESTAMP,
    updator       bigint      not null default 0,
    update_time   datetime    not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    version       int         not null default 1,
    is_delete     tinyint     not null default 2,
    primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_rol_role_org_id                             */
/*==============================================================*/
create index index_ppm_rol_role_org_id on ppm_rol_role
    (
     org_id
        );

/*==============================================================*/
/* Index: index_ppm_rol_role_role_group_id                      */
/*==============================================================*/
create index index_ppm_rol_role_role_group_id on ppm_rol_role
    (
     role_group_id
        );

/*==============================================================*/
/* Index: index_ppm_rol_role_create_time                        */
/*==============================================================*/
create index index_ppm_rol_role_create_time on ppm_rol_role
    (
     create_time
        );

/*==============================================================*/
/* Table: ppm_rol_role_group                                    */
/*==============================================================*/
create table if not exists ppm_rol_role_group
(
    id          bigint      not null,
    org_id      bigint      not null default 0,
    name        varchar(64) not null default '',
    creator     bigint      not null default 0,
    create_time datetime    not null default CURRENT_TIMESTAMP,
    updator     bigint      not null default 0,
    update_time datetime    not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    version     int         not null default 1,
    is_delete   tinyint     not null default 2,
    primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_rol_role_group_org_id                       */
/*==============================================================*/
create index index_ppm_rol_role_group_org_id on ppm_rol_role_group
    (
     org_id
        );

/*==============================================================*/
/* Index: index_ppm_rol_role_group_create_time                  */
/*==============================================================*/
create index index_ppm_rol_role_group_create_time on ppm_rol_role_group
    (
     create_time
        );

/*==============================================================*/
/* Table: ppm_rol_role_user                                     */
/*==============================================================*/
create table if not exists ppm_rol_role_user
(
    id          bigint   not null,
    org_id      bigint   not null default 0,
    role_id     bigint   not null default 0,
    user_id     bigint   not null default 0,
    creator     bigint   not null default 0,
    create_time datetime not null default CURRENT_TIMESTAMP,
    updator     bigint   not null default 0,
    update_time datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    version     int      not null default 1,
    is_delete   tinyint  not null default 2,
    primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_rol_role_user_org_id                        */
/*==============================================================*/
create index index_ppm_rol_role_user_org_id on ppm_rol_role_user
    (
     org_id
        );

/*==============================================================*/
/* Index: index_ppm_rol_role_user_role_id                       */
/*==============================================================*/
create index index_ppm_rol_role_user_role_id on ppm_rol_role_user
    (
     role_id
        );

/*==============================================================*/
/* Index: index_ppm_rol_role_user_user_id                       */
/*==============================================================*/
create index index_ppm_rol_role_user_user_id on ppm_rol_role_user
    (
     user_id
        );

/*==============================================================*/
/* Index: index_ppm_rol_role_user_create_time                   */
/*==============================================================*/
create index index_ppm_rol_role_user_create_time on ppm_rol_role_user
    (
     create_time
        );


/* Table: lc_per_manage_group                                   */
create table if not exists lc_per_manage_group
(
    id              bigint      not null comment 'ID',
    org_id          bigint      not null default 0 comment '组织ID',
    lang_code       varchar(64) not null default '' comment 'LANG_CODE',
    name            varchar(64) not null default '' comment '管理组名称',
    user_ids        json comment '成员',
    app_package_ids json comment '管理应用id列表',
    usage_ids       json comment '操作权限列表',
    dept_ids        json comment '管理部门ID列表',
    role_ids        json comment '管理角色ID列表',
    creator         bigint      not null default 0,
    create_time     datetime    not null default CURRENT_TIMESTAMP,
    updator         bigint      not null default 0,
    update_time     datetime    not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    version         int         not null default 1,
    is_delete       tinyint     not null default 2,
    primary key (id)
) comment '管理组';
alter table lc_per_manage_group
    add app_ids json null comment 'AppId列表' after app_package_ids;

create index index_per_manage_group_org_id on lc_per_manage_group (org_id);
create index index_per_manage_group_name on lc_per_manage_group (name);
CREATE INDEX index_per_manage_group_user_ids ON lc_per_manage_group ((CAST(user_ids -> '$' AS UNSIGNED ARRAY)));
CREATE INDEX index_per_manage_group_app_package_ids ON lc_per_manage_group ((CAST(app_package_ids -> '$' AS UNSIGNED ARRAY)));
CREATE INDEX index_per_manage_group_app_ids ON lc_per_manage_group ((CAST(app_ids -> '$' AS UNSIGNED ARRAY)));
CREATE INDEX index_per_manage_group_dept_ids ON lc_per_manage_group ((CAST(dept_ids -> '$' AS UNSIGNED ARRAY)));
CREATE INDEX index_per_manage_group_role_ids ON lc_per_manage_group ((CAST(role_ids -> '$' AS UNSIGNED ARRAY)));
