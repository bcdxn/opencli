import schema from "./spec.schema.json";

export const schemaVersion = String(schema.properties.opencliVersion.enum[0]);
