package com.polaris.lesscode.uc.internal.resp;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import lombok.Data;

import java.io.Serializable;
import java.util.List;

/**
 * 用户信息 响应模型(由用户服务返回)
 *
 * @author roamer
 * @version v1.0
 * @date 2020-09-10 18:34
 */
@ApiModel("用户信息 响应模型(由用户服务返回)")
@Data
public class UserInfoResp implements Serializable {

    private static final long serialVersionUID = 795870311167326011L;

    @ApiModelProperty("ID")
    private Long id;

    @ApiModelProperty("昵称")
    private String name;

    @ApiModelProperty("昵称拼音")
    private String namePy;

    @ApiModelProperty("用户头像")
    private String avatar;

    @ApiModelProperty("状态")
    private Integer status;

    @ApiModelProperty("是否删除")
    private Integer isDelete;

    @ApiModelProperty("用户类型，1：内部，2：外部")
    private Integer type;

    private List<RoleBindData> roleList;

    private List<DeptBindData> deptList;

    public UserInfoResp() {
        status = 1;
        isDelete = 2;
    }

}
