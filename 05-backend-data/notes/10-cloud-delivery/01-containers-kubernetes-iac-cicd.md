---
stage: intermediate
direction: backend-data
topic: containers-kubernetes-iac-cicd
---

# 容器、Kubernetes、Terraform 与 CI/CD：把服务变成可验证的交付物

这一篇把一个 Go HTTP 服务从源码交付到 Kubernetes 集群。重点不是记住命令，而是建立一条可追溯链路：提交确定源码，构建产生不可变镜像，部署清单声明运行状态，Terraform 声明基础设施，流水线保存每次验证和发布的证据。

## 能力边界与前置知识

容器解决进程运行环境的一致性；Kubernetes 持续把集群实际状态收敛到声明状态；Terraform 管理云资源的期望状态；CI 持续验证变更，CD 把已验证版本按受控策略交付。它们不能替代应用的授权、业务约束、备份或监控。

前置知识是 Linux 进程、HTTP 健康接口、Git、镜像仓库和配置文件。示例使用 `api` 服务监听 `8080`，提供 `/healthz` 与 `/readyz`。前者表示进程仍可处理健康检查；后者表示服务已完成启动并可接收业务流量。

```mermaid
flowchart LR
    A["Commit"] --> B["测试与静态检查"]
    B --> C["构建镜像"]
    C --> D["镜像扫描与推送"]
    D --> E["部署清单引用 digest"]
    E --> F["Kubernetes rollout"]
    F --> G["Smoke test 与指标观察"]
    H["Terraform plan/apply"] --> I["集群、网络、Registry"]
    I --> F
```

同一个 Git commit 可以构建多个镜像标签，但生产发布应记录不可变的镜像 digest，例如 `registry.example.com/acme/api@sha256:...`。tag 是可移动名称，digest 才精确指向内容。

## Docker 的对象、层与运行边界

Dockerfile 是构建说明；镜像是只读层与配置的集合；容器是镜像加上可写层和运行时隔离后的进程实例。每条会修改文件系统的 Dockerfile 指令通常会形成可复用构建层。改变某层的输入会使依赖它的后续层失去缓存，因此先复制依赖清单、下载依赖，再复制频繁变化的源代码，通常能缩短迭代构建时间。

| 对象 | 实际作用 | 不应承担的责任 |
| --- | --- | --- |
| Dockerfile | 定义构建输入和运行镜像 | 保存生产密码 |
| Image | 可分发的运行文件系统与元数据 | 保存持久业务数据 |
| Container | 一个受隔离约束的进程 | 作为数据库备份位置 |
| Volume | 容器外持久化数据卷 | 替代数据库复制与备份策略 |
| Network | 容器间连通与端口发布 | 业务访问控制的唯一防线 |
| Registry | 分发镜像和清单 | 代替发布审批与扫描 |

容器的 PID 1 接收停止信号。若在 shell 中用 `sh -c` 启动而没有 `exec`，信号可能先交给 shell，应用不能及时退出。应用应处理 `SIGTERM`：停止接收新请求、等待有限时间完成在途请求、关闭连接后退出。Kubernetes 的 `terminationGracePeriodSeconds` 是这段时间的上限，不是应用可以无限等待的承诺。

## 案例一：构建一个最小且可停止的 Go API 镜像

以下 Dockerfile 使用多阶段构建。第一阶段保留编译工具，第二阶段只复制运行所需二进制和 CA 证书。`COPY --from` 的边界使最终镜像不包含 Go 编译器、模块缓存和源码。

```dockerfile
# syntax=docker/dockerfile:1
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags='-s -w' -o /out/api ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/api /app/api
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/api"]
```

`USER nonroot` 让进程不以 root 身份运行；它不能修复应用漏洞，也不能替代 Kubernetes 的安全上下文。distroless 镜像没有 shell，排障应使用日志、临时调试容器或专门的调试镜像，而不是为方便而把 shell 工具永久带入生产镜像。

构建与本地验证：

```sh
docker build -t acme/api:dev .
docker run --rm --name api -p 8080:8080 acme/api:dev
curl -fsS http://127.0.0.1:8080/healthz
docker stop --time 20 api
```

成功条件是 `curl` 得到 2xx，`docker stop` 在 20 秒内退出，日志中没有被强制杀死的未完成请求。失败分支是应用只监听 `127.0.0.1`：容器内部可访问，但端口发布无法使宿主访问；服务应监听 `0.0.0.0:8080` 或未指定地址的 `:8080`。

镜像内的临时文件在容器删除后丢失。需要持久化的 PostgreSQL 数据可在本地 Compose 使用 named volume；生产中应使用受备份和恢复策略保护的托管数据库或 PVC。不要把用户上传文件写进 Deployment Pod 的文件系统，因为滚动更新或重新调度都会改变该文件系统。

## Compose：开发环境的一致启动，而非生产集群替代品

Compose 将多个服务、网络和卷定义在一个项目中，适合本地启动 API、数据库和依赖。服务名会成为默认网络内的 DNS 名称，因此 API 应连接 `postgres:5432`，不是宿主的 `localhost`。

```yaml
services:
  api:
    build: .
    ports: ["8080:8080"]
    environment:
      DATABASE_URL: postgres://app:dev-password@postgres:5432/app?sslmode=disable
    depends_on:
      postgres:
        condition: service_healthy
  postgres:
    image: postgres:17
    environment:
      POSTGRES_USER: app
      POSTGRES_PASSWORD: dev-password
      POSTGRES_DB: app
    volumes: ["postgres-data:/var/lib/postgresql/data"]
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U app -d app"]
      interval: 5s
      timeout: 3s
      retries: 10
volumes:
  postgres-data: {}
```

`depends_on` 的健康条件只编排启动顺序，不能让网络故障消失。应用仍要配置连接超时、有限重试，并在数据库不可用时明确报错或保持未就绪。示例密码仅限本地开发；Compose 文件若进入仓库，生产凭据应从受控 Secret 注入，不能写入 YAML。

## Kubernetes：控制器、工作负载与网络入口

Kubernetes API 接收期望状态。控制器持续比较期望和实际状态并创建、删除或替换对象。Pod 是可共同调度、共享网络命名空间和卷的一组容器；它不是长期稳定的机器。Deployment 管理无状态 Pod 的 ReplicaSet 和滚动更新；StatefulSet 为有序、稳定网络标识或稳定存储声明的工作负载服务，但不自动把任意数据库变成安全高可用数据库。

| 资源 | 何时使用 | 关键限制 |
| --- | --- | --- |
| Pod | 调试或控制器创建的最小单元 | 单独创建的 Pod 不会自动重建 |
| Deployment | 无状态 HTTP/API/worker | Pod 名称和本地磁盘不稳定 |
| StatefulSet | 每副本需要稳定身份或 PVC | 仍需应用级复制和恢复设计 |
| Service | 为后端 Pod 提供稳定虚拟地址 | 不承担 HTTP 路由规则 |
| Ingress | HTTP/HTTPS 路由声明 | 必须存在匹配的 Ingress controller |
| Job/CronJob | 一次性迁移或定时任务 | 幂等和并发策略由任务设计保证 |
| ConfigMap/Secret | 注入非敏感/敏感配置 | Secret 默认只是编码，不等于加密管理 |
| PV/PVC | 抽象存储供给与请求 | 访问模式、扩容和备份取决于 StorageClass |

Service 通过 label selector 找到可用 Endpoint。应用调用 `api.default.svc.cluster.local` 时，DNS 解析到 Service 的虚拟 IP 或 Endpoint；流量不需要、也不应该依赖某个 Pod IP。Ingress 只是一组规则，实际 TLS、负载均衡、限流能力由集群安装的控制器决定。

## 健康检查是流量与重启决策

startup probe 在成功前阻止 liveness 和 readiness 触发；它适合启动慢的进程。readiness probe 失败时，Pod 从对应 Service 的可用 Endpoint 中移除，但容器不会仅因未就绪而重启。liveness probe 连续失败会触发容器重启；把依赖数据库的短暂不可用直接判断为 liveness 失败，可能使所有应用实例同时重启并放大故障。

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 0
      maxSurge: 1
  selector:
    matchLabels: {app: api}
  template:
    metadata:
      labels: {app: api}
    spec:
      terminationGracePeriodSeconds: 30
      containers:
        - name: api
          image: registry.example.com/acme/api@sha256:REPLACE_WITH_REAL_DIGEST
          ports: [{containerPort: 8080}]
          resources:
            requests: {cpu: "100m", memory: "128Mi"}
            limits: {cpu: "500m", memory: "256Mi"}
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
            runAsNonRoot: true
            capabilities: {drop: ["ALL"]}
          startupProbe:
            httpGet: {path: /healthz, port: 8080}
            failureThreshold: 30
            periodSeconds: 2
          readinessProbe:
            httpGet: {path: /readyz, port: 8080}
            periodSeconds: 5
          livenessProbe:
            httpGet: {path: /healthz, port: 8080}
            periodSeconds: 10
---
apiVersion: v1
kind: Service
metadata: {name: api}
spec:
  selector: {app: api}
  ports: [{port: 80, targetPort: 8080}]
```

`requests` 参与调度：调度器以它判断节点是否还有可分配资源。`limits` 限制容器资源；内存超限可能触发 OOM kill，CPU 限制可能造成节流。数值必须来自压测、历史用量和 SLO 目标，不要复制固定的“标准配置”。

验证滚动升级：

```sh
kubectl apply --server-side -f k8s/api.yaml
kubectl rollout status deployment/api --timeout=5m
kubectl get pods -l app=api
kubectl get endpointslices -l kubernetes.io/service-name=api
kubectl rollout undo deployment/api
```

`rollout status` 成功只表示控制器认为副本已就绪；它不是端到端业务验收。发布后还要用真实认证请求检查关键路径、错误率、延迟和版本标签。回滚 Deployment 只回滚工作负载模板，不能自动回滚已经执行的破坏性数据库迁移。

## 配置、密钥、存储和批处理的边界

ConfigMap 适合公开配置，如功能开关默认值和公共域名；Secret 适合敏感值的传递，但 base64 不是加密。集群至少应限制 Secret 的 RBAC 读取权限，并在需要时使用外部密钥管理、静态加密和轮换。不要把 Secret 放进镜像层、ConfigMap、日志、错误消息或 Helm values 的明文提交。

数据库 schema migration 应作为受控 Job 运行，并限制并发。例如一个版本只允许一个迁移 Job；迁移先做向后兼容的 expand，再发布可读写新旧格式的应用，最后在确认旧版本已下线后 contract。CronJob 适用于定时清理、报表等可重试任务；必须设计幂等键、超时、`concurrencyPolicy` 和失败告警。

HPA 根据 CPU、内存或外部指标改变副本数。CPU 利用率常以 request 为基准，因此未设置 request 会使含义不完整。HPA 只能扩展应用副本，不能增加数据库连接容量、第三方 API 配额或队列消费者的幂等性；扩容前应设置连接池上限和下游背压。

Helm 将模板、values 和 release 记录组合为可安装 Chart。它适合复用有参数差异的 Kubernetes 清单，但过多条件分支会让渲染结果难以审阅。提交部署前使用 `helm template` 或 GitOps 渲染结果审查最终 YAML；values 文件中不得保存真实生产密钥。

## Terraform：声明、状态与漂移

Terraform 配置中的 provider 调用外部 API，resource 描述对象，module 封装一组可复用资源。Terraform state 记录配置地址与远端对象标识、属性和依赖关系，使 Terraform 能计算从当前状态到目标状态的变更。state 可能含敏感值，即使配置标记为 `sensitive` 也不能把 state 当公开文件。

```hcl
terraform {
  required_version = ">= 1.9.0"
  backend "s3" {
    bucket = "acme-terraform-state"
    key    = "production/network/terraform.tfstate"
    region = "ap-southeast-1"
  }
}

module "network" {
  source = "./modules/network"
  name   = "production"
  cidr   = "10.20.0.0/16"
}
```

远端 state 需要版本、访问控制、加密和锁，避免两次 apply 同时写入。不同环境使用独立 state 边界，减少测试变更影响生产的可能。`terraform plan -out=tfplan` 生成待审查计划；应在相同、受控的运行环境执行对应 `terraform apply tfplan`，不要把计划文件传到不可信位置，因为它可能含敏感数据。

Drift 指配置以外有人或系统改变了远端资源。`terraform plan` 会刷新可读属性并显示差异；它不等价于已修复。先确认漂移是否合理，再选择把远端改回声明状态、更新配置，或把对象从 Terraform 管理中明确移除。对生产网络、IAM 和删除操作，plan 审阅与最小权限执行身份是必要边界。

## 案例二：从提交到灰度发布的流水线

目标：每次合并到主分支都生成可追溯镜像；预发布环境通过 smoke test 后才允许生产灰度；任何健康或业务指标异常都能停止或回滚。

```mermaid
sequenceDiagram
    participant Dev as 开发者
    participant CI as CI
    participant Reg as Registry
    participant K8s as 预发布集群
    participant Prod as 生产集群
    Dev->>CI: push commit
    CI->>CI: lint、test、build、scan
    CI->>Reg: push tag 与 digest
    CI->>K8s: 部署 digest
    CI->>K8s: smoke test
    CI->>Prod: 批准后小比例发布
    Prod-->>CI: 指标与健康结果
    CI->>Prod: 扩大或回滚
```

GitHub Actions workflow 是事件触发的一组 jobs；job 内步骤顺序执行，不同 job 需显式用 `needs` 表达依赖。下面展示结构，不把生产凭据写入仓库。OIDC 或短期云身份应替代长期云访问密钥；部署 job 应使用受保护环境和审批规则。

```yaml
name: verify-and-release
on:
  pull_request: {}
  push:
    branches: [main]
jobs:
  verify:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: {go-version: '1.25.x'}
      - run: go vet ./...
      - run: go test -race ./...
  build:
    if: github.event_name == 'push'
    needs: verify
    runs-on: ubuntu-latest
    permissions: {contents: read, packages: write, id-token: write}
    steps:
      - uses: actions/checkout@v4
      - run: ./scripts/build-scan-and-push.sh "$GITHUB_SHA"
      - run: ./scripts/deploy-staging.sh "$GITHUB_SHA"
      - run: ./scripts/smoke-test.sh https://staging.example.com
  production:
    needs: build
    environment: production
    runs-on: ubuntu-latest
    steps:
      - run: ./scripts/canary-release.sh "$GITHUB_SHA"
```

流水线必须区分验证与发布。PR 可运行无密钥测试、lint、依赖检查和构建；来自 fork 的代码不应直接获得可写 Registry、部署凭据或生产网络。镜像扫描发现漏洞后不能机械地“一律阻断”：应根据可利用路径、修复可用性、基础镜像和组织风险门槛分类，但高危可利用漏洞不能靠忽略标签绕过。

Smoke test 是部署后的短路径验证，例如获取 `/readyz`、以测试账号创建并查询一个资源、验证版本响应头。它不能证明高峰容量或所有权限规则。灰度期间观察请求错误率、延迟分位数、关键业务成功率和资源饱和度；这些阈值、观察窗口、停止条件和负责人必须预先写入发布策略。

失败分支：把 `latest` 部署到生产。之后即使查看 Kubernetes manifest 也无法确定它实际拉取了哪个构建，回滚也可能拉到新的 `latest`。修正是把 commit、构建 run、SBOM/扫描结果、镜像 digest、部署 revision、配置版本与 smoke test 结果关联保存。

## 发布策略的取舍

| 策略 | 流量变化 | 适合 | 主要风险 |
| --- | --- | --- |
| Rolling | 新旧副本逐步替换 | 兼容的新版本 API | 新旧版本并存期间的 schema 不兼容 |
| Blue/Green | 预备完整新环境后切换 | 切换可逆、容量充足 | 两套环境成本与连接状态迁移 |
| Canary | 先给一部分真实流量 | 需要用指标判断风险 | 需要可靠分流、指标和回滚自动化 |
| Feature flag | 同一部署内按规则启用 | 可独立关闭的业务能力 | 旗标债务、权限和默认值错误 |

不要把“灰度”理解为不会出错。它降低爆炸半径，但数据库写入、缓存污染、消息副作用和外部调用可能已发生。对于不可逆操作，先做 shadow、只读验证、幂等设计或人工审批。

## 调试路径与生产边界

先从部署 revision、镜像 digest、Pod event、容器退出码、readiness 状态开始，再看应用结构化日志与 trace。`ImagePullBackOff` 常指向镜像名称、权限或网络；`CrashLoopBackOff` 是反复退出后的退避状态，不是根因；`OOMKilled` 需要检查实际峰值、limit、泄漏和并发，而不是只提高 limit。

常用只读检查：

```sh
kubectl describe pod -l app=api
kubectl logs deploy/api --all-containers=true --tail=200
kubectl get events --sort-by=.lastTimestamp
kubectl rollout history deployment/api
terraform plan -detailed-exitcode
```

`terraform plan -detailed-exitcode` 的退出码 0 表示无差异，2 表示存在差异，1 表示错误；在 CI 中应显式处理它，不能把 2 当成执行失败。生产排障避免直接编辑受 Terraform 或 GitOps 管理的资源；临时缓解若必须修改，应记录并把最终决定同步回声明配置，否则漂移会在下一次发布被覆盖。

## 综合练习与验收

为一个带 PostgreSQL 的 API 建立本地 Compose、生产 Deployment/Service、迁移 Job、Terraform 网络模块和 CI 流水线。

验收条件：

- 镜像不以 root 运行，且构建产物不含源码和开发密钥。
- Pod 在数据库未就绪时不接收流量，数据库短暂错误不会导致所有 Pod 因 liveness 重启。
- 变更可通过 digest 定位到 commit、测试、扫描和部署 revision。
- 部署失败可停止 rollout；应用与数据库的迁移路径向后兼容。
- Terraform state 存在受控远端，生产 apply 经过 plan 审阅，手工变更可被检测。
- 在预发布环境模拟错误版本，确认 smoke test、指标阈值和回滚步骤确实生效。

## 来源

- [Docker：Multi-stage builds](https://docs.docker.com/build/building/multi-stage/)（访问日期：2026-07-23）
- [Kubernetes：Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/)（访问日期：2026-07-23）
- [Kubernetes：Pod 生命周期与探针](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#container-probes)（访问日期：2026-07-23）
- [Terraform：State](https://developer.hashicorp.com/terraform/language/state)（访问日期：2026-07-23）
- [GitHub Actions：Workflows](https://docs.github.com/en/actions/concepts/workflows-and-actions/workflows)（访问日期：2026-07-23）
