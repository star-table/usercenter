package com.polaris.lesscode.uc.internal.resp;

import lombok.Data;

import java.util.List;

@Data
public class GetUserDeptIdsResp {

    private List<Long> deptIds;
}
