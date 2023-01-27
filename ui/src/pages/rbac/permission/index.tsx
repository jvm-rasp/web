import React, { useRef, useState } from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import type { ActionType, ProColumns } from '@ant-design/pro-table';
import ProTable from '@ant-design/pro-table';
import { addPermission, permissionIndex, removePermission, updatePermission } from './service';
import { Button, Form, message, Modal, Space } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { ModalForm, ProFormRadio, ProFormText } from '@ant-design/pro-form';
import PermissionTreeSelect from '@/components/Permission/tree';
import { Access, useAccess } from '@@/plugin-access/access';

const PermissionList: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const [createModalVisible, handleCreateModalVisible] = useState<boolean>();
  const [updateModalVisible, handleUpdateModalVisible] = useState<boolean>();
  const [current, setCurrent] = useState<Partial<PermissionListItem>>({});
  const access: API.UserAccessItem = useAccess();
  const permissionDelete = (id: number) => {
    Modal.confirm({
      title: '删除权限',
      content: '您确定要删除该权限吗？删除后无法恢复！',
      centered: true,
      okType: 'danger',
      onOk: async () => {
        const response = await removePermission(id);
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
  const columns: ProColumns<PermissionListItem>[] = [
    {
      title: 'id',
      dataIndex: 'id',
    },
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
    {
      title: '操作',
      render: (_, record) => {
        return (
          <Space>
            <Access accessible={access.rbacPermissionCreate!}>
              <Button
                type={'primary'}
                size={'small'}
                onClick={() => {
                  handleUpdateModalVisible(true);
                  setCurrent(record);
                }}
              >
                编辑
              </Button>
            </Access>
            <Access accessible={access.rbacPermissionDelete!}>
              <Button
                type={'primary'}
                size={'small'}
                danger={true}
                onClick={async () => {
                  permissionDelete(record.id!);
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
      <ProTable<PermissionListItem>
        headerTitle={'权限列表'}
        actionRef={actionRef}
        rowKey={'id'}
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
        toolBarRender={() => [
          <Access accessible={access.rbacPermissionCreate!}>
            <Button
              type="primary"
              key="primary"
              onClick={() => {
                handleCreateModalVisible(true);
              }}
            >
              <PlusOutlined /> 新建
            </Button>
          </Access>,
        ]}
      />
      {createModalVisible && (
        <ModalForm
          title={'添加权限'}
          visible={createModalVisible}
          onVisibleChange={handleCreateModalVisible}
          modalProps={{
            centered: true,
          }}
          onFinish={async (value) => {
            const response = await addPermission(value as PermissionListItem);
            const { code, msg } = response;
            if (code !== 200) {
              message.error(msg);
              return;
            }
            handleCreateModalVisible(false);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }}
        >
          <ProFormText
            label={'权限名称'}
            placeholder={'请填写权限名称'}
            rules={[{ required: true, message: '请填写权限名称' }]}
            name={'name'}
          />
          <Form.Item label={'父级权限'} name={'parent_id'}>
            <PermissionTreeSelect />
          </Form.Item>
          <ProFormText
            label={'权限网址'}
            placeholder={'请填写权限网址'}
            rules={[{ required: true, message: '请填写权限网址' }]}
            name={'url'}
          />
          <ProFormText label={'权限图标'} placeholder={'请填写图标(iconfont)'} name={'icon'} />
          <ProFormRadio.Group
            name="hide_in_menu"
            label="导航菜单"
            initialValue={2}
            rules={[{ required: true, message: '请选择菜单是否显示' }]}
            options={[
              {
                label: '显示',
                value: 2,
              },
              {
                label: '隐藏',
                value: 1,
              },
            ]}
          />
          <ProFormRadio.Group
            name="hide_children_in_menu"
            label="隐藏子菜单"
            initialValue={2}
            rules={[{ required: true, message: '请选择是否隐藏子菜单' }]}
            options={[
              {
                label: '否',
                value: 2,
              },
              {
                label: '是',
                value: 1,
              },
            ]}
          />
          <ProFormText label={'优先级'} placeholder={'请填写权限优先级'} name={'priority'} />
        </ModalForm>
      )}
      {updateModalVisible && (
        <ModalForm
          title={'添加权限'}
          visible={updateModalVisible}
          onVisibleChange={handleUpdateModalVisible}
          modalProps={{
            centered: true,
          }}
          onFinish={async (value) => {
            const payload = {
              id: current.id,
              ...value,
            };
            const response = await updatePermission(payload as PermissionListItem);
            const { code, msg } = response;
            if (code !== 200) {
              message.error(msg);
              return;
            }
            handleUpdateModalVisible(false);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }}
        >
          <ProFormText
            label={'权限名称'}
            placeholder={'请填写权限名称'}
            initialValue={current.name}
            rules={[{ required: true, message: '请填写权限名称' }]}
            name={'name'}
          />
          <Form.Item label={'父级权限'} name={'parent_id'} initialValue={current.parent_id}>
            <PermissionTreeSelect />
          </Form.Item>
          <ProFormText
            initialValue={current.url}
            label={'权限网址'}
            placeholder={'请填写权限网址'}
            rules={[{ required: true, message: '请填写权限网址' }]}
            name={'url'}
          />
          <ProFormText
            initialValue={current.icon}
            label={'权限图标'}
            placeholder={'请填写图标(iconfont)'}
            name={'icon'}
          />
          <ProFormRadio.Group
            name="hide_in_menu"
            label="导航菜单"
            initialValue={current.hide_in_menu}
            rules={[{ required: true, message: '请选择菜单是否显示' }]}
            options={[
              {
                label: '显示',
                value: 2,
              },
              {
                label: '隐藏',
                value: 1,
              },
            ]}
          />
          <ProFormRadio.Group
            name="hide_children_in_menu"
            label="隐藏子菜单"
            initialValue={current.hide_children_in_menu}
            rules={[{ required: true, message: '请选择是否隐藏子菜单' }]}
            options={[
              {
                label: '否',
                value: 2,
              },
              {
                label: '是',
                value: 1,
              },
            ]}
          />
          <ProFormText
            label={'优先级'}
            initialValue={current.priority}
            placeholder={'请填写权限优先级'}
            name={'priority'}
          />
        </ModalForm>
      )}
    </PageContainer>
  );
};

export default PermissionList;
