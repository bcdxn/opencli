#!/usr/bin/env node

import { run } from "./gencli/run";
import { Actions } from "./actions";

// Parse arguments using yargs
async function main() {
  const actions = new Actions();
  await run(process.argv, actions);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
