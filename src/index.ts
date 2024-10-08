/// @ts-ignore
import { default as mod } from "../out/lib.cjs";

declare namespace types {
  type TemplateFunc = (...args: unknown[]) => unknown;

  type TemplateFuncs = Record<string, TemplateFunc>;

  class Template {
    constructor(name: string);
    option(...options: string[]): Template;
    definedTemplates(): string;
    new(name: string): Template;
    funcs(funcs: TemplateFuncs): Template;
    parse(text: string): Template;
    execute(context: unknown): string;
  }
}

export const Template: typeof types.Template = mod.Template;
