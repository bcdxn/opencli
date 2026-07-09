#!/usr/bin/env node

import yargs from "yargs";
import { hideBin } from "yargs/helpers";
import { run } from "./gencli/run";
import { Actions } from "./actions";

// Parse arguments using yargs
async function main() {
  const actions = new Actions();
  await run(yargs(hideBin(process.argv)), actions);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
