const purgecss = require('@fullhuman/postcss-purgecss')({
    content: ['./../templates/**/*.tmpl'],
    defaultExtractor: content => content.match(/[A-Za-z0-9-_:/]+/g) || [],
});

module.exports = {
    plugins: [
        require('tailwindcss'),
        require('autoprefixer'),
        require('postcss-nested-ancestors'),
        require('postcss-nested'),
        ...process.env.NODE_ENV === 'production'
            ? [purgecss, require('cssnano')]
            : [],
    ],
};
