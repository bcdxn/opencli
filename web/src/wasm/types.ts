export interface ValidationError {
  message: string;
  path: string;
}

export interface ValidationResult {
  valid: boolean;
  errors: ValidationError[];
}

export interface GenerationResult {
  output: string;
  error: string;
}

declare global {
  interface Window {
    validateOCS: (input: string, format: string) => ValidationResult;
    generateOCSDocs: (
      input: string,
      inputFormat: string,
      outputFormat: string,
      htmlFlavor: string,
    ) => GenerationResult;
  }
}
