# CSS Architecture Documentation

## Overview

The Schola Theologiae CSS architecture uses a centralized system of CSS custom properties (variables) defined in `application.css`. This approach ensures consistency, maintainability, and ease of theming across the entire application.

## CSS Variables Reference

### Colors

#### Primary Colors
- `--color-primary: #8b4513` - Main brown color for headers and important text
- `--color-primary-hover: #a0522d` - Hover state for primary color
- `--color-accent: #d4af37` - Golden accent color
- `--color-accent-dark: #bfa14a` - Darker shade of accent
- `--color-accent-darker: #afa11b` - Even darker accent shade

#### Text Colors
- `--color-text-primary: #333` - Main body text
- `--color-text-secondary: #666` - Secondary/muted text
- `--color-text-tertiary: #999` - Tertiary/very muted text
- `--color-text-light: #555` - Light text variant

#### Background Colors
- `--color-background-primary: #fcfad9` - Main page background
- `--color-background-secondary: #f8f5e4` - Secondary background
- `--color-background-tertiary: #f0ead6` - Tertiary background
- `--color-background-white: #ffffff` - Pure white background

#### Border Colors
- `--color-border-light: #e8dcc0` - Light border color
- `--color-border-medium: #d4c49a` - Medium border color

#### Navigation Colors
- `--color-nav-bg: #2c2c2c` - Navigation bar background
- `--color-nav-bg-dark: #1a1a1a` - Darker navigation background
- `--color-nav-text: #ffffff` - Navigation text color
- `--color-nav-text-muted: #cccccc` - Muted navigation text

### Typography

- `--font-primary: "Cascadia Code", "Roboto", sans-serif` - Main body font
- `--font-decorative: 'Tagesschrift', 'Eagle Lake', serif` - Decorative headers
- `--font-medieval: 'MedievalSharp', serif` - Medieval-style headers
- `--font-skranji: 'Skranji', serif` - Alternative decorative font

### Spacing Scale

- `--spacing-xs: 0.25rem` (4px)
- `--spacing-sm: 0.5rem` (8px)
- `--spacing-md: 1rem` (16px)
- `--spacing-lg: 1.5rem` (24px)
- `--spacing-xl: 2rem` (32px)
- `--spacing-2xl: 3rem` (48px)

### Border Radius

- `--radius-sm: 4px` - Small corners
- `--radius-md: 8px` - Medium corners
- `--radius-lg: 12px` - Large corners
- `--radius-xl: 15px` - Extra large corners
- `--radius-2xl: 20px` - Double extra large corners
- `--radius-full: 50%` - Circular/pill shape

### Shadows

- `--shadow-sm: 0 2px 4px rgba(0, 0, 0, 0.05)` - Subtle shadow
- `--shadow-md: 0 4px 12px rgba(0, 0, 0, 0.1)` - Medium shadow
- `--shadow-lg: 0 8px 24px rgba(0, 0, 0, 0.15)` - Large shadow
- `--shadow-xl: 0 12px 35px rgba(0, 0, 0, 0.2)` - Extra large shadow

### Transitions

- `--transition-fast: 0.2s ease` - Quick transitions
- `--transition-normal: 0.3s ease` - Standard transitions
- `--transition-slow: 0.4s ease` - Slower transitions

## File Structure

```
app/assets/stylesheets/
├── application.css       # Main entry point with CSS variables
├── nav_bar.css          # Navigation bar styles
├── home.css             # Home page specific styles
├── footer.css           # Footer styles
├── books.css            # Books listing page
├── book.css             # Individual book pages
├── pages.css            # Common page styles
├── search.css           # Search modal and results
└── exame_maturidade.css # Exame de Maturidade page
```

## Usage Guidelines

### Using CSS Variables

Always use CSS variables instead of hardcoded values:

```css
/* ✅ Good */
.button {
  background: var(--color-primary);
  padding: var(--spacing-md);
  border-radius: var(--radius-md);
  transition: var(--transition-normal);
}

/* ❌ Bad */
.button {
  background: #8b4513;
  padding: 16px;
  border-radius: 8px;
  transition: 0.3s ease;
}
```

### Performance Optimization

Use `will-change` for elements that will animate:

```css
.animated-card {
  transition: transform var(--transition-normal);
  will-change: transform;
}

.animated-card:hover {
  transform: translateY(-5px);
}
```

### Accessibility

Always include focus-visible styles:

```css
.interactive-element:focus-visible {
  outline: 2px solid var(--color-accent);
  outline-offset: 2px;
}
```

### Responsive Design

Use consistent breakpoints:

```css
/* Mobile first approach */
@media (max-width: 480px) { /* Small mobile */ }
@media (max-width: 768px) { /* Tablet */ }
@media (max-width: 1024px) { /* Small desktop */ }
@media (min-width: 1400px) { /* Large desktop */ }
```

## Best Practices

1. **Mobile-First**: Always design for mobile first, then add larger breakpoints
2. **Semantic HTML**: Use proper HTML5 semantic elements
3. **BEM Naming**: Consider BEM (Block Element Modifier) for complex components
4. **Comments**: Add comments to explain complex CSS sections
5. **Specificity**: Keep specificity low; avoid `!important`
6. **DRY**: Don't repeat yourself - use variables for repeated values
7. **Performance**: 
   - Avoid expensive properties like `box-shadow` on scroll
   - Use `transform` and `opacity` for animations (GPU accelerated)
   - Add `will-change` for elements that will animate

## Common Patterns

### Card Component
```css
.card {
  background: var(--color-background-white);
  border-radius: var(--radius-xl);
  padding: var(--spacing-xl);
  box-shadow: var(--shadow-lg);
  border: 2px solid var(--color-border-light);
  transition: transform var(--transition-normal);
}

.card:hover {
  transform: translateY(-5px);
  box-shadow: var(--shadow-xl);
}
```

### Button Component
```css
.btn {
  background: linear-gradient(135deg, var(--color-primary), var(--color-primary-hover));
  color: var(--color-background-white);
  padding: var(--spacing-md) var(--spacing-xl);
  border-radius: 25px;
  border: none;
  font-weight: 600;
  cursor: pointer;
  transition: all var(--transition-normal);
}

.btn:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}
```

## Maintenance

When adding new features:

1. Check if needed colors/spacing already exist in variables
2. Add new variables to `application.css` if creating new design tokens
3. Follow existing patterns and naming conventions
4. Test on multiple screen sizes
5. Ensure keyboard accessibility with focus-visible styles
6. Document any new patterns or components

## Browser Support

The CSS uses modern features that are supported in:
- Chrome/Edge 88+
- Firefox 85+
- Safari 14+

CSS Custom Properties (variables) are used throughout, which are well-supported in all modern browsers.

## Future Improvements

Potential areas for enhancement:

1. **Dark Mode**: Add alternative color scheme using CSS variables
2. **CSS Grid**: Expand use of CSS Grid for complex layouts
3. **Container Queries**: Use when browser support improves
4. **CSS Nesting**: Consider when native CSS nesting is widely supported
5. **Logical Properties**: Use `margin-inline`, `padding-block` for better RTL support
