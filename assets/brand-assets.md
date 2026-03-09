# Barca-Strategos Brand Assets

## 🐘⚔️ Complete Brand Package

This directory contains all official brand assets for Barca-Strategos Phoenix, implementing the "Le Piège de Cannes" (Battle of Cannae) logo concept.

## 📁 File Structure

```
assets/
├── logo/
│   ├── barca-strategos-logo.svg      # Primary logo
│   ├── barca-strategos-icon.svg      # Icon version
│   ├── barca-strategos-mono.svg      # Monochrome version
│   ├── barca-strategos-logo.png      # PNG for web use
│   ├── barca-strategos-icon.png      # PNG for icons
│   └── favicon.ico                # Windows favicon
├── colors/
│   ├── brand-palette.svg           # Color palette reference
│   └── usage-guidelines.md        # Color usage rules
├── typography/
│   ├── fonts/                     # Font files
│   └── typography-guide.md        # Typography rules
└── templates/
    ├── presentations/               # PowerPoint templates
    ├── documents/                 # Word templates
    └── social-media/              # Social media graphics
```

## 🎨 Logo Usage Guidelines

### Primary Logo (barca-strategos-logo.svg)
- **Use**: Website headers, main applications, marketing materials
- **Minimum Size**: 64px height
- **Clear Space**: Equal to logo height on all sides
- **Background**: Ensure contrast (light or dark versions available)

### Icon Version (barca-strategos-icon.svg)
- **Use**: Mobile apps, favicons, notifications
- **Sizes**: 16px, 32px, 64px, 128px
- **Simplification**: Maintains recognizability at small sizes

### Monochrome Version (barca-strategos-mono.svg)
- **Use**: Security documents, technical specifications
- **Single Color**: Noir Profond (#1A1A1A)
- **Purpose**: High contrast, professional appearance

## 🎨 Color Implementation

### Primary Colors
```css
/* CSS Variables */
--rouge-carthaginois: #E2725B;  /* Tunisia earth, battle passion */
--bronze-antique: #CD7F32;      /* Rust robustness, historical armor */
--noir-profond: #1A1A1A;      /* Black-box invisibility, technical elegance */
```

### Usage Rules
- **Rouge Carthaginois (60%)**: Primary actions, alerts, CTAs
- **Bronze Antique (30%)**: Headers, icons, accents
- **Noir Profond (10%)**: Text, backgrounds, borders

## 📝 Typography Implementation

### Font Stack
```css
/* Font Imports */
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&family=Open+Sans:wght@400;600&display=swap');

/* CSS Variables */
--font-primary: 'Inter', sans-serif;
--font-secondary: 'Open Sans', sans-serif;
```

### Typography Scale
```css
/* Headings */
--text-4xl: 2.5rem;   /* 40px */
--text-3xl: 2rem;      /* 32px */
--text-2xl: 1.5rem;    /* 24px */
--text-xl: 1.25rem;    /* 20px */
--text-lg: 1.125rem;   /* 18px */
--text-base: 1rem;     /* 16px */
--text-sm: 0.875rem;   /* 14px */
--text-xs: 0.75rem;    /* 12px */
```

## 🖼️ Application Examples

### Website Header
```html
<header class="phoenix-header">
  <div class="logo-container">
    <img src="assets/logo/barca-strategos-logo.svg" alt="Barca-Strategos Phoenix" class="logo">
    <nav class="navigation">
      <a href="/features" class="nav-link">Features</a>
      <a href="/docs" class="nav-link">Documentation</a>
      <a href="/deploy" class="nav-link">Deploy</a>
    </nav>
  </div>
</header>
```

### Security Document Header
```html
<div class="security-header">
  <img src="assets/logo/barca-strategos-mono.svg" alt="Barca-Strategos" class="security-logo">
  <div class="classification">CONFIDENTIAL</div>
  <div class="document-title">Security Assessment Report</div>
</div>
```

### Mobile App Icon
```html
<!-- Android -->
<link rel="icon" type="image/png" sizes="192x192" href="assets/logo/barca-strategos-icon-192.png">

<!-- iOS -->
<link rel="apple-touch-icon" sizes="180x180" href="assets/logo/barca-strategos-icon-180.png">

<!-- Favicon -->
<link rel="icon" type="image/x-icon" href="assets/logo/favicon.ico">
```

## 📱 Social Media Assets

### Profile Pictures
- **Twitter**: 400x400px logo with full branding
- **LinkedIn**: 300x300px professional version
- **Discord**: 128x128px server icon
- **GitHub**: Profile avatar and repository header

### Banner Templates
- **Twitter Header**: 1500x500px campaign banner
- **LinkedIn Banner**: 1200x627px company banner
- **Website Hero**: 1920x1080px hero section background

## 🔐 Security Considerations

### Asset Protection
- **Watermarking**: Subtle logo placement on sensitive documents
- **Classification**: Color-coded security levels with brand integration
- **Access Control**: Tiered access to brand asset repository
- **Version Control**: Track all asset versions and modifications

### Digital Security
- **Metadata**: Proper titles, descriptions, and alt text
- **Certificates**: SSL validation for brand domains
- **Anti-Counterfeiting**: Official distribution channels only

## 📊 Asset Performance

### File Formats
- **SVG**: Vector graphics for scalability
- **PNG**: Raster for web use with transparency
- **JPG**: Print materials without transparency
- **ICO**: Windows favicon with multiple sizes

### Optimization
- **SVG**: Optimized for web (compressed paths)
- **PNG**: Progressive loading, appropriate compression
- **File Size**: Logo under 10KB, icon under 5KB

This brand package ensures consistent, professional representation of Barca-Strategos Phoenix across all platforms and applications while maintaining the tactical intelligence and historical legacy that defines our framework.
