package com.polaris.lesscode.uc.internal.resp;

import lombok.Data;

@Data
public class MemberSimpleInfo {

    private Long id;

    private String name;

    private Long parentId;

    private String avatar;

    private Integer status;
}
