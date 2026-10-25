import js from '@eslint/js'
import globals from 'globals'
import reactHooks from 'eslint-plugin-react-hooks'
import boundaries from 'eslint-plugin-boundaries'
import tseslint from 'typescript-eslint'

// FSD layers (ADR-0002): app > pages > widgets > features > entities > shared.
// A slice may import only lower layers, only through their public index.ts,
// and never a sibling slice on the same layer.
const layers = ['app', 'pages', 'widgets', 'features', 'entities', 'shared']
const below = (l) => layers.slice(layers.indexOf(l) + 1)

export default tseslint.config(
  { ignores: ['dist', 'coverage', 'node_modules', 'test-results', 'playwright-report', 'src/shared/api/schema.gen.ts'] },
  {
    files: ['**/*.{ts,tsx}'],
    extends: [js.configs.recommended, ...tseslint.configs.recommended],
    languageOptions: { ecmaVersion: 2022, globals: globals.browser },
    plugins: { 'react-hooks': reactHooks },
    rules: {
      'react-hooks/rules-of-hooks': 'error',
      'react-hooks/exhaustive-deps': 'error',
      'no-restricted-syntax': [
        'error',
        { selector: "JSXAttribute[name.name='dangerouslySetInnerHTML']", message: 'Render user content as text only (SEC).' },
        { selector: "MemberExpression[object.name='sessionStorage']", message: 'No Web Storage for session data (SEC).' },
      ],
      'no-restricted-imports': ['error', { patterns: [{ group: ['../../*'], message: 'Use @/ aliases across slices.' }] }],
    },
  },
  {
    files: ['src/**/*.{ts,tsx}'],
    plugins: { boundaries },
    settings: {
      'import/resolver': { typescript: { project: './tsconfig.json' } },
      'boundaries/include': ['src/**/*'],
      'boundaries/elements': [
        { type: 'app', pattern: 'src/app' },
        { type: 'shared', pattern: 'src/shared/*', capture: ['segment'] },
        ...['pages', 'widgets', 'features', 'entities'].map((l) => ({ type: l, pattern: `src/${l}/*`, capture: ['slice'] })),
      ],
    },
    rules: {
      // Same-slice (internal) imports are not checked; everything else must go down a layer
      // through the slice's public index.ts. Sibling slices on one layer are disallowed.
      'boundaries/dependencies': [
        'error',
        {
          default: 'disallow',
          policies: [
            ...['app', 'pages', 'widgets', 'features', 'entities'].map((l) => ({
              from: { element: { type: l } },
              allow: { to: { element: { type: below(l).filter((t) => t !== 'shared'), fileInternalPath: 'index.ts' } } },
            })),
            { from: { element: { type: '*' } }, allow: { to: { element: { type: 'shared' } } } },
          ],
        },
      ],
    },
  },
)
