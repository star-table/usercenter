package com.polaris.lesscode.uc.internal.enums;

public enum PayLevel {
    FREE(1, "免费版"),
    STANDARD(2, "标准版")
    ;

    private Integer code;

    private String desc;

    PayLevel(Integer code, String desc) {
        this.code = code;
        this.desc = desc;
    }

    public Integer getCode() {
        return code;
    }

    public String getDesc() {
        return desc;
    }
}
