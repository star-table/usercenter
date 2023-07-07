package com.polaris.lesscode.uc.internal.resp;

import lombok.Data;

import java.util.List;

@Data
public class RepeatMemberInfo {
    private Long id;

    private String name;

    private List<String> department;
}
