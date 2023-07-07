package com.polaris.lesscode.uc.internal.resp;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import lombok.Data;

import java.io.Serializable;

/**
 * 部门信息 响应模型(由用户服务返回)
 *
 * @author roamer
 * @version v1.0
 * @date 2020-09-10 18:35
 */
@ApiModel("部门信息 响应模型(由用户服务返回)")
@Data
public class DeptInfoResp implements Serializable {
    private static final long serialVersionUID = 8205685181898168488L;

    @ApiModelProperty("ID")
    private Long id;

    @ApiModelProperty("父部门")
    private Long parentId;

    @ApiModelProperty("部门名称")
    private String name;

    @ApiModelProperty("状态")
    private Integer status;

    @ApiModelProperty("是否删除")
    private Integer isDelete;

    public DeptInfoResp() {
        status = 1;
        isDelete = 2;
    }
}
