package com.polaris.lesscode.uc.internal.req;

import io.swagger.annotations.ApiModelProperty;
import lombok.Data;

@Data
public class GetRepeatMemberReq {
    @ApiModelProperty("组织id")
    private Long orgId;
}
