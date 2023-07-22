module.exports = {
  root: true,
  parser: "vue-eslint-parser",
  parserOptions: {
    parser: "@typescript-eslint/parser",
  },
  extends: ["@nuxt/eslint-config", "plugin:prettier/recommended"],
  rules: {
    "vue/max-attributes-per-line": [
      "error", { 
        "singleline": { "max": 4 },      
        "multiline": { "max": 4 } 
      }
    ],
  },
  overrides: [
    { 
      files: './components/MarkdownHTML.vue',
      rules: {
        'vue/no-v-html': 'off'
      }
    }
  ]
};
