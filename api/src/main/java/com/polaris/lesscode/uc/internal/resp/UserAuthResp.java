package com.polaris.lesscode.uc.internal.resp;

import lombok.Data;

@Data
public class UserAuthResp {

	private String corpId;
	
	private Long orgId;
	
	private String outUserId;
	
	private String sourceChannel;
	
	private Long userId;
}
