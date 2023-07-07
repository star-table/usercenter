/**
 * 
 */
package com.polaris.lesscode.uc.internal.feign;

import org.springframework.cloud.openfeign.FeignClient;

import com.polaris.lesscode.consts.ApplicationConsts;
import com.polaris.lesscode.uc.internal.api.UserCenterApi;
import com.polaris.lesscode.uc.internal.fallback.UserCenterFallbackFactory;

/**
 * @author Bomb.
 *
 */
@FeignClient(name = ApplicationConsts.APPLICATION_USERCENTER, fallbackFactory = UserCenterFallbackFactory.class)
public interface UserCenterProvider extends UserCenterApi {

}
