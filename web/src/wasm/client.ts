import type { ValidationResult, GenerationResult } from "./types";

let _wasmReady = false;

/**
 * Loads and instantiates the OpenCLI WASM binary.
 * After this resolves, window.validateOCS and window.generateOCSDocs are available.
 */
export async function loadWasm(): Promise<void> {
  // Go class is provided by wasm_exec.js loaded in index.html
  const go = new Go();
  const result = await WebAssembly.instantiateStreaming(
    fetch(`${import.meta.env.BASE_URL}opencli.wasm`),
    go.importObject,
  );
  // go.run() starts the Go scheduler; main() registers JS globals synchronously
  // before blocking on select{}, so functions are available after this call.
  go.run(result.instance);
  _wasmReady = true;
}

export function isWasmReady(): boolean {
  return _wasmReady;
}

export function validateOCS(input: string, format: string): ValidationResult {
  return window.validateOCS(input, format);
}

export function generateOCSDocs(
  input: string,
  inputFormat: string,
  outputFormat: string,
  htmlFlavor: string,
): GenerationResult {
  return window.generateOCSDocs(input, inputFormat, outputFormat, htmlFlavor);
}
