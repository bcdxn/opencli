export {};

declare global {
  interface Window {
    // Replace 'any' with a specific function signature if you know it, e.g., () => void
    OcliDocs?: any;
  }

  // Provided by wasm_exec.js loaded as a beforeInteractive script
  class Go {
    importObject: WebAssembly.Imports;
    run(instance: WebAssembly.Instance): void;
  }
}
