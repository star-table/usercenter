package com.polaris.lesscode.uc.internal.req;

import io.swagger.annotations.ApiModelProperty;
import lombok.Data;

@Data
public class AddOutCollaboratorReq {

    @ApiModelProperty("组织id")
    private Long orgId;

    @ApiModelProperty("要加入的用户id")
    private Long userId;
}
