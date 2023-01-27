import {useState, useEffect} from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import ProCard, {StatisticCard} from '@ant-design/pro-card';
import {fetchWeekData} from '@/pages/dashboard/service';
import RcResizeObserver from 'rc-resize-observer';
import {Pie} from '@ant-design/charts';

const {Statistic, Divider} = StatisticCard;

const DashBoard = () => {
  const [data, setData] = useState<any>({});
  const [responsive, setResponsive] = useState(false);

  const getfetchDashboard = async () => {
    //发送请求获取统计数据
    const res = await fetchWeekData();
    if (res.code === 200) {
      setData(res!.data);
      return;
    }
  };

  useEffect(() => {
    getfetchDashboard();
  }, []);

  return (
    <PageContainer>
      <RcResizeObserver
        key="resize-observer"
        onResize={(offset) => {
          setResponsive(offset.width < 596);
        }}
      >
        <ProCard
          title="最近7天攻击数据概览"
          extra={data.time}
          split={responsive ? 'horizontal' : 'vertical'}
          headerBordered
          bordered
        >
          <ProCard split="horizontal">
            <ProCard split="horizontal">
              <ProCard split="vertical">
                <StatisticCard.Group direction={responsive ? 'column' : 'row'}>
                  <StatisticCard
                    statistic={{
                      title: '总攻击次数',
                      value: data.all_attack_cnt,
                      description: (
                        <Statistic
                          title="对比7天前"
                          value={data.all_attack_diff_value}
                          trend={data.all_attack_trend}
                        />
                      ),
                    }}
                  />
                  <Divider type={responsive ? 'horizontal' : 'vertical'}/>
                  <StatisticCard
                    statistic={{
                      title: '高危漏洞',
                      value: data.week_high_attack_cnt,
                      description: (
                        <StatisticCard.Statistic
                          title="占比"
                          value={
                            data.all_attack_cnt === 0
                              ? '100%'
                              : ''.concat(
                                ((data.week_high_attack_cnt / data.all_attack_cnt) * 100).toFixed(
                                  2,
                                ),
                                '%',
                              )
                          }
                        />
                      ),
                    }}
                    chart={
                      <Pie
                        height={100}
                        width={100}
                        innerRadius={0.7}
                        renderer="svg"
                        data={[
                          {name: '高危漏洞', value: data.week_high_attack_cnt},
                          {
                            name: '',
                            value: data.all_attack_cnt - data.week_high_attack_cnt,
                          },
                        ]}
                        angleField="value"
                        colorField="name"
                        statistic={undefined}
                        legend={false}
                        tooltip={false}
                        color={['#f60505', '#f5f2f2']}
                      />
                    }
                    chartPlacement="left"
                  />
                  <StatisticCard
                    statistic={{
                      title: '拦截漏洞',
                      value: data.week_attack_block_cnt,
                      description: (
                        <StatisticCard.Statistic
                          title="占比"
                          value={
                            data.all_attack_cnt === 0
                              ? '100%'
                              : ''.concat(
                                (
                                  (data.week_attack_block_cnt / data.all_attack_cnt) *
                                  100
                                ).toFixed(2),
                                '%',
                              )
                          }
                        />
                      ),
                    }}
                    chart={
                      <Pie
                        height={100}
                        width={100}
                        innerRadius={0.7}
                        renderer="svg"
                        data={[
                          {name: '拦截漏洞', value: data.week_attack_block_cnt},
                          {
                            name: '',
                            value: data.all_attack_cnt - data.week_attack_block_cnt,
                          },
                        ]}
                        angleField="value"
                        colorField="name"
                        statistic={undefined}
                        legend={false}
                        tooltip={false}
                        color={['#6394f9', '#f5f2f2']}
                      />
                    }
                    chartPlacement="left"
                  />
                </StatisticCard.Group>
              </ProCard>
            </ProCard>
          </ProCard>
        </ProCard>
        <br/>
      </RcResizeObserver>
    </PageContainer>
  );
};

export default DashBoard;
