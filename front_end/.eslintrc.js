module.exports = {
    root: true,
    ignorePatterns: ["node_modules/"],
    parser: "@typescript-eslint/parser",
    plugins: ["react-hooks"],
    extends: [
      "eslint:recommended",
      "plugin:prettier/recommended",
      "plugin:react/recommended",
      "plugin:jsx-a11y/strict",
    ],
    settings: {
      react: {
        // Tells eslint-plugin-react to automatically detect the version
        // of React to use
        version: "detect",
      },
    },
    rules: {
      "react-hooks/rules-of-hooks": "error",
      "react-hooks/exhaustive-deps": ["error", {"additionalHooks": "useRecoilCallback"}],
      "react/prop-types": "off",
      "jsx-a11y/anchor-is-valid": "off",
      "react/display-name": "off",
    },
  };
