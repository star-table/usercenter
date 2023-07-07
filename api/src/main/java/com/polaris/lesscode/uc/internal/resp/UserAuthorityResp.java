package com.polaris.lesscode.uc.internal.resp;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import lombok.Data;

import java.io.Serializable;
import java.util.ArrayList;
import java.util.List;

/**
 * 用户关联信息 响应模型(由用户服务返回)
 *
 * @author roamer
 * @version v1.0
 * @date 2020-09-10 18:34
 */
@ApiModel("用户关联信息 响应模型(由用户服务返回)")
@Data
public class UserAuthorityResp implements Serializable {

    private static final long serialVersionUID = 795870311167326011L;
    @ApiModelProperty("部门信息")
    private Long orgId;

    @ApiModelProperty("用户信息")
    private Long userId;

    @ApiModelProperty("是否是组织拥有者")
    private Boolean isOrgOwner;

    @ApiModelProperty("是否是系统管理员")
    private Boolean isSysAdmin;

    @ApiModelProperty("是否是子管理员")
    private Boolean isSubAdmin;

    @ApiModelProperty("是否是外部协作人")
    private Boolean isOutCollaborator;

    @ApiModelProperty("所属部门Id列表")
    private List<Long> refDeptIds;

    @ApiModelProperty("所属角色Id列表")
    private List<Long> refRoleIds;

    @ApiModelProperty("管理部门的权限")
    private Boolean hasDeptOptAuth;

    @ApiModelProperty("管理角色的权限")
    private Boolean hasRoleOptAuth;

    @ApiModelProperty("管理应用包的权限")
    private Boolean hasAppPackageOptAuth;

    @ApiModelProperty("管理的角色ID列表")
    private List<Long> manageRoles;

    @ApiModelProperty("管理的部门ID列表")
    private List<Long> manageDepts;

    @ApiModelProperty("管理的应用包ID列表")
    private List<Long> manageAppPackages;

    @ApiModelProperty("管理的应用ID列表")
    private List<Long> manageApps;

    @ApiModelProperty("操作权限项")
    private List<String> optAuth;

    public UserAuthorityResp() {
        refDeptIds = new ArrayList<>();
        refDeptIds = new ArrayList<>();
        isOrgOwner = false;
        isSysAdmin = false;
        isSubAdmin = false;
        hasDeptOptAuth = false;
        hasRoleOptAuth = false;
        hasAppPackageOptAuth = false;
        manageRoles = new ArrayList<>();
        manageDepts = new ArrayList<>();
        manageAppPackages = new ArrayList<>();
        manageApps = new ArrayList<>();
    }

    @ApiModel("管理项信息 响应模型(由用户服务返回)")
    @Data
    public static class ObjOptAuthItem {
        @ApiModelProperty("ID")
        private Long id;

        @ApiModelProperty("是否具有管理权限")
        private Boolean hasManage;

        public ObjOptAuthItem() {
            hasManage = false;
        }
    }

}
