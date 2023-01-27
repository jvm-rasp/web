import React, { useRef, useState } from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import type { ActionType, ProColumns } from '@ant-design/pro-table';
import ProTable from '@ant-design/pro-table';
import { Button, message, Modal, Space } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { addRole, removeRole, roleIndex, updateRole } from './service';
import { ModalForm, ProFormText, ProFormTextArea } from '@ant-design/pro-form';
import { history } from 'umi';
import { Access, useAccess } from '@@/plugin-access/access';

const RoleList: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const [createModalVisible, handleModalVisible] = useState<boolean>(false);
  const [current, handleCurrent] = useState<Partial<RoleListItem>>({});
  const [updateModalVisible, handleUpdateModalVisible] = useState<boolean>(false);
  const access: API.UserAccessItem = useAccess();

  const roleDelete = (id: number) => {
    Modal.confirm({
      title: '删除角色',
      content: '您确定要删除该角色吗？删除后不能恢复！',
      centered: true,
      okType: 'danger',
      onOk: async () => {
        const response = await removeRole(id);
        const { code, msg } = response;
        if (code !== 200) {
          message.error(msg);
          return;
        }
        message.success(msg);
        if (actionRef.current) {
          actionRef.current?.reload();
        }
      },
    });
  };
  const columns: ProColumns<RoleListItem>[] = [
    {
      title: 'id',
      dataIndex: 'id',
      search: false,
    },
    {
      title: '名称',
      dataIndex: 'name',
      search: false,
    },
    {
      title: '编码',
      dataIndex: 'code',
      search: false,
    },
    {
      title: '备注',
      dataIndex: 'remark',
      search: false,
    },
    {
      title: '操作',
      render: (_, record) => {
        return (
          <Space>
            <Access accessible={access.rbacRoleUpdate!}>
              <Button
                type={'primary'}
                size={'small'}
                onClick={() => {
                  handleCurrent(record);
                  handleUpdateModalVisible(true);
                }}
              >
                编辑
              </Button>
            </Access>
            <Access accessible={access.rbacRolePermissionIndex!}>
              <Button
                type={'dashed'}
                size={'small'}
                onClick={() => {
                  history.push(`/rbac/role/permission/index`, { role: record });
                }}
              >
                权限
              </Button>
            </Access>
            <Access accessible={access.rbacRoleDelete!}>
              <Button
                type={'primary'}
                size={'small'}
                danger={true}
                onClick={async () => {
                  roleDelete(record.id!);
                }}
              >
                删除
              </Button>
            </Access>
          </Space>
        );
      },
    },
  ];
  return (
    <PageContainer>
      <ProTable<RoleListItem>
        headerTitle={'角色列表'}
        actionRef={actionRef}
        rowKey="id"
        request={async (params = {}) => {
          const response = await roleIndex(params);
          return {
            success: response.code === 200,
            data: response.data,
          };
        }}
        columns={columns}
        pagination={false}
        search={false}
        toolBarRender={() => [
          <Access accessible={access.rbacRoleCreate!}>
            <Button
              type="primary"
              key="primary"
              onClick={() => {
                handleModalVisible(true);
              }}
            >
              <PlusOutlined /> 新建
            </Button>
          </Access>,
        ]}
      />
      {createModalVisible && (
        <ModalForm
          title="添加角色"
          width="450px"
          visible={createModalVisible}
          onVisibleChange={handleModalVisible}
          modalProps={{
            centered: true,
          }}
          onFinish={async (value) => {
            const response = await addRole(value as RoleListItem);
            const { code, msg } = response;
            if (code !== 200) {
              message.error(msg);
              return;
            }
            handleModalVisible(false);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }}
        >
          <ProFormText
            label={'角色名称'}
            placeholder={'请填写角色名称'}
            rules={[
              {
                required: true,
                message: '请填写角色名名称',
              },
            ]}
            name="name"
          />
          <ProFormText
            label={'角色编码'}
            placeholder={'请填写角色编码'}
            rules={[
              {
                required: true,
                message: '请填写角色名编码',
              },
            ]}
            name="code"
          />
          <ProFormTextArea
            label={'备注'}
            placeholder={'请输入备注'}
            rules={[
              {
                required: true,
                message: '请填写角色备注',
              },
            ]}
            name="remark"
          />
        </ModalForm>
      )}
      {updateModalVisible && (
        <ModalForm
          title="编辑角色"
          width="450px"
          visible={updateModalVisible}
          onVisibleChange={handleUpdateModalVisible}
          modalProps={{
            centered: true,
          }}
          onFinish={async (value) => {
            // payload 是新的值
            const payload: RoleListItem = {
              id: current.id,
              ...value, // 包含 value 全部key:value
            };
            // 更新新值
            const response = await updateRole(payload);
            const { code, msg } = response;
            if (code !== 200) {
              message.error(msg);
              return;
            }
            handleUpdateModalVisible(false);
            handleCurrent({}); // 当前对象清空
            if (actionRef.current) {
              actionRef.current.reload(); // 当前对象重新加载
            }
          }}
        >
          <ProFormText
            initialValue={current.name}
            label={'角色名称'}
            placeholder={'请填写角色名称'}
            rules={[
              {
                required: true,
                message: '请填写角色名名称',
              },
            ]}
            name="name"
          />
          <ProFormText
            initialValue={current.code}
            label={'角色编码'}
            placeholder={'请填写角色编码'}
            rules={[
              {
                required: true,
                message: '请填写角色名编码',
              },
            ]}
            name="code"
          />
          <ProFormTextArea
            initialValue={current.remark}
            label={'备注'}
            placeholder={'请输入备注'}
            rules={[
              {
                required: true,
                message: '请填写角色备注',
              },
            ]}
            name="remark"
          />
        </ModalForm>
      )}
    </PageContainer>
  );
};
export default RoleList;
