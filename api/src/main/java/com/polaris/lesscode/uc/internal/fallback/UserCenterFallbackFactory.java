/**
 *
 */
package com.polaris.lesscode.uc.internal.fallback;

import com.polaris.lesscode.consts.ApplicationConsts;
import com.polaris.lesscode.feign.AbstractBaseFallback;
import com.polaris.lesscode.uc.internal.api.UserCenterApi;
import com.polaris.lesscode.uc.internal.req.*;
import com.polaris.lesscode.uc.internal.resp.*;
import com.polaris.lesscode.vo.Result;
import feign.hystrix.FallbackFactory;
import org.springframework.stereotype.Component;
import org.springframework.web.bind.annotation.RequestBody;

import java.util.ArrayList;
import java.util.List;

/**
 * @author admin
 *
 */
@Component
public class UserCenterFallbackFactory extends AbstractBaseFallback implements FallbackFactory<UserCenterApi> {

    @Override
    public UserCenterApi create(Throwable cause) {
        return new UserCenterApi() {

            @Override
            public Result<List<UserInfoResp>> getUserListByIds(UserListByIdsReq req) {
                return wrappDeal(ApplicationConsts.APPLICATION_USERCENTER, cause, () -> {
                    return Result.ok(new ArrayList<UserInfoResp>());
                });
            }

            @Override
            public Result<List<UserInfoResp>> getAllUserListByIds(UserListByIdsReq req) {
                return Result.ok(new ArrayList<UserInfoResp>());
            }

            @Override
            public Result<UserAuthorityResp> getUserAuthority(UserAuthorityReq req) {
                return wrappDeal(ApplicationConsts.APPLICATION_USERCENTER, cause, () -> {
                    return Result.ok(new UserAuthorityResp());
                });
            }

            @Override
            public Result<UserAuthResp> apiKeyAuth(String apiKey) {
                return wrappDeal(ApplicationConsts.APPLICATION_USERCENTER, cause, () -> {
                    return Result.ok(new UserAuthResp());
                });
            }

            @Override
            public Result<Boolean> addPkgToManageGroup(AddPkgReq reqParam) {
                return wrappDeal(ApplicationConsts.APPLICATION_USERCENTER, cause, () -> {
                    return Result.ok(true);
                });
            }

            @Override
            public Result<Boolean> deletePkgFromManageGroup(DeletePkgReq reqParam) {
                return wrappDeal(ApplicationConsts.APPLICATION_USERCENTER, cause, () -> {
                    return Result.ok(true);
                });
            }

            @Override
            public Result<Boolean> addAppToManageGroup(AddAppReq reqParam) {
                return wrappDeal(ApplicationConsts.APPLICATION_USERCENTER, cause, () -> {
                    return Result.ok(true);
                });
            }

            @Override
            public Result<Boolean> deleteAppFromManageGroup(DeleteAppReq reqParam) {
                return wrappDeal(ApplicationConsts.APPLICATION_USERCENTER, cause, () -> {
                    return Result.ok(true);
                });
            }

            @Override
            public Result<GetLeadersByDeptIdsResp> getLeadersByDeptIds(GetLeadersByDeptIdsReq req) {
                return wrappDeal(ApplicationConsts.APPLICATION_USERCENTER, cause, () -> {
                    return Result.ok(new GetLeadersByDeptIdsResp());
                });
            }

            @Override
            public Result<GetUserDeptIdsResp> getUserDeptIds(GetUserDeptIdsReq req) {
                return wrappDeal(ApplicationConsts.APPLICATION_USERCENTER, cause, () -> {
                    return Result.ok(new GetUserDeptIdsResp());
                });
            }

            @Override
            public Result<GetMemberSimpleInfoResp> getMemberSimpleInfo(GetMemberSimpleInfoReq req) {
                return wrappDeal(ApplicationConsts.APPLICATION_USERCENTER, cause, () -> {
                    return Result.ok(new GetMemberSimpleInfoResp());
                });
            }

            @Override
            public Result<RepeatMemberInfoResp> getRepeatMember(@RequestBody GetRepeatMemberReq req) {
                return wrappDeal(ApplicationConsts.APPLICATION_USERCENTER, cause, () -> {
                    return Result.ok(new RepeatMemberInfoResp());
                });
            }

            @Override
            public Result<GetManagerResp> getManager(GetManagerReq req) {
                return wrappDeal(ApplicationConsts.APPLICATION_USERCENTER, cause, () -> {
                    return Result.ok(new GetManagerResp());
                });
            }

            @Override
            public Result<GetRoleUserIdsResp> getRoleUserIds(GetRoleUserIdsReq req) {
                return wrappDeal(ApplicationConsts.APPLICATION_USERCENTER, cause, () -> {
                    return Result.ok(new GetRoleUserIdsResp());
                });
            }

            @Override
            public Result<GetDeptUserIdsResp> getDeptUserIds(GetDeptUserIdsReq req) {
                return wrappDeal(ApplicationConsts.APPLICATION_USERCENTER, cause, () -> {
                    return Result.ok(new GetDeptUserIdsResp());
                });
            }

            @Override
            public Result<GetOrgInfoResp> getOrgInfo(Long orgId) {
                return wrappDeal(ApplicationConsts.APPLICATION_USERCENTER, cause, () -> {
                    return Result.ok(new GetOrgInfoResp());
                });
            }

            @Override
            public Result<?> addOutCollaborator(AddOutCollaboratorReq req) {
                return wrappDeal(ApplicationConsts.APPLICATION_USERCENTER, cause, () -> {
                    return Result.ok(null);
                });
            }

            @Override
            public Result<List<RoleInfoResp>> getRoleListByIds(RoleListByIdsReq req) {
                return wrappDeal(ApplicationConsts.APPLICATION_USERCENTER, cause, () -> {
                    return Result.ok(new ArrayList<RoleInfoResp>());
                });
            }

            @Override
            public Result<List<RoleInfoResp>> getAllRoleListByIds(RoleListByIdsReq req) {
                return wrappDeal(ApplicationConsts.APPLICATION_USERCENTER, cause, () -> {
                    return Result.ok(new ArrayList<RoleInfoResp>());
                });
            }

            @Override
            public Result<List<DeptInfoResp>> getDeptListByIds(DeptListByIdsReq req) {
                return wrappDeal(ApplicationConsts.APPLICATION_USERCENTER, cause, () -> {
                    return Result.ok(new ArrayList<DeptInfoResp>());
                });
            }

            @Override
            public Result<List<DeptInfoResp>> getAllDeptListByIds(DeptListByIdsReq req) {
                return Result.ok(new ArrayList<DeptInfoResp>());
            }

            @Override
            public Result<UserAuthResp> auth(String token) {
                return wrappDeal(ApplicationConsts.APPLICATION_USERCENTER, cause, () -> {
                    return Result.ok(new UserAuthResp());
                });
            }

            @Override
            public Result<UserAuthResp> authCheckStatus(String token) {
                return wrappDeal(ApplicationConsts.APPLICATION_USERCENTER, cause, () -> {
                    return Result.ok(new UserAuthResp());
                });
            }
        };
    }

}
