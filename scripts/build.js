import { dtsPlugin } from "esbuild-plugin-d.ts";
import { build } from "esbuild";

build({
  entryPoints: ["./src/index.ts"],
  outdir: "./dist",
  bundle: true,
  format: "esm",
  minify: true,
  plugins: [
    dtsPlugin({
      experimentalBundling: true,
    }),
  ],
});
