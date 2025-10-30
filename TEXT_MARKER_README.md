# Text Marker Feature - Schola Theologiae

## Overview
The Text Marker feature allows users to highlight and save important text passages while reading book chapters. It's designed with a mobile-first approach and provides persistent storage across browser sessions.

## Features

### 🖍️ Text Highlighting
- **Floating Marker Button**: A circular button appears in the top-right corner on book chapter pages
- **Toggle Marking Mode**: Click the button to activate/deactivate text selection mode
- **Visual Feedback**: The button changes color and shows animation when marking mode is active

### 📱 Mobile-First Design
- **Touch-Friendly**: Optimized for mobile touch interactions
- **Responsive**: Button and UI elements adapt to different screen sizes
- **Accessibility**: Proper touch targets and visual feedback

### 💾 Persistent Storage & Restoration
- **Enhanced localStorage**: All marked texts are saved with contextual information
- **Intelligent Restoration**: Automatically restores highlights when page refreshes
- **Multiple Restoration Methods**: 
  1. Context-aware matching (using surrounding text)
  2. Exact text matching 
  3. Flexible partial matching for edge cases
- **Visual Distinction**: Restored markings have a green tint vs. yellow for new markings
- **Session Persistence**: Markings survive page refreshes, browser restarts, and navigation

### 🎯 User Interactions

#### Marking Text
1. Click the marker button to enter marking mode
2. Select any text in the chapter
3. The text becomes highlighted with a yellow background
4. Click the button again to exit marking mode

#### Removing Markings
- Click on any highlighted text to remove the marking
- Confirm the removal in the dialog that appears

#### Viewing Marked Texts
- **Long-press** (800ms) the marker button to see a summary of all marked texts on the current page
- The summary shows marked text snippets and their creation dates

### 🎨 Visual Design
- **Highlighted Text**: Yellow gradient background with left border accent
- **Hover Effects**: Highlighted text has subtle hover animations
- **Toast Notifications**: Non-intrusive feedback messages for user actions
- **Color Coding**: 
  - Blue button: Normal state
  - Gold button: Active marking mode
  - Green button: Long-press action

## Technical Implementation

### Files Structure
```
app/javascript/text_marker.js     # Main functionality
app/assets/stylesheets/text_marker.css  # Styling
config/importmap.rb               # JavaScript module registration
app/assets/stylesheets/application.css  # CSS import
```

### Browser Compatibility
- Modern browsers with localStorage support
- Touch events for mobile devices
- CSS grid and flexbox support

### Storage Format
```javascript
{
  "/books/catecismo-pio-x/parte1/1": [
    {
      id: "marker_1234567890_abc123",
      text: "Selected text content",
      timestamp: 1703123456789,
      page: "/books/catecismo-pio-x/parte1/1",
      context: {
        before: "text appearing before selection",
        after: "text appearing after selection", 
        parentTag: "P",
        parentClass: "book-content"
      }
    }
  ]
}
```

## New Features in v2.0

### 🔄 **Intelligent Text Restoration**
The system now uses multiple strategies to restore marked texts:

1. **Context-Aware Matching**: Uses surrounding text to locate markings precisely
2. **Exact Text Search**: Falls back to direct text matching
3. **Flexible Partial Matching**: Handles cases where text may have changed slightly

### 🎨 **Visual Improvements**
- **Green highlights** for restored markings vs. **yellow** for new ones
- **Smooth animations** when markings are restored
- **Partial match indicators** with dashed borders for incomplete restorations

### 📊 **Enhanced Feedback**
- Toast notifications show how many markings were restored
- Console logging for debugging restoration issues
- Better error handling and recovery

## Usage Instructions

### For Readers
1. Navigate to any book chapter (Catecismo or Summa Theologiae)
2. Look for the circular marker button in the top-right corner
3. Tap the button to activate marking mode (button turns gold)
4. Select text you want to highlight
5. Tap the button again to deactivate marking mode
6. Long-press the button to see your marked texts

### For Developers
The TextMarker class is automatically initialized on book chapter pages and provides:
- `getMarkedTextsForPage()`: Get markings for current page
- `getAllMarkedTexts()`: Get all markings across all pages
- `clearAllMarkedTexts()`: Remove all markings (with confirmation)
- `showMarkedTextsSummary()`: Display summary dialog

## Future Enhancements
- Export marked texts to notes
- Search within marked texts
- Categories/tags for markings
- Sharing marked passages
- Offline synchronization
- Better text restoration across page changes
