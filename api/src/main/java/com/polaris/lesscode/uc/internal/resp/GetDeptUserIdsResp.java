package com.polaris.lesscode.uc.internal.resp;

import lombok.Data;

import java.util.List;
import java.util.Map;

@Data
public class GetDeptUserIdsResp {

    private Map<Long, List<Long>> data;
}
