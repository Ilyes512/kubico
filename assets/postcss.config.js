module.exports = {
    plugins: {
        'postcss-nested-ancestors': {},
        'tailwindcss/nesting': 'postcss-nested',
        tailwindcss: {},
        ...process.env.NODE_ENV === 'production'
            ? {
                'postcss-preset-env': {
                    features: {
                        'nesting-rules': false
                    },
                },
                '@fullhuman/postcss-purgecss': {
                    content: ['./../templates/**/*.tmpl'],
                    defaultExtractor: content => content.match(/[A-Za-z0-9-_:/]+/g) || [],
                },
                'cssnano': {},
            }
            : {},
      },
};
