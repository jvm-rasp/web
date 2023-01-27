import React, { useEffect, useRef, useState } from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import ProTable, { ActionType, ProColumns } from '@ant-design/pro-table';
import { Button, message, Modal, Space, Form, Tag, Switch } from 'antd';
import { Access, useAccess } from '@@/plugin-access/access';
import ProCard from '@ant-design/pro-card';

import type { TableListItem } from '@/pages/module/data';
import {
  queryModule,
  deleteModule,
  addModule,
  updateModule,
  statusModule,
} from '@/pages/module/service';
import { PlusOutlined } from '@ant-design/icons';
import ProForm, {
  ModalForm,
  ProFormList,
  ProFormRadio,
  ProFormGroup,
  ProFormText,
  ProFormTextArea,
} from '@ant-design/pro-form';

const RaspdModuleList: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const [createModalVisible, handleCreateModalVisible] = useState<boolean>(false);
  const [updateModalVisible, handleUpdateModalVisible] = useState<boolean>(false);
  const [row, handleRow] = useState<Partial<TableListItem>>({});
  const access: API.UserAccessItem = useAccess();

  const handleDeleteConfig = (id: number) => {
    Modal.confirm({
      title: '删除插件',
      content: '您确定要删除该插件吗？删除后不能恢复！',
      centered: true,
      okType: 'danger',
      onOk: async () => {
        const response = await deleteModule(id);
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

  const columns: ProColumns<TableListItem>[] = [
    {
      title: 'id',
      dataIndex: 'id',
      hideInTable: true,
      search: false,
    },
    {
      title: '名称',
      dataIndex: 'name',
    },
    {
      title: '类型',
      dataIndex: 'type',
    },
    {
      title: '描述',
      ellipsis: true,
      copyable: true,
      dataIndex: 'message',
      search: false,
    },
    {
      title: '版本',
      dataIndex: 'version',
      tooltip: '版本号命名规范：x.y.z',
      search: false,
      render: (dataIndex) => {
        return (
          <Space>
            <Tag color="blue">{dataIndex}</Tag>
          </Space>
        );
      },
    },
    {
      title: '下载链接',
      dataIndex: 'url',
      ellipsis: true,
      hideInTable: true,
      search: false,
    },
    {
      title: '插件hash',
      dataIndex: 'hash',
      ellipsis: true,
      hideInTable: true,
      search: false,
    },
    {
      title: '目标中间件',
      dataIndex: 'middleware_name',
      ellipsis: true,
      search: false,
    },
    {
      title: '中间件版本',
      dataIndex: 'middleware_version',
      hideInTable: true,
      ellipsis: true,
      search: false,
    },
    {
      title: '模块状态',
      dataIndex: 'status',
      valueType: 'select',
      tooltip: '禁用后，新建策略无法选中',
      valueEnum: {
        0: {
          text: '全部',
          status: 'Success',
        },
        1: {
          text: '可用',
          status: 'Success',
        },
        2: {
          text: '禁用',
          status: 'Processing',
        },
      },
      render: (dom, record) => {
        return (
          <Access
            accessible={access.moduleStatus!}
            fallback={() => {
              return dom;
            }}
          >
            <Switch
              defaultChecked={record.status === 1}
              checkedChildren={'可用'}
              unCheckedChildren={'禁用'}
              onChange={async (checked) => {
                const status = checked ? 1 : 2;
                const response = await statusModule({ id: record.id!, status });
                const { code, msg } = response;
                if (code === 200) {
                  message.success(msg);
                  return;
                }
                message.error(msg);
              }}
            />
          </Access>
        );
      },
    },
    {
      title: '操作人',
      dataIndex: 'username',
      search: false,
    },
    {
      title: '更新时间',
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
            <a
              type={'primary'}
              onClick={async () => {
                await handleDeleteConfig(record.id!);
              }}
            >
              删除
            </a>
            <a
              type={'primary'}
              onClick={() => {
                handleRow(record);
                handleUpdateModalVisible(true);
              }}
            >
              编辑
            </a>
          </Space>
        );
      },
    },
  ];
  const [form] = Form.useForm();
  useEffect(() => {
    form.setFieldsValue({ ...row });
  }, []);
  return (
    <PageContainer>
      <ProTable
        headerTitle={'插件列表'}
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        columns={columns}
        toolBarRender={() => [
          <Button type="primary" onClick={() => handleCreateModalVisible(true)}>
            <PlusOutlined /> 新建
          </Button>,
        ]}
        request={async (params = {}) => {
          const response = await queryModule(params);
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
          title="新增防护模块"
          width={'lg'}
          visible={createModalVisible}
          onVisibleChange={handleCreateModalVisible}
          modalProps={{
            centered: true,
          }}
          initialValues={{
            type: 'hook',
          }}
          onFinish={async (value) => {
            const response = await addModule(value as TableListItem);
            const { code, msg } = response;
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
              width="sm"
              name="name"
              label="模块名称"
              placeholder="请输入模块名称"
              rules={[
                {
                  required: true,
                  message: '请输入模块名称',
                },
              ]}
            />
            <ProFormText
              width="sm"
              name="version"
              label="模块版本"
              tooltip="版本号命名规范：x.y.z"
              placeholder="请输入模块版本号"
              rules={[
                {
                  required: true,
                  message: '请输入模块版本号',
                },
              ]}
            />
            <ProFormRadio.Group
              name="type"
              width="sm"
              label="模块类型"
              tooltip="是否需要修改字节码"
              options={[
                {
                  label: 'hook模块',
                  value: 'hook',
                },
                {
                  label: 'algorithm模块',
                  value: 'algorithm',
                },
                {
                  label: '其他模块',
                  value: 'other',
                },
              ]}
            />
          </ProForm.Group>
          <ProForm.Group>
          <ProFormText
            name={'url'}
            width="lg"
            label={'下载链接'}
            tooltip="模块下载地址，如对象存储服务"
            placeholder={'请输入输模块信息'}
            rules={[
              {
                required: true,
                message: '请输入输模块信息',
              },
            ]}
          />
          <ProFormText
            name={'hash'}
            width="lg"
            label={'插件hash'}
            tooltip="jar包md5"
            placeholder={'请入输插件hash'}
            rules={[
              {
                required: true,
                message: '请入输插件hash',
              },
            ]}
          />
           </ProForm.Group>
          <ProFormList name="parameters" label="模块参数信息">
            <ProFormGroup key="group">
              <ProFormText
                rules={[
                  {
                    required: true,
                  },
                ]}
                name="key"
                label="参数名称"
              />
              <ProFormTextArea
                rules={[
                  {
                    required: true,
                  },
                ]}
                name="info"
                label="参数含义"
              />
              <ProFormTextArea
                rules={[
                  {
                    required: true,
                  },
                ]}
                width={'lg'}
                name="value"
                label="默认值"
              />
            </ProFormGroup>
          </ProFormList>
          <ProFormTextArea
            name={'message'}
            label={'描述'}
            placeholder={'请输模块信息'}
            rules={[
              {
                required: true,
                message: '请输模块信息',
              },
            ]}
          />
        </ModalForm>
      )}
      {updateModalVisible && (
        <ModalForm
          title="编辑模块信息"
          visible={updateModalVisible}
          onVisibleChange={handleUpdateModalVisible}
          width={'lg'}
          modalProps={{
            centered: true,
          }}
          onFinish={async (value: TableListItem) => {
            const payload: TableListItem = {
              id: row!.id!,
              ...value,
            };
            const response = await updateModule(payload);
            const { code, msg } = response;
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
          <ProForm.Group>
            <ProFormText
              width="sm"
              name="name"
              label="模块名称"
              placeholder="请输入模块名称"
              initialValue={row.name}
              rules={[
                {
                  required: true,
                  message: '请输入模块名称',
                },
              ]}
            />
            <ProFormText
              width="sm"
              name="version"
              initialValue={row.version}
              label="模块版本"
              placeholder="请输入模块版本号"
              tooltip="版本号命名规范：x.y.z"
              rules={[
                {
                  required: true,
                  message: '请输入模块版本号',
                },
              ]}
            />
            <ProFormRadio.Group
              name="type"
              width="sm"
              label="模块类型"
              tooltip="是否需要修改字节码"
              initialValue={row.type}
              options={[
                {
                  label: 'hook模块',
                  value: 'hook',
                },
                {
                  label: 'algorithm模块',
                  value: 'algorithm',
                },
                {
                  label: '其他模块',
                  value: 'other',
                },
              ]}
            />
          </ProForm.Group>
          <ProFormText
            name={'url'}
            initialValue={row.url}
            label={'下载链接'}
            tooltip="模块下载地址，如对象存储服务"
            placeholder={'请输入输模块信息'}
            rules={[
              {
                required: true,
                message: '请输入输模块信息',
              },
            ]}
          />
          <ProFormText
            name={'hash'}
            initialValue={row.hash}
            label={'插件hash'}
            tooltip="jar包md5"
            placeholder={'请输入插件hash'}
            rules={[
              {
                required: true,
                message: '请输入插件hash',
              },
            ]}
          />

          <ProFormList name="parameters" initialValue={row.parameters} label="模块参数信息">
            <ProCard headerBordered={true} hoverable={true} bordered={true}>
              <ProFormGroup key="group">
                <ProFormText
                  rules={[
                    {
                      required: true,
                    },
                  ]}
                  name="key"
                  width="md"
                  placeholder={'请输参数名称'}
                  label="参数名称"
                />
                <ProFormText
                  rules={[
                    {
                      required: true,
                    },
                  ]}
                  name="info"
                  width="md"
                  placeholder={'请输参数含义'}
                  label="参数含义"
                />
              </ProFormGroup>
              <ProFormTextArea
                rules={[
                  {
                    required: true,
                  },
                ]}
                width={'xl'}
                name="value"
                placeholder={'请输参数值'}
                tooltip={'数组类型参数每行一个'}
                label="默认值"
              />
            </ProCard>
          </ProFormList>
          <ProFormTextArea
            name={'message'}
            initialValue={row.message}
            label={'描述'}
            placeholder={'请输入模块信息'}
            rules={[
              {
                required: true,
                message: '请输入模块信息',
              },
            ]}
          />
        </ModalForm>
      )}
    </PageContainer>
  );
};
export default RaspdModuleList;
