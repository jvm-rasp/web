import { request } from 'umi';

export async function roleIndex(
  params: {
    // query
    /** 当前的页码 */
    current?: number;
    /** 页面的容量 */
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  return request<API.Response<RoleListItem[]>>('/api/rbac/role/index', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

export async function updateRole(payload: RoleListItem) {
  return request<API.Response<void>>('/api/rbac/role/update', {
    method: 'POST',
    data: {
      ...payload,
    },
  });
}
export async function addRole(payload: RoleListItem) {
  return request<API.Response<void>>('/api/rbac/role/create', {
    method: 'POST',
    data: {
      ...payload,
    },
  });
}
export async function removeRole(id: number) {
  return request<API.Response<void>>('/api/rbac/role/delete', {
    method: 'POST',
    params: {
      id,
    },
  });
}
export async function authRolePermission(payload: { role_id: number; permission_ids: number[] }) {
  return request<API.Response<void>>(`/api/rbac/role/auth/permission`, {
    method: 'POST',
    data: {
      ...payload,
    },
  });
}

export async function rolePermissionIds(role_id: number) {
  return request<API.Response<number[]>>(`/api/rbac/role/permission/ids`, {
    method: 'GET',
    params: {
      role_id,
    },
  });
}
