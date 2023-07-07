package consts

// 运行环境
const (
	//ENV
	RunEnv = "RUN_ENV"
	// 本地
	RunEnvLocal = "local"
	// 开发
	RunEnvDev = "dev"
	// 测试
	RunEnvTest = "test"
	// 预发布
	RunEnvStag = "stag"
	// 生产
	RunEnvProd = "prod"
	// gray
	RunEnvGray = "gray"
	// 生产
	RunEnvProdTw = "prodtw"
	// NON
	RunEnvNull = "null"
	// 极星本地
	RunEnvBjxLocal = "bjxlocal"
	// 极星测试
	RunEnvBjxTest = "bjxtest"
	// 极星线上
	RunEnvBjxProd    = "bjxprod"
	RunEnvBjxFuseK8s = "fuse_k8s"
	// uat环境
	RunEnvBjxFuat_k8s = "uat_k8s"
)

// Jaeger tracing key
const JaegerTracingKey = "jaeger-tracing-key"

// IdentityUserHeader 传递userId的header
const IdentityUserHeader = "X-POLARIS-IDENTITY-USER"

// IdentityOrgHeader 传递orgId的header
const IdentityOrgHeader = "X-POLARIS-IDENTITY-ORG"

const ContextTracerKey = "Tracer-context"

const OpenTraceKey = "opentraceId"
