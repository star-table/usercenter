package com.polaris.lesscode.uc.internal.resp;

import io.swagger.annotations.ApiModelProperty;
import lombok.Data;

/**
 * @author Nico
 * @date 2021/4/26 17:11
 */
@Data
public class DepartmentLeader {

    @ApiModelProperty("部门id")
    private Long departmentId;

    @ApiModelProperty("leaderId")
    private Long leaderId;
}
