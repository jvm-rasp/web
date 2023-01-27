import React, { useEffect, useRef, useState } from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import ProTable, { ActionType, ProColumns } from '@ant-design/pro-table';
import { message, Modal, Space, Form, Tag, Table } from 'antd';
import type { TableListItem } from '@/pages/host/data';
import { batchUpdateConfig, deleteHost, queryHost, updateConfig } from '@/pages/host/service';
import { Access, useAccess } from '@@/plugin-access/access';
import { history } from '@@/core/history';
import { ModalForm, ProFormSelect, ProFormText, ProFormTextArea } from '@ant-design/pro-form';
import { getSimpleConfigList } from '@/pages/config/service';
import { UpdateConfigDetail } from '@/pages/host/data';

const RaspHostList: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const [createBatchUpdateModalVisible, handleCreateBatchUpdateModalVisible] =
    useState<boolean>(false);
  const [createModalVisible, handleCreateModalVisible] = useState<boolean>(false);
  const [row, handleRow] = useState<Partial<TableListItem>>({});
  const [selectRow, handlSelectRow] = useState<Partial<TableListItem[]>>([]);

  const access: API.UserAccessItem = useAccess();
  const handleDeleteHost = (id: number) => {
    Modal.confirm({
      title: '删除实例',
      content: '您确定要删除该实例吗？删除后不能恢复！',
      centered: true,
      okType: 'danger',
      onOk: async () => {
        const response = await deleteHost(id);
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
      hideInForm: true,
      search: false,
    },
    {
      title: '实例名称',
      dataIndex: 'host_name',
      copyable: true,
    },
    {
      title: 'IP',
      dataIndex: 'ip',
      tip: '内网ip',
    },
    {
      title: 'RASP版本',
      dataIndex: 'version',
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
      title: '接入方式',
      dataIndex: 'agent_mode',
      valueType: 'select',
      filters: true,
      onFilter: true,
      valueEnum: {
        dynamic: {
          text: '动态',
          status: 'Success',
        },
        static: {
          text: '静态',
          status: 'Processing',
        },
        disable: {
          text: '关闭',
          status: 'Default',
        },
      },
    },
    {
      title: '配置id',
      dataIndex: 'config_id',
      search: false,
    },
    {
      title: '在线状态',
      dataIndex: 'online',
      tip: '5分钟上报一次',
      search: false,
      render: (dataIndex) => {
        if (dataIndex == true) {
          return (
            <Space>
              <Tag color="green">在线</Tag>
            </Space>
          );
        } else {
          return (
            <Space>
              <Tag color="grey">离线</Tag>
            </Space>
          );
        }
      },
    },
    {
      title: '配置生效时间',
      dataIndex: 'agent_config_update_time',
      valueType: 'dateTime',
      search: false,
    },
    {
      title: '注册时间',
      search: false,
      sorter: true,
      valueType: 'dateTime',
      dataIndex: 'create_time',
    },
    {
      title: '操作',
      hideInForm: true,
      search: false,
      render: (_, record) => {
        return (
          <Space>
            <Access accessible={access.hostDelete!}>
              <a
                type={'primary'}
                onClick={async () => {
                  await handleDeleteHost(record.id!);
                }}
              >
                删除
              </a>
            </Access>
            <Access accessible={access.hostDetail!}>
              <a
                type={'primary'}
                onClick={async () => {
                  handleRow(record);
                  handleCreateModalVisible(true);
                }}
              >
                更新策略
              </a>
            </Access>
            <Access accessible={access.hostDetail!}>
              <a
                type={'dashed'}
                onClick={async () => {
                  history.push(`/host/detail/index`, { host: record });
                }}
              >
                详情
              </a>
            </Access>
          </Space>
        );
      },
    },
  ];
  const [form] = Form.useForm();
  useEffect(() => {
    form.setFieldsValue({ ...row });
  }, []);
  // @ts-ignore
  return (
    <PageContainer>
      <ProTable
        headerTitle={'实例列表'}
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        pagination={{
          pageSize: 10,
          pageSizeOptions: [10, 20, 50, 100, 500],
          showSizeChanger: true,
        }}
        rowSelection={{
          selections: [Table.SELECTION_ALL, Table.SELECTION_INVERT],
        }}
        tableAlertRender={({ selectedRowKeys, onCleanSelected }) => (
          <Space size={24}>
            <span>
              已选 {selectedRowKeys.length} 项
              <a style={{ marginLeft: 8 }} onClick={onCleanSelected}>
                取消选择
              </a>
            </span>
          </Space>
        )}
        tableAlertOptionRender={({ selectedRows }) => {
          return (
            <Space size={16}>
              <a
                type={'dashed'}
                onClick={async () => {
                  handlSelectRow(selectedRows);
                  handleCreateBatchUpdateModalVisible(true);
                }}
              >
                批量更新配置
              </a>
            </Space>
          );
        }}
        columns={columns}
        request={async (params = {}) => {
          const response = await queryHost(params);
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
          title="更新策略"
          visible={createModalVisible}
          onVisibleChange={handleCreateModalVisible}
          modalProps={{
            centered: true,
          }}
          onFinish={async (value) => {
            const response = await updateConfig(value as UpdateConfigDetail);
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
          <ProFormText
            name={'hostName'}
            initialValue={row.host_name}
            label={'实例名称'}
            readonly={true}
          />
          <ProFormSelect
            name="configId"
            label="选择策略"
            rules={[
              {
                required: true,
                message: '请选择一个',
              },
            ]}
            request={async () => {
              const response = await getSimpleConfigList();
              return response!.data!;
            }}
          />
        </ModalForm>
      )}
      {createBatchUpdateModalVisible && (
        <ModalForm
          title="批量更新策略"
          visible={createBatchUpdateModalVisible}
          onVisibleChange={handleCreateBatchUpdateModalVisible}
          modalProps={{
            centered: true,
          }}
          onFinish={async (value) => {
            const response = await batchUpdateConfig(value as UpdateConfigDetail);
            const { code, msg } = response;
            if (code !== 200) {
              message.error(msg);
              return;
            }
            message.success(msg);

            handleCreateBatchUpdateModalVisible(false);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }}
        >
          <ProFormTextArea
            name={'hostNames'}
            // @ts-ignore
            initialValue={selectRow.map((v) => v.host_name).join('\n')}
            label={'实例名称列表'}
            colSize={10}
            readonly={false}
          />
          <ProFormSelect
            name="configId"
            label="选择策略"
            rules={[
              {
                required: true,
                message: '请选择一个',
              },
            ]}
            request={async () => {
              const response = await getSimpleConfigList();
              return response!.data!;
            }}
          />
        </ModalForm>
      )}
    </PageContainer>
  );
};
export default RaspHostList;
