import React, { useEffect, useRef, useState } from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import ProTable, { ActionType, ProColumns } from '@ant-design/pro-table';
import { message, Modal, Space, Form, Tag, Table } from 'antd';
import type { TableListItem } from '@/pages/attack/data';
import { deleteAttack, batchDeleteAttack, markAttack, queryAttack } from '@/pages/attack/service';
import { history } from '@@/core/history';

const RaspAttackInfoList: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const [row] = useState<Partial<TableListItem>>({});
  const handleDeleteAttack = (id: number) => {
    Modal.confirm({
      title: '删除日志',
      content: '您确定要删除该日志吗？删除后不能恢复！',
      centered: true,
      okType: 'danger',
      onOk: async () => {
        const response = await deleteAttack(id);
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

  const handleBatchDeleteAttack = (ids: (number | string)[]) => {
    Modal.confirm({
      title: '批量删除日志',
      content: '您确定要批量删除日志吗？删除后不能恢复！',
      centered: true,
      okType: 'danger',
      onOk: async () => {
        const response = await batchDeleteAttack(ids);
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
      hideInTable: true,
      search: false,
    },
    {
      title: '攻击时间',
      search: false,
      valueType: 'dateTime',
      dataIndex: 'attack_time',
    },
    {
      title: '实例名称',
      dataIndex: 'host_name',
      copyable: true,
    },
    {
      title: '攻击类型',
      dataIndex: 'attack_type',
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
      title: '攻击ip',
      search: false,
      dataIndex: 'remote_ip',
    },
    {
      title: '风险等级',
      tooltip: '范围：0-100',
      dataIndex: 'level',
      search: false,
      render(dataIndex) {
        if (dataIndex! >= 90) {
          return (
            <Space>
              <Tag color="red">高危</Tag>
            </Space>
          );
        } else if (dataIndex! <= 50) {
          return (
            <Space>
              <Tag color="">低危</Tag>
            </Space>
          );
        } else {
          return (
            <Space>
              <Tag color="orange">中危</Tag>
            </Space>
          );
        }
      },
    },
    {
      title: '阻断状态',
      dataIndex: 'is_blocked',
      valueType: 'select',
      valueEnum: {
        false: {
          text: '放行',
          status: 'Success',
        },
        true: {
          text: '阻断',
          status: 'Error',
        },
      },
    },
    {
      title: '处理状态',
      dataIndex: 'handle_status',
      valueType: 'select',
      valueEnum: {
        2147483647: {
          text: '全部',
          status: 'Success',
        },
        0: {
          text: '未处理',
          status: 'Default',
        },
        1: {
          text: '处理中',
          status: 'Processing',
        },
        2: {
          text: '已处理',
          status: 'Success',
        },
        3: {
          text: '误报',
          status: 'Warning',
        },
        4: {
          text: '忽略',
          color: 'black',
        },
      },
    },
    {
      title: '操作',
      hideInForm: true,
      search: false,
      render: (_, record) => {
        return (
          <Space>
            <a
              type={'dashed'}
              onClick={async () => {
                await handleDeleteAttack(record.id!);
              }}
            >
              删除
            </a>
            <a
              type={'dashed'}
              onClick={async () => {
                // 如果是未处理状态，则标记为处理中的状态
                if (record.handle_status == 0) {
                  await markAttack(record.id, 1);
                }
                history.push(`/attack/detail/index`, { attackInfo: record });
              }}
            >
              详情
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
        headerTitle={'攻击日志'}
        actionRef={actionRef}
        rowKey="id"
        pagination={{
          showSizeChanger: true,
        }}
        search={{
          labelWidth: 120,
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
        tableAlertOptionRender={({ selectedRowKeys }) => {
          return (
            <Space size={16}>
              <a
                type={'dashed'}
                onClick={async () => {
                  handleBatchDeleteAttack(selectedRowKeys);
                }}
              >
                批量删除
              </a>
              <a type={'dashed'} onClick={async () => {}}>
                批量忽略
              </a>
            </Space>
          );
        }}
        columns={columns}
        request={async (params = {}) => {
          const response = await queryAttack(params);
          return {
            success: response.code === 200,
            data: response!.data!.list,
            page: response!.data!.page_num,
            total: response!.data!.total,
          };
        }}
      />
    </PageContainer>
  );
};
export default RaspAttackInfoList;
