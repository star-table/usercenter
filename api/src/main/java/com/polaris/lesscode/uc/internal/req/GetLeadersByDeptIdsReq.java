package com.polaris.lesscode.uc.internal.req;

import lombok.Data;

import java.util.List;

/**
 * 通过部门id批量获取leaders
 *
 * @author Nico
 * @date 2021/4/26 17:09
 */
@Data
public class GetLeadersByDeptIdsReq {

    /**
     * 组织id
     **/
    private Long orgId;

    /**
     * 部门id列表
     **/
    private List<Long> deptIds;
}
