package com.polaris.lesscode.uc.internal.req;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serializable;

/**
 * 用户权限信息 请求参数（内部调用）
 *
 * @author roamer
 * @version v1.0
 * @date 2020-09-10 18:34
 */
@Data
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class UserAuthorityReq implements Serializable {

    private static final long serialVersionUID = 795870311167326011L;

    @ApiModelProperty("组织ID")
    private Long orgId;

    @ApiModelProperty("用户信息")
    private Long userId;

}
