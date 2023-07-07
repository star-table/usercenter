package com.polaris.lesscode.uc.internal.resp;

import lombok.Data;

import java.util.List;

@Data
public class GetManagerData {

    private String memberType;

    private Long memberId;

    private String langCode;

    private List<Long> appIds;

}
