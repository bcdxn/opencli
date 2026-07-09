import { ActionsInterface } from "./gencli/actions";
import { CommandPrintData } from "./gencli/types";
import {
  PleasantriesFarewellArgs,
  PleasantriesFarewellFlags,
  PleasantriesGreetArgs,
  PleasantriesGreetFlags,
} from "./gencli/params";
import { defaultHelpFn, defaultUsageFn } from "./gencli/help";
import { help, usage } from "yargs";

export class Actions implements ActionsInterface {
  async PleasantriesGreet(
    args: PleasantriesGreetArgs,
    flags: PleasantriesGreetFlags,
  ): Promise<void> {
    if (flags.language == "english") {
      console.log("hello", args.name);
    } else {
      console.log("hola", args.name);
    }
  }

  async PleasantriesFarewell(
    args: PleasantriesFarewellArgs,
    flags: PleasantriesFarewellFlags,
  ): Promise<void> {
    if (flags.language == "english") {
      console.log("good bye", args.name);
    } else {
      console.log("adios", args.name);
    }
  }

  help(cmd: CommandPrintData): void {
    defaultHelpFn(cmd);
  }
  usage(cmd: CommandPrintData): void {
    defaultUsageFn(cmd);
  }
}
