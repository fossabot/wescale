package tabletserver

type WasmVM interface {
	GetRuntimeType() string
	InitRuntime() error
	GetWasmModule(key string) (bool, WasmModule)
	InitWasmModule(key string, wasmBytes []byte) (WasmModule, error)
	ClearWasmModule(key string)
}

type WasmModule interface {
	NewInstance(qre *QueryExecutor) (WasmInstance, error)
}

type WasmInstance interface {
	RunWASMPlugin() error
	RunWASMPluginAfter(args *WasmPluginExchangeAfter) (*WasmPluginExchangeAfter, error)
}
