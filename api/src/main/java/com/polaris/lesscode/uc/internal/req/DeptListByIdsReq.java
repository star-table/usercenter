package com.polaris.lesscode.uc.internal.req;

import io.swagger.annotations.ApiModelProperty;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

/**
 * 组织和ID参数
 *
 * @author roamer
 * @version v1.0
 * @date 2020-09-14 15:37
 */
@Data
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class DeptListByIdsReq {

    @ApiModelProperty("组织ID")
    private Long orgId;

    @ApiModelProperty("ID列表")
    private List<Long> ids;
}
