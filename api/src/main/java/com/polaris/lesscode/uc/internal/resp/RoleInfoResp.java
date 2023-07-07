package com.polaris.lesscode.uc.internal.resp;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import lombok.Data;

import java.io.Serializable;

/**
 * 角色信息 响应模型(由用户服务返回)
 *
 * @author roamer
 * @version v1.0
 * @date 2020-09-10 18:35
 */
@ApiModel("角色信息 响应模型(由用户服务返回)")
@Data
public class RoleInfoResp implements Serializable {


    private static final long serialVersionUID = 5220074020337895239L;
    @ApiModelProperty("ID")
    private Long id;

    @ApiModelProperty("角色名称")
    private String name;

    @ApiModelProperty("状态")
    private Integer status;

    @ApiModelProperty("是否删除")
    private Integer isDelete;

    public RoleInfoResp(){
          status = 1;
          isDelete = 2;
    }
}
