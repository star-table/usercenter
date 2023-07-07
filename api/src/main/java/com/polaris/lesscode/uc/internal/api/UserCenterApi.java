package com.polaris.lesscode.uc.internal.api;

import com.polaris.lesscode.uc.internal.req.*;
import com.polaris.lesscode.uc.internal.resp.*;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.*;

import com.polaris.lesscode.vo.Result;

import java.util.List;

@RequestMapping("/usercenter/inner/api/v1")
public interface UserCenterApi {

    /**
     * code:
     * 300301: 身份认证异常
     * 0: 成功
     *
     * @param token
     * @return
     */
    @PostMapping(value = "/user/auth", consumes = MediaType.APPLICATION_FORM_URLENCODED_VALUE)
    Result<UserAuthResp> auth(@RequestParam("token") String token);

    /**
     * code:
     * 300301: 身份认证异常
     * 0: 成功
     *
     * @param token
     * @return
     */
    @PostMapping(value = "/user/auth-check-status", consumes = MediaType.APPLICATION_FORM_URLENCODED_VALUE)
    Result<UserAuthResp> authCheckStatus(@RequestParam("token") String token);

    /**
     * code:
     * 400: 参数异常
     * 500: 系统异常
     * 0: 成功
     *
     * @param req
     * @return
     */
    @PostMapping(value = "/user/getListByIds", consumes = MediaType.APPLICATION_JSON_VALUE)
    Result<List<UserInfoResp>> getUserListByIds(@RequestBody() UserListByIdsReq req);

       /**
     * code:
     * 400: 参数异常
     * 500: 系统异常
     * 0: 成功
     *
     * @param req
     * @return
     */
    @PostMapping(value = "/user/getAllListByIds", consumes = MediaType.APPLICATION_JSON_VALUE)
    Result<List<UserInfoResp>> getAllUserListByIds(@RequestBody() UserListByIdsReq req);

    /**
     * code:
     * 400: 参数异常
     * 500: 系统异常
     * 0: 成功
     *
     * @param req
     * @return
     */
    @PostMapping(value = "/role/getListByIds", consumes = MediaType.APPLICATION_JSON_VALUE)
    Result<List<RoleInfoResp>> getRoleListByIds(@RequestBody() RoleListByIdsReq req);

    /**
     * code:
     * 400: 参数异常
     * 500: 系统异常
     * 0: 成功
     *
     * @param req
     * @return
     */
    @PostMapping(value = "/role/getAllListByIds", consumes = MediaType.APPLICATION_JSON_VALUE)
    Result<List<RoleInfoResp>> getAllRoleListByIds(@RequestBody() RoleListByIdsReq req);

    /**
     * code:
     * 400: 参数异常
     * 500: 系统异常
     * 0: 成功
     *
     * @param req
     * @return
     */
    @PostMapping(value = "/dept/getListByIds", consumes = MediaType.APPLICATION_JSON_VALUE)
    Result<List<DeptInfoResp>> getDeptListByIds(@RequestBody() DeptListByIdsReq req);

    /**
         * code:
         * 400: 参数异常
         * 500: 系统异常
         * 0: 成功
         *
         * @param req
         * @return
         */
        @PostMapping(value = "/dept/getAllListByIds", consumes = MediaType.APPLICATION_JSON_VALUE)
        Result<List<DeptInfoResp>> getAllDeptListByIds(@RequestBody() DeptListByIdsReq req);

    /**
     * code:
     * 400: 参数异常
     * 500: 系统异常
     * 0: 成功
     *
     * @param req
     * @return
     */
    @PostMapping(value = "/user/getUserAuthority", consumes = MediaType.APPLICATION_JSON_VALUE)
    Result<UserAuthorityResp> getUserAuthority(@RequestBody() UserAuthorityReq req);


    /**
     * code:
     * 204000: OpenAPI签名校验失败
     * 0: 成功
     *
     * @param apiKey Open ApiKey
     * @return {@code UserAuthResp}
     */
    @PostMapping(value = "/auth/api-key-auth", consumes = MediaType.APPLICATION_FORM_URLENCODED_VALUE)
    Result<UserAuthResp> apiKeyAuth(@RequestParam("apiKey") String apiKey);


    /**
     * code:
     * 400: 参数异常
     * 500: 系统异常
     * 0: 成功
     *
     * @param reqParam 参数
     * @return {@code true}
     */
    @PostMapping(value = "manage-group/add-pkg", consumes = MediaType.APPLICATION_JSON_VALUE)
    Result<Boolean> addPkgToManageGroup(@RequestBody AddPkgReq reqParam);


    /**
     * code:
     * 400: 参数异常
     * 500: 系统异常
     * 0: 成功
     *
     * @param reqParam 参数
     * @return {@code true}
     */
    @PostMapping(value = "manage-group/delete-pkg", consumes = MediaType.APPLICATION_JSON_VALUE)
    Result<Boolean> deletePkgFromManageGroup(@RequestBody DeletePkgReq reqParam);


    /**
     * code:
     * 400: 参数异常
     * 500: 系统异常
     * 0: 成功
     *
     * @param reqParam 参数
     * @return {@code true}
     */
    @PostMapping(value = "manage-group/add-app", consumes = MediaType.APPLICATION_JSON_VALUE)
    Result<Boolean> addAppToManageGroup(@RequestBody AddAppReq reqParam);


    /**
     * code:
     * 400: 参数异常
     * 500: 系统异常
     * 0: 成功
     *
     * @param reqParam 参数
     * @return {@code true}
     */
    @PostMapping(value = "manage-group/delete-app", consumes = MediaType.APPLICATION_JSON_VALUE)
    Result<Boolean> deleteAppFromManageGroup(@RequestBody DeleteAppReq reqParam);

    /**
     * 通过部门id批量获取部门下的leaders
     *
     * @Author Nico
     * @Date 2021/4/26 17:13
     **/
    @PostMapping(value = "/dept/getLeadersByDeptIds", consumes = MediaType.APPLICATION_JSON_VALUE)
    Result<GetLeadersByDeptIdsResp> getLeadersByDeptIds(@RequestBody GetLeadersByDeptIdsReq req);

    /**
     * 获取用户部门ids
     **/
    @PostMapping(value = "/dept/getUserDeptIds", consumes = MediaType.APPLICATION_JSON_VALUE)
    Result<GetUserDeptIdsResp> getUserDeptIds(@RequestBody GetUserDeptIdsReq req);

    /**
     * 获取用户、部门、角色的k，v结构，k为名称，v为id
     *
     * @Author Nico
     * @Date 2021/5/12 13:58
     **/
    @PostMapping(value = "/user/getMemberSimpleInfo", consumes = MediaType.APPLICATION_JSON_VALUE)
    Result<GetMemberSimpleInfoResp> getMemberSimpleInfo(@RequestBody GetMemberSimpleInfoReq req);

    /**
     * 获取重复成员信息 用户、部门、角色
     *
     */
    @PostMapping(value = "/user/getRepeatMember", consumes = MediaType.APPLICATION_JSON_VALUE)
    Result<RepeatMemberInfoResp> getRepeatMember(@RequestBody GetRepeatMemberReq req);

    /**
     * 获取管理组
     *
     **/
    @PostMapping(value = "/manage-group/getManagerInfo")
    Result<GetManagerResp> getManager(@RequestBody GetManagerReq req);

    /**
     * 获取角色成员id列表
     **/
    @PostMapping(value = "/role/getUserIds")
    Result<GetRoleUserIdsResp> getRoleUserIds(@RequestBody GetRoleUserIdsReq req);

    /**
     * 获取部门成员id列表
     **/
    @PostMapping(value = "/dept/getUserIds")
    Result<GetDeptUserIdsResp> getDeptUserIds(@RequestBody GetDeptUserIdsReq req);

    /**
     * 获取组织信息
     **/
    @GetMapping(value = "/org/info")
    Result<GetOrgInfoResp> getOrgInfo(@RequestParam(value="orgId") Long orgId);

    /**
     * 添加外部协作人
     */
    @PostMapping(value = "/org/add-out-collaborator")
    Result<?> addOutCollaborator(@RequestBody AddOutCollaboratorReq req);

}
