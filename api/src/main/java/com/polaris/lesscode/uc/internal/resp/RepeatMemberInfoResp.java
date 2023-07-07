package com.polaris.lesscode.uc.internal.resp;

import io.swagger.annotations.ApiModelProperty;
import lombok.Data;

import java.util.List;

@Data
public class RepeatMemberInfoResp {
    @ApiModelProperty("成员信息")
    private List<RepeatMemberInfo> user;

    @ApiModelProperty("部门信息")
    private List<RepeatMemberInfo> department;

    @ApiModelProperty("角色信息")
    private List<RepeatMemberInfo> role;
}
