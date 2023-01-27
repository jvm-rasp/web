type WeekDataResult = {
  code?: number;
  msg?: string;
  data?: {
    time?: string;
    all_attack_cnt?: number;
    all_attack_diff_value?: number;
    all_attack_trend?: string;
    week_high_attack_cnt?: number;
    week_high_attack_diff_value?: number;
    week_high_attack_trend?: string;
    week_attack_block_cnt?: number;
    week_attack_block_diff_value?: number;
    week_attack_block_trend?: string;
  };
};
