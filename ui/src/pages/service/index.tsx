import React, {useEffect, useRef, useState} from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import ProTable, {ActionType, ProColumns} from '@ant-design/pro-table';
import {Button, message, Modal, Space, Form, Tag} from 'antd';
import {Access, useAccess} from '@@/plugin-access/access';
import {PlusOutlined} from '@ant-design/icons';
import ProForm, {
  ModalForm,
  ProFormText,
  ProFormRadio,
} from '@ant-design/pro-form';

import type {TableListItem, AddServiceInfo} from '@/pages/service/data';
import {queryService, deleteService, addService, updateService} from '@/pages/service/service';

const RaspConfigList: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const [createModalVisible, handleCreateModalVisible] = useState<boolean>(false);
  const [updateModalVisible, handleUpdateModalVisible] = useState<boolean>(false);

  const [row, handleRow] = useState<Partial<TableListItem>>({});
  const access: API.UserAccessItem = useAccess();
  const handleDeleteService = (id: number) => {
    Modal.confirm({
      title: '删除服务',
      content: '您确定要删除该服务吗？删除后不能恢复！',
      centered: true,
      okType: 'danger',
      onOk: async () => {
        const response = await deleteService(id);
        const {code, msg} = response;
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

  const columns: ProColumns<TableListItem>[] = [
    {
      title: 'id',
      dataIndex: 'id',
      hideInForm: true,
      hideInTable: true,
      search: false,
    },
    {
      title: '服务名称',
      dataIndex: 'service_name',
      tooltip:'支持模糊搜索',
      copyable: true,
    },
    {
      title: '服务描述',
      dataIndex: 'description',
      search: false,
    },
    {
      title: '服务等级',
      dataIndex: 'level',
      search: false,
      render(dataIndex) {
        if (dataIndex == true) {
          return (
            <Space>
              <Tag color="green">核心服务</Tag>
            </Space>
          );
        } else {
          return (
            <Space>
              <Tag color="blue">非核心服务</Tag>
            </Space>
          );
        }
      },
    },
    {
      title: '节点数量',
      dataIndex: 'node_number',
      search: false,
    },
    {
      title: '低峰期时间',
      tooltip:'0～23小时',
      dataIndex: 'low_period',
    },
    {
      title: '开发语言',
      dataIndex: 'language',
      search: false,
    },
    {
      title: '负责人',
      dataIndex: 'owners',
      search: false,
    },
    {
      title: '更新时间',
      search: false,
      hideInForm: true,
      dataIndex: 'update_time',
    },
    {
      title: '创建时间',
      search: false,
      hideInForm: true,
      dataIndex: 'create_time',
    },
    {
      title: '操作',
      hideInForm: true,
      search: false,
      render: (_, record) => {
        return (
          <Space>
            <Access accessible={access.serviceDelete!}>
              <a
                type={'primary'}
                onClick={async () => {
                  await handleDeleteService(record.id!);
                }}
              >
                删除
              </a>
            </Access>
            <Access accessible={access.serviceUpdate!}>
              <a
                type={'primary'}
                onClick={async () => {
                  handleRow(record);
                  handleUpdateModalVisible(true);
                }}
              >
                更新
              </a>
            </Access>
          </Space>
        );
      },
    },
  ];
  const [form] = Form.useForm();
  useEffect(() => {
    form.setFieldsValue({...row});
  }, []);
  return (
    <PageContainer>
      <ProTable
        headerTitle={'服务列表'}
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        columns={columns}
        toolBarRender={() => [
          <Button type="primary" onClick={() => handleCreateModalVisible(true)}>
            <PlusOutlined/> 新建
          </Button>,
        ]}
        request={async (params = {}) => {
          const response = await queryService(params);
          return {
            success: response.code === 200,
            data: response!.data!.list,
            page: response!.data!.page_num,
            total: response!.data!.total,
          };
        }}
      />
      {createModalVisible && (
        <ModalForm
          title="新增服务"
          visible={createModalVisible}
          onVisibleChange={handleCreateModalVisible}
          modalProps={{
            centered: true,
          }}
          initialValues={{
            language: 'java',
            level: 'false',
          }}
          layout={'horizontal'}
          onFinish={async (value) => {
            const response = await addService(value as AddServiceInfo);
            const {code, msg} = response;
            if (code !== 200) {
              message.error(msg);
              return;
            }
            message.success(msg);
            handleCreateModalVisible(false);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }}
        >
          <ProForm.Group>
            <ProFormText
              width="md"
              name="service_name"
              label="服务名称"
              tooltip="最长为 64 位"
              placeholder="请输入服务名称"
              rules={[
                {
                  required: true,
                  message: '请输入服务名称',
                },
              ]}
            />
            <ProFormText
              width="md"
              name="description"
              label="服务描述"
              tooltip="最长为 128 位"
              placeholder="请输入服务描述"
              rules={[
                {
                  required: true,
                  message: '请输入服务描述',
                },
              ]}
            />
          </ProForm.Group>
          <ProFormRadio.Group
            name="language"
            label="开发语言"
            options={[
              {
                label: 'java',
                value: 'java',
              },
              {
                label: 'go',
                value: 'go',
              },
              {
                label: 'python',
                value: 'python',
              },
              {
                label: 'js',
                value: 'js',
              },
              {
                label: 'c++',
                value: 'c++',
              },
              {
                label: '其他',
                value: '其他',
              },
            ]}
            rules={[
              {
                required: true,
                message: '请选择一个',
              },
            ]}
          />
          <ProFormRadio.Group
            name="level"
            label="服务等级"
            radioType="button"
            rules={[
              {
                required: true
              },
            ]}
            options={[
              {
                label: '非核心服务',
                value: 'false',
              },
              {
                label: '核心服务',
                value: 'true',
              }
            ]}
          />
        </ModalForm>
      )}
      {updateModalVisible && (
        <ModalForm
          title="更新服务"
          visible={updateModalVisible}
          onVisibleChange={handleUpdateModalVisible}
          modalProps={{
            centered: true,
          }}
          initialValues={{
            language: 'java',
            level: 'false',
          }}
          layout={'horizontal'}
          onFinish={async (value) => {
            const response = await updateService(value as AddServiceInfo);
            const {code, msg} = response;
            if (code !== 200) {
              message.error(msg);
              return;
            }
            message.success(msg);
            handleUpdateModalVisible(false);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }}
        >
          <ProFormText
            name={'id'}
            width="sm"
            initialValue={row.id}
            label={'服务ID'}
            readonly={true}
          />
          <ProFormText
            name={'service_name'}
            width="sm"
            initialValue={row.service_name}
            label={'服务名称'}
          />
          <ProFormText
            name={'low_period'}
            width="sm"
            initialValue={row.low_period}
            label={'低峰时间'}
          />
          <ProFormRadio.Group
            name="level"
            width="sm"
            label="服务等级"
            radioType="button"
            initialValue={row.level}
            options={[
              {
                label: '非核心服务',
                value: 'false',
              },
              {
                label: '核心服务',
                value: 'true',
              }
            ]}
          />

          <ProFormText
            name={'owners'}
            width="sm"
            initialValue={row.owners}
            label={'管理人员'}
          />
        </ModalForm>
      )}
    </PageContainer>
  );
};
export default RaspConfigList;
