import { request } from 'umi';

export async function permissionIndex(payload: API.PageParams) {
  return request<API.Response<PermissionListItem[]>>(`/api/rbac/permission/index`, {
    method: 'GET',
    params: {
      ...payload,
    },
  });
}
export async function addPermission(payload: PermissionListItem) {
  return request<API.Response<void>>(`/api/rbac/permission/create`, {
    method: 'POST',
    data: {
      ...payload,
    },
  });
}

export async function updatePermission(payload: PermissionListItem) {
  return request<API.Response<void>>(`/api/rbac/permission/update`, {
    method: 'POST',
    data: {
      ...payload,
    },
  });
}

export async function removePermission(id: number) {
  return request<API.Response<void>>('/api/rbac/permission/delete', {
    method: 'POST',
    params: {
      id,
    },
  });
}
