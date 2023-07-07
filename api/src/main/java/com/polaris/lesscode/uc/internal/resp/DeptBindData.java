package com.polaris.lesscode.uc.internal.resp;

import lombok.Data;

@Data
public class DeptBindData {

    private Long departmentId;

    private String departmentName;

    private String outOrgDepartmentId;

    private String OutOrgDepartmentCode;

    private String outOrgDepartmentParentId;

    private Integer isLeader;

    private Long positionId;

    private String positionName;

    private Integer positionLevel;
}
