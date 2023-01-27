import React, { useEffect, useRef, useState } from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import ProTable, { ActionType, ProColumns } from '@ant-design/pro-table';
import { Button, message, Modal, Space, Form, Switch } from 'antd';
import { Access, useAccess } from '@@/plugin-access/access';
import { PlusOutlined } from '@ant-design/icons';
import ProForm, {
  ModalForm,
  ProFormText,
  ProFormRadio,
  ProFormCheckbox,
} from '@ant-design/pro-form';

import type { TableListItem } from '@/pages/config/data';
import { queryConfig, deleteConfig, addConfig, copyConfig } from '@/pages/config/service';
import { ConfigDetail } from '@/pages/config/data';
import { getModuleType } from '@/pages/module/service';
import { history } from '@@/core/history';

const RaspConfigList: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const [createModalVisible, handleCreateModalVisible] = useState<boolean>(false);
  const [blockConfigVisible, handleBlockConfigVisible] = useState<boolean>(false);
  const [row] = useState<Partial<TableListItem>>({});
  const access: API.UserAccessItem = useAccess();
  const handleDeleteConfig = (id: number) => {
    Modal.confirm({
      title: '删除配置',
      content: '您确定要删除该配置吗？删除后不能恢复！',
      centered: true,
      okType: 'danger',
      onOk: async () => {
        const response = await deleteConfig(id);
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

  const handleConfigCopy = (id: number) => {
    Modal.confirm({
      title: '复制配置',
      content: '您确定要复制该配置吗？',
      centered: true,
      okType: 'default',
      onOk: async () => {
        const response = await copyConfig(id);
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
      title: '配置名称',
      dataIndex: 'config_name',
    },
    {
      title: '描述',
      dataIndex: 'message',
      search: false,
    },
    {
      title: '配置状态',
      dataIndex: 'status',
      valueType: 'select',
      valueEnum: {
        2147483647: {
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
            <Access accessible={access.configDelete!}>
              <a
                type={'primary'}
                onClick={async () => {
                  await handleDeleteConfig(record.id!);
                }}
              >
                删除
              </a>
            </Access>
            <Access accessible={access.configCopy!}>
              <a
                type={'primary'}
                onClick={async () => {
                  await handleConfigCopy(record.id!);
                }}
              >
                复制
              </a>
            </Access>
            <a
              type={'dashed'}
              onClick={async () => {
                history.push(`/config/detail/index`, { config: record });
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
        headerTitle={'配置列表'}
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
          const response = await queryConfig(params);
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
          title="新增配置"
          visible={createModalVisible}
          onVisibleChange={handleCreateModalVisible}
          modalProps={{
            centered: true,
          }}
          width={'lg'}
          initialValues={{
            agent_mode: 'dynamic',
          }}
          onFinish={async (value) => {
            const response = await addConfig(value as ConfigDetail);
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
              name="config_name"
              label="配置名称"
              tooltip="最长为 24 位"
              placeholder="请输入配置名称"
              rules={[
                {
                  required: true,
                  message: '请输入配置名称',
                },
                {
                  pattern: /^[0-9a-zA-Z_]{4,24}$/,
                  message: '4~24位字母数字下划线',
                },
              ]}
            />
            <ProFormRadio.Group
              name="agent_mode"
              label="Agent接入模式"
              tooltip="静态接入需要修改jvm参数并重启进程"
              options={[
                {
                  label: '动态接入',
                  value: 'dynamic',
                },
                {
                  label: '静态接入',
                  value: 'static',
                },
                {
                  label: '关闭',
                  value: 'disable',
                },
              ]}
              rules={[
                {
                  required: true,
                  message: '请至少选择一个',
                },
              ]}
            />
          </ProForm.Group>

          <ProFormCheckbox.Group
            name="modules"
            label="模块列表"
            rules={[
              {
                required: true,
                message: '请至少选择一个',
              },
            ]}
            request={async () => {
              const response = await getModuleType();
              return response!.data!;
            }}
          />
          <ProForm.Group>
            <Switch
              defaultChecked={false}
              checkedChildren={'高级配置'}
              unCheckedChildren={'高级配置'}
              onChange={async (checked) => {
                const status = checked ? 1 : 2;
                if (status == 1) {
                  handleBlockConfigVisible(true);
                } else {
                  handleBlockConfigVisible(false);
                }
              }}
            />
          </ProForm.Group>
          <Access accessible={blockConfigVisible!}>
            <ProForm.Group>
              <ProFormText
                width="md"
                name="redirect_url"
                label="阻断反馈链接"
                placeholder="请输入阻断反馈链接"
                initialValue={'https://www.jrasp.com/block.html'}
                rules={[
                  {
                    required: false,
                    message: '请输入阻断反馈链接',
                  },
                ]}
              />
              <ProFormText
                width="xs"
                name="block_status_code"
                label="阻断状态码"
                placeholder="请输入阻断状态码"
                initialValue={302}
                rules={[
                  {
                    required: false,
                    message: '请输入阻断状态码',
                  },
                  {
                    pattern: /^3\d{2}$/,
                    message: '不合法的状态码！(300~399)',
                  },
                ]}
              />
            </ProForm.Group>
            <ProFormText
              name={'json_block_content'}
              label={'json格式的阻断文本'}
              width="xl"
              placeholder={'请输入json格式的阻断内容'}
              initialValue={
                '{"error":true, "reason": "Request blocked by JRASP (https://www.jrasp.com)"}'
              }
              rules={[
                {
                  required: false,
                  message: '请输入json格式的阻断内容',
                },
              ]}
            />
            <ProFormText
              name={'xml_block_content'}
              label={'xml格式的阻断文本'}
              width="xl"
              placeholder={'请输入xml格式的阻断内容'}
              initialValue={
                '<?xml version="1.0"?><doc><error>true</error><reason>Request blocked by JRASP</reason></doc>'
              }
              rules={[
                {
                  required: false,
                  message: '请输入xml格式的阻断内容',
                },
              ]}
            />
            <ProFormText
              name={'html_block_content'}
              label={'html格式的阻断文本'}
              width="xl"
              placeholder={'请输入html格式的阻断文本'}
              initialValue={
                '</script><script>location.href="https://www.jrasp.com/block.html"</script>'
              }
              rules={[
                {
                  required: false,
                  message: '请输入html格式的阻断文本',
                },
              ]}
            />

            <ProForm.Group>
              <ProFormText
                width="lg"
                name="bin_file_url"
                label="守护进程下载链接"
                placeholder="请输入守护进程下载链接(空默认不更新)"
                rules={[
                  {
                    required: false,
                    message: '请输入配置名称(空默认不更新)',
                  },
                ]}
              />
            </ProForm.Group>
            <ProForm.Group>
              <ProFormText
                width="lg"
                name="bin_file_hash"
                label="守护进程文件hash"
                tooltip={'查看文件的hash: linux上执行md5sum命令,macos上执行md5命令'}
                placeholder="请输入守护进程文件hash(空默认不更新)"
                rules={[
                  {
                    required: false,
                    message: '请输入守护进程文件hash(空默认不更新)',
                  },
                ]}
              />
            </ProForm.Group>
          </Access>
          <ProFormText
            name={'message'}
            label={'备注'}
            width="xl"
            placeholder={'请输入备注信息'}
            rules={[
              {
                required: true,
                message: '请输入备注信息',
              },
              {
                pattern: /^[^\u4e00-\u9fa5]+$/,
                message: '输入英文备注',
              },
            ]}
          />
        </ModalForm>
      )}
    </PageContainer>
  );
};
export default RaspConfigList;
