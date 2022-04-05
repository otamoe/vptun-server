module.exports = {
  root: true,
  env: {
    node: true,
    es6: true,
    browser: true,
    mocha: true,
  },
  extends: [
    "alloy",
    "alloy/typescript",
    "plugin:vue/essential",
    "@vue/typescript",

    "prettier",
    "plugin:prettier/recommended",
  ],
  rules: {
    "no-console": "off",
    "no-debugger": process.env.NODE_ENV === "production" ? "error" : "off",
    complexity: ["error", 40],
    "max-params": ["error", 10],
    "@typescript-eslint/no-empty-interface": "off",
    "@typescript-eslint/prefer-for-of": "off",
    "guard-for-in": "off",
  },
  parserOptions: {
    parser: "@typescript-eslint/parser",
  },
  plugins: ["@typescript-eslint", "prettier"],
  overrides: [
    {
      files: [
        "**/__tests__/*.{j,t}s?(x)",
        "**/tests/unit/**/*.spec.{j,t}s?(x)",
      ],
      env: {
        jest: true,
      }
    }
  ]
}
