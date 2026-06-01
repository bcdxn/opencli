/// <reference types="vite/client" />

// Go class injected into global scope by wasm_exec.js (loaded in index.html)
declare class Go {
  importObject: WebAssembly.Imports;
  run(instance: WebAssembly.Instance): Promise<void>;
}
