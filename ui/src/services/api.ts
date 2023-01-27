import { request } from 'umi';
import type { MenuDataItem } from '@ant-design/pro-layout';

/** 获取当前的用户 GET /api/currentUser */
export async function currentUser() {
  return request<API.Response<API.CurrentUser>>('/api/rbac/user/current', {
    method: 'GET',
  });
}
export async function userMenu() {
  return request<API.Response<MenuDataItem[]>>(`/api/rbac/user/menu`, {
    method: 'GET',
  });
}
