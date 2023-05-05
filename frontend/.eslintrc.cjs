module.exports = {
  root: true,
  extends: ["@nuxt/eslint-config"],
  rules: {
    "vue/max-attributes-per-line": ["error", {
      "singleline": {
        "max": 2
      },      
      "multiline": {
        "max": 3
      }
    }],
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
