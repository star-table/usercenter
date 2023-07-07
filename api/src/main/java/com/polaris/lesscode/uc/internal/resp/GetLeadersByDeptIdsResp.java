package com.polaris.lesscode.uc.internal.resp;

import lombok.Data;

import java.util.List;

/**
 * 通过bumenid批量获取leaders响应结构体
 *
 * @author Nico
 * @date 2021/4/26 17:10
 */
@Data
public class GetLeadersByDeptIdsResp {

    private List<DepartmentLeader> leaders;

}
