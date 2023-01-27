package com.jrasp.admin.server.kafka;

import com.alibaba.fastjson.JSONObject;
import com.jrasp.admin.common.vo.PageResult;
import com.jrasp.admin.server.pojo.AttackInfo;
import com.jrasp.admin.server.pojo.JavaProcessInfo;
import com.jrasp.admin.server.pojo.RaspHost;
import com.jrasp.admin.server.service.JavaProcessInfoService;
import com.jrasp.admin.server.service.RaspAttackInfoService;
import com.jrasp.admin.server.service.RaspHostService;
import io.krakens.grok.api.Grok;
import io.krakens.grok.api.GrokCompiler;
import io.krakens.grok.api.Match;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;
import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;
import org.springframework.transaction.annotation.Isolation;
import org.springframework.transaction.annotation.Transactional;

import java.text.SimpleDateFormat;
import java.util.*;


@Component
@Slf4j
public class KafkaConsumerService {

    private static final String pattern = "%{TIMESTAMP_ISO8601:time}\\s*%{LOGLEVEL:level}\\s*\\[%{DATA:thread}\\]\\s*(?<api>([\\S+]*))\\s*%{GREEDYDATA:message}";

    public static Grok grok = null;

    static {
        /* Create a new grokCompiler instance */
        GrokCompiler grokCompiler = GrokCompiler.newInstance();
        grokCompiler.registerDefaultPatterns();

        /* Grok pattern to compile, here httpd logs */
        grok = grokCompiler.compile(pattern);
    }

    @Autowired
    private RaspHostService raspHostService;

    @Autowired
    private JavaProcessInfoService javaProcessInfoService;

    @Autowired
    private RaspAttackInfoService raspAttackInfoService;

    @KafkaListener(topics = "jrasp-daemon", groupId = "jrasp")
    public void listenerDaemon(ConsumerRecord<?, ?> record) {
        log.debug("topic is {}, offset is {}, value is {}", record.topic(), record.offset(), record.value());
        String body = (String) record.value();
        DaemonMessage daemonMessage = JSONObject.parseObject(body, DaemonMessage.class);
        // error 日志打印
        if (daemonMessage != null && "ERROR".equals(daemonMessage.getLevel())) {
            log.error("jrasp-daemon err:" + body);
        }
        switch (Objects.requireNonNull(daemonMessage).getLogId()) {
            case LogConstant.DAEMON_STARTUP_LOGID:
                handleStartupLog(daemonMessage);
                break;
            case LogConstant.HOST_ENV_LOGID:
                handleHostEnvLog(daemonMessage);
                break;
            case LogConstant.HEART_BEAT_LOGID:
                handleHeartbeatLog(daemonMessage);
                break;
            case LogConstant.AGENT_SUCCESS_INIT:
                handleAgentInitAndUnloadLog(daemonMessage);
                break;
            case LogConstant.AGENT_SUCCESS_UNLOAD:
                handleAgentInitAndUnloadLog(daemonMessage);
                break;
            case LogConstant.JAVA_PROCESS_STARTUP:
                handleFindJavaProcessLog(daemonMessage);
                break;
            case LogConstant.JAVA_PROCESS_SHUTDOWN:
                handleRemoveJavaProcessLog(daemonMessage);
                break;
            case LogConstant.NACOS_INIT_INFO:
                handleNacosInitLog(daemonMessage);
                break;
            case LogConstant.Agent_CONFIG_UPDATE:
                handleAgentConfigUpdateLog(daemonMessage);
                break;
            case LogConstant.CONFIG_ID:
                handleUpdateConfigId(daemonMessage);
                break;
            default:
        }
    }

    @KafkaListener(topics = {"jrasp-agent", "jrasp-module"}, groupId = "jrasp")
    public void listenerAgent(ConsumerRecord<?, ?> record) {
        // TODO 暂时不处理这部分日志
    }

    @KafkaListener(topics = {"jrasp-attack"}, groupId = "jrasp")
    public void listenerAttack(ConsumerRecord<?, ?> record) {
        log.debug("topic is {}, offset is {}, value is {}", record.topic(), record.offset(), record.value());
        String body = (String) record.value();
        JSONObject jsonObject = JSONObject.parseObject(body);
        JSONObject hostObject = (JSONObject) jsonObject.get("host");
        String hostName = (String) hostObject.get("name");
        String message = (String) jsonObject.get("message");

        Match gm = grok.match(message);
        Map<String, Object> capture = gm.capture();
        // String logTime = (String) capture.get("time");
        // String threadName = (String) capture.get("thread");
        String level = (String) capture.get("level");
        if ("ERROR".equals(level) || "ERR".equals(level)) {
            log.error("jrasp-attack err:" + body);
        }
        String message2 = (String) capture.get("message");
        AgentMessage agentMessage = new AgentMessage();
        agentMessage.setLevel(level);
        agentMessage.setHostName(hostName);
        agentMessage.setMessage(message2);
        handleAttackLog(agentMessage);
    }

    private void handleStartupLog(DaemonMessage message) {
        RaspHost hostEntity = new RaspHost(message.getHostName(), message.getIp(), message.getTs());
        // detail字段转成json map
        JSONObject jsonObject = JSONObject.parseObject(message.getDetail());
        // 提取json map中的agentMode的值
        hostEntity.setAgentMode(jsonObject.getString("agentMode"));
        raspHostService.updateOrInsert(hostEntity);
        // 理论上kafka消息都不丢的情况，不需要定时补偿
    }

    @Transactional(isolation = Isolation.READ_COMMITTED, rollbackFor = Exception.class)
    void handleHeartbeatLog(DaemonMessage message) {
        String hostName = message.getHostName();
        RaspHost raspHost = new RaspHost(hostName, message.getTs());
        raspHost.setAgentInfo(message.getDetail());
        raspHost.setHeartbeatTime(message.getTs());
        raspHostService.updateOrInsert(raspHost);
        // 同时更新 java info 信息，对账信息
        PageResult<JavaProcessInfo> result = javaProcessInfoService.getIndex(hostName, 0L, 100L, null);
        List<JavaProcessInfo> list = result.getList();
        if (list == null) {
            return;
        }
        String detail = message.getDetail();
        Map map = JSONObject.parseObject(detail, Map.class);
        for (JavaProcessInfo javaProcessInfo : list) {
            if (map.size() == 0 || !map.containsKey(String.valueOf(javaProcessInfo.getPid()))) {
                // 需要删除的进程信息
                javaProcessInfoService.removeByHostNameAndPid(hostName, javaProcessInfo.getPid());
            }
        }
    }

    public void handleFindJavaProcessLog(DaemonMessage message) {
        JSONObject messageObject = JSONObject.parseObject(message.getDetail());
        int pid = messageObject.getInteger("javaPid");
        String startTime = messageObject.getString("startTime");
        String status = messageObject.getString("injectedStatus");
        String cmdLines = messageObject.getString("cmdLines");
        JavaProcessInfo javaInfo = new JavaProcessInfo(message.getHostName(), cmdLines, startTime, pid, status);
        javaProcessInfoService.updateOrInsert(javaInfo);
    }

    private void handleHostEnvLog(DaemonMessage message) {
        // 主机相关的信息，安装之后不发生变化，相对固定
        RaspHost hostEntity = new RaspHost(message.getHostName(), message.getIp(), message.getTs());
        JSONObject jsonObject = JSONObject.parseObject(message.getDetail());
        String installDir = jsonObject.getString("installDir");
        String version = jsonObject.getString("version");
        // 字段更新，修复问题
        String exeFileHash = jsonObject.getString("binFileHash");
        String osType = jsonObject.getString("osType");
        int totalMem = jsonObject.getInteger("totalMem");
        int cpuCounts = jsonObject.getInteger("cpuCounts");
        // 前端换成"剩余空间"
        int freeDisk = jsonObject.getInteger("freeDisk");
        // 新增 buildDateTime、buildGitBranch、buildGitCommit
        String buildDateTime = jsonObject.getString("buildDateTime");
        String buildGitBranch = jsonObject.getString("buildGitBranch");
        String buildGitCommit = jsonObject.getString("buildGitCommit");
        hostEntity.setBuildDateTime(buildDateTime);
        hostEntity.setBuildGitBranch(buildGitBranch);
        hostEntity.setBuildGitCommit(buildGitCommit);

        hostEntity.setInstallDir(installDir);
        hostEntity.setVersion(version);
        hostEntity.setExeFileHash(exeFileHash);
        hostEntity.setOsType(osType);
        hostEntity.setTotalMem(totalMem);
        hostEntity.setCpuCounts(cpuCounts);
        hostEntity.setFreeDisk(freeDisk);
        raspHostService.updateOrInsert(hostEntity);
    }

    public void handleRemoveJavaProcessLog(DaemonMessage message) {
        String pidStr = message.getDetail();
        javaProcessInfoService.removeByHostNameAndPid(message.getHostName(), Integer.parseInt(pidStr));
    }

    /**
     * 处理agent初始化成功/退出消息
     */
    private void handleAgentInitAndUnloadLog(DaemonMessage daemonMessage) {
        // agent 退出会上报一条日志，用来更新注入状态
        JSONObject jsonObject = JSONObject.parseObject(daemonMessage.getDetail());
        String hostName = daemonMessage.getHostName();
        Integer pid = jsonObject.getInteger("pid");
        String startTime = jsonObject.getString("startTime");
        String status = jsonObject.getString("status");
        JavaProcessInfo javaProcessInfo = new JavaProcessInfo(hostName, pid, startTime, status);
        javaProcessInfoService.updateStatus(javaProcessInfo);
    }

    /**
     * 处理攻击消息
     */
    private void handleAttackLog(AgentMessage agentMessage) {
        if (agentMessage == null) {
            return;
        }

        AttackInfo attackInfoEntity = new AttackInfo(agentMessage.getHostName());
        String message = agentMessage.getMessage();
        if ("".equals(message) || message.length() == 0) {
            return;
        }

        JSONObject attackInfoMap = JSONObject.parseObject(message);
        JSONObject context = JSONObject.parseObject(attackInfoMap.getString("context"));
        if (context != null && !context.isEmpty()) {
            attackInfoEntity.setRequestProtocol(context.getString("protocol"));
            attackInfoEntity.setHttpMethod(context.getString("method"));
            attackInfoEntity.setRemoteIp(context.getString("remoteHost"));
            // 上报字段修改
            attackInfoEntity.setRequestURI(context.getString("requestURL"));
            attackInfoEntity.setLocalIp(context.getString("localAddr"));
            String parameters = context.getString("parameters");
            attackInfoEntity.setRequestParameters(parameters == null ? "null" : parameters);
            // attackInfoEntity.setCookies(httpInfo.getString("cookies"));
            String header = context.getString("header");
            if (header != null) {
                String[] headerArray = header.split(";");
                String headerArrayJson = JSONObject.toJSONString(headerArray);
                attackInfoEntity.setHeader(headerArrayJson);
            }
            // 新增body信息
            String body = context.getString("body");
            if (body != null && !"".equals(body) && !"null".equals(body)) {
                String[] byteStringArray = body.split(",");
                byte[] byteBody = new byte[byteStringArray.length];
                for (int i = 0; i < byteStringArray.length; i++) {
                    byteBody[i] = Byte.parseByte(byteStringArray[i]);
                }
                String bodyStr = new String(byteBody);
                attackInfoEntity.setBody(bodyStr);
            }
        }
        attackInfoEntity.setIsBlocked(attackInfoMap.getBoolean("isBlocked"));
        attackInfoEntity.setLevel(attackInfoMap.getInteger("level"));
        attackInfoEntity.setAttackType(attackInfoMap.getString("attackType"));
        attackInfoEntity.setCheckType(attackInfoMap.getString("algorithm"));
        String stackTrace = attackInfoMap.getString("stackTrace");
        if (stackTrace != null) {
            String[] stackTraceArray = stackTrace.split(",");
            String stackTraceArrayJson = JSONObject.toJSONString(stackTraceArray);
            attackInfoEntity.setStackTrace(stackTraceArrayJson);
        }
        Long attackTime = attackInfoMap.getLong("attackTime");
        SimpleDateFormat df = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
        String attackTimeStr = df.format(new Date(attackTime));
        attackInfoEntity.setAttackTime(attackTimeStr);
        String attackParameters = attackInfoMap.getString("payload");
        attackInfoEntity.setAttackParameters(attackParameters);
        try {
            raspAttackInfoService.addNewAttackInfo(attackInfoEntity);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    /**
     * 处理依赖消息
     */
    private void handleDenpendencyLog(AgentMessage agentMessage) {
        if (!StringUtils.isEmpty(agentMessage.getPid())) {
            String hostName = agentMessage.getHostName();
            int javaPid = Integer.parseInt(agentMessage.getPid());
            //javaProcessInfoService.updateDependency(hostName, javaPid, agentMessage.getMessage());
        }
    }

    /**
     * 处理nacos状态消息
     */
    private void handleNacosInitLog(DaemonMessage daemonMessage) {
        if (!StringUtils.isEmpty(daemonMessage.getHostName())) {
            String hostName = daemonMessage.getHostName();
            String detail = daemonMessage.getDetail();
            raspHostService.updateHostNacosInfo(hostName, detail);
        }
    }

    /**
     * agent 配置更新状态
     */
    private void handleAgentConfigUpdateLog(DaemonMessage daemonMessage) {
        if (!StringUtils.isEmpty(daemonMessage.getHostName())) {
            String hostName = daemonMessage.getHostName();
            String agentConfigUpdateTime = daemonMessage.getTs();
            raspHostService.updateAgentConfigUpdateTime(hostName, agentConfigUpdateTime);
        }
    }

    /**
     * 配置更新ID
     */
    private void handleUpdateConfigId(DaemonMessage daemonMessage) {
        if (!StringUtils.isEmpty(daemonMessage.getHostName())) {
            String hostName = daemonMessage.getHostName();
            String detail = daemonMessage.getDetail();
            JSONObject configIdObject = JSONObject.parseObject(detail);
            Integer configId = configIdObject.getInteger("configId");
            raspHostService.updateHostConfigId(hostName, configId);
        }
    }


}
