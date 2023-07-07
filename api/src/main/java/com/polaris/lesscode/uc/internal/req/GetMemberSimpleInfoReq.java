package com.polaris.lesscode.uc.internal.req;

import io.swagger.annotations.ApiModelProperty;
import lombok.Data;

@Data
public class GetMemberSimpleInfoReq {

    @ApiModelProperty("组织id")
    private Long orgId;

    @ApiModelProperty("类型(1成员2部门3角色)")
    private Integer type;

    private Integer needDelete;
}
