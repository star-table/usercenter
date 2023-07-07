package com.polaris.lesscode.uc.internal.resp;

import io.swagger.annotations.ApiModelProperty;
import lombok.Data;

import java.lang.reflect.Member;
import java.util.List;
import java.util.Map;

@Data
public class GetMemberSimpleInfoResp {

    @ApiModelProperty("成员信息列表")
    private List<MemberSimpleInfo> data;

}
