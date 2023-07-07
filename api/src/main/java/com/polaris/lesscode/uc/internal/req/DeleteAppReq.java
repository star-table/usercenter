package com.polaris.lesscode.uc.internal.req;

import io.swagger.annotations.ApiModelProperty;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * @author roamer
 * @version v1.0
 * @date 2020-11-24 14:19
 */
@Data
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class DeleteAppReq {

    @ApiModelProperty("组织ID")
    private Long orgId;

    @ApiModelProperty("用户ID")
    private Long userId;

    @ApiModelProperty("应用ID")
    private Long appId;

}
