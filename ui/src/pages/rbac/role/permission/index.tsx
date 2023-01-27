import React, { useEffect, useRef, useState } from 'react';
import { Link, useHistory } from 'umi';
import { FooterToolbar, PageContainer } from '@ant-design/pro-layout';
import ProTable, { ActionType, ProColumns } from '@ant-design/pro-table';
import IconFont from '@/components/Font/Iconfont';
import { authRolePermission, rolePermissionIds } from '@/pages/rbac/role/service';
import { Button, message } from 'antd';
import { permissionIndex } from '@/pages/rbac/permission/service';

const authPermissionHandler = async (roleId: number, permissionIds: number[]) => {
  if (!permissionIds) return;
  const payload = { role_id: roleId, permission_ids: permissionIds };
  const hide = message.loading('正在分配权限...');
  const response = await authRolePermission(payload);
  hide();
  const { code, msg } = response;
  if (code === 200) {
    message.success(msg);
    return;
  }
  message.error(msg);
};

const RolePermissionIndex: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const history = useHistory();
  // @ts-ignore
  const { role } = history.location.state;
  const [selectedRowKeysState, setSelectedRowKeys] = useState<number[]>([]);
  const getRolePermissionIds = async (roleId: number) => {
    const response = await rolePermissionIds(roleId);
    setSelectedRowKeys(response.data || []);
  };
  useEffect(() => {
    getRolePermissionIds(role.id);
  }, []);
  const columns: ProColumns<PermissionListItem>[] = [
    {
      title: '名称',
      dataIndex: 'name',
    },
    {
      title: '网址',
      dataIndex: 'url',
      search: false,
    },
    {
      title: ' 图标',
      dataIndex: 'icon',
      search: false,
      render: (_, record) => {
        return <IconFont type={record!.icon} />;
      },
    },
    {
      title: '导航菜单',
      dataIndex: 'hide_in_menu',
      valueEnum: {
        1: { text: '隐藏', status: 'Success' },
        2: { text: '显示', status: 'Processing' },
      },
    },
    {
      title: '优先级',
      dataIndex: 'priority',
      search: false,
    },
  ];

  return (
    <PageContainer
      header={{
        title: `${role.name}权限分配`,
        breadcrumb: {
          routes: [
            {
              path: '/rbac',
              breadcrumbName: '权限管理',
            },
            {
              path: '/rbac/role/index',
              breadcrumbName: '角色管理',
            },
            {
              path: '',
              breadcrumbName: '权限分配',
            },
          ],
          itemRender: (route, params, routes) => {
            const last = routes.indexOf(route) === routes.length - 1;
            return last ? (
              <span>{route.breadcrumbName}</span>
            ) : (
              <Link to={route.path}>{route.breadcrumbName}</Link>
            );
          },
        },
      }}
    >
      <ProTable<PermissionListItem>
        headerTitle={'分配权限'}
        actionRef={actionRef}
        rowKey={'id'}
        rowSelection={{
          selectedRowKeys: selectedRowKeysState,
          onSelect: (record, selected) => {
            let result = [];
            let childrenIds: number[] = [];
            if (selected) {
              if (record.children) {
                childrenIds = record.children.map((item) => item.id!);
              }
              result = [...selectedRowKeysState, record.id!, ...childrenIds];
            } else {
              if (record.children) {
                childrenIds = record.children.map((item) => item.id!);
              }
              result = selectedRowKeysState.filter((item) => item !== record.id!);
              result = [...result].filter((x) => [...childrenIds].every((y) => y !== x));
            }
            setSelectedRowKeys(result);
          },
        }}
        request={async (params) => {
          const response = await permissionIndex(params);
          return {
            success: response.code === 200,
            data: response.data,
          };
        }}
        search={false}
        columns={columns}
        pagination={false}
      />
      {selectedRowKeysState.length > 0 && (
        <FooterToolbar
          extra={
            <div>
              已选择 <a style={{ fontWeight: 600 }}>{selectedRowKeysState.length}</a> 项权限
            </div>
          }
        >
          <Button
            type="primary"
            onClick={async () => {
              await authPermissionHandler(role.id, selectedRowKeysState);
            }}
          >
            提交
          </Button>
        </FooterToolbar>
      )}
    </PageContainer>
  );
};

export default RolePermissionIndex;
