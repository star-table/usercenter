package com.polaris.lesscode.uc.internal.req;

import io.swagger.annotations.ApiModelProperty;
import lombok.Data;

@Data
public class GetUserDeptIdsReq {

    @ApiModelProperty("组织ID")
    private Long orgId;

    @ApiModelProperty("用户信息")
    private Long userId;
}
