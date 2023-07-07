package com.polaris.lesscode.uc.internal.resp;

import lombok.Data;

import java.util.List;
import java.util.Map;

@Data
public class GetOrgInfoResp {

    private Long orgId;

    private String orgName;

    private Long orgOwnerId;

    private String outOrgId;

    private String sourceChannel;

    private Integer payLevel;
}
