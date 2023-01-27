package com.jrasp.admin.server.kafka;

public interface LogConstant {
    // daemon 1000~1999
    int DAEMON_STARTUP_LOGID = 1000;
    int HOST_ENV_LOGID = 1002;
    int HEART_BEAT_LOGID = 1011;
    int DEPENDENCY_JAR_LOGID = 1016;
    int AGENT_SUCCESS_UNLOAD = 1017;
    int JAVA_PROCESS_STARTUP = 1018;
    int JAVA_PROCESS_SHUTDOWN = 1019;
    int AGENT_SUCCESS_INIT = 1020;
    // nacos 初始化信息
    int NACOS_INIT_INFO = 1024;

    int Agent_CONFIG_UPDATE = 1025;

    int CONFIG_ID = 1030;

    // agent 2000~2999

    // module 3000~3999
    // 依赖日志
    int DENPENDENCY_MODULE_JAR_INFO_LOG_ID = 3011;
    // 攻击日志
    int SEC_MODULE_ATTACK_LOG_ID = 3012;
    // 检测耗时
    int SEC_MODULE_HOOK_TIME_LOG_ID = 3013;
}
