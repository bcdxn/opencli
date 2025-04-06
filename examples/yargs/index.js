#!/usr/bin/env node
const { Cli } = require("./cli.gen");
const { CliHandlersInterface } = require("./cli-interface.gen")

/**
 * @class
 * @implements {CliHandlersInterface}
 */
class HandlersImpl extends CliHandlersInterface {
  async pleasantriesFarewell(yargv, args, flags) {
    if (flags.language === "spanish") {
      console.log("hola", args.name);
    } else {
      console.log("hello", args.name);
    }
  }

  async pleasantriesFarewell(yargv, args, flags) {
    if (flags.language === "spanish") {
      console.log("adios", args.name);
    } else {
      console.log("bye", args.name);
    }
  }
}

async function main() {
  let version = "1.0.0"

  let cli = new Cli(new HandlersImpl(), version);
  try {
    await cli.run();
  } catch (err) {
    console.log(err);
  }
}

main();