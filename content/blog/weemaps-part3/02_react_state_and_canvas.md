
# Building a "Human-in-the-Loop" Interface for Local AI (Part 2)

**Topic:** React State, Canvas, and Coordinate Systems
**Stack:** React, SVG, LocalStorage

In the final part of this series, we look at the "glue" that holds the application together. How do we make the grid line up perfectly? How do we save our work if the browser refreshes?

## 1. The Grid Calibration Challenge

Overlaying a perfect mathematical hex grid onto a pixel-art screenshot is surprisingly hard.

*   **The Problem:** The user's image might be zoomed in, cropped weirdly, or use non-standard hex dimensions.
*   **The Solution:** We need a calibration tool that allows sub-pixel precision.

![Image: Grid Calibration Sidebar](path/to/calibration_sidebar.png)
*> Caption: The Calibration controls allow floating point adjustments. On the map, the green lines update in real-time.*

We separated the **Logical Grid** from the **Visual Grid**.

### The Logical Grid (The Math)
We store the configuration in a simple state object:
```typescript
interface HexConfig {
  originX: 32.5, // Float for precision
  originY: 40.0,
  width: 64.0,
  height: 72.0
}
```
We use `input type="number" step="0.1"` to allow users to nudge the grid by fractions of a pixel. In pixel art analysis, being off by 0.5px can mean the difference between cropping a unit's head or cutting it off.

### The Visual Grid (The Render)
We render the grid using SVG overlaid on top of the image.

Why SVG and not Canvas?
*   **Crispness:** SVG lines remain sharp at any zoom level.
*   **Interactivity:** We don't need it here, but SVG allows attaching click handlers to specific grid cells easier than calculating mouse positions on a generic canvas.
*   **CSS Styling:** We can easily add drop-shadows or glow effects to the grid lines (like the "Golden Pulse" on the highlighted tile) using standard CSS classes.

## 2. The Zoom Problem

We added a "Zoom" slider. This introduces a coordinate system headache.

If the user zooms to `2.0x` and clicks the screen at `x: 500`, the actual image pixel is at `x: 250`.

We handled this by applying the zoom via CSS `transform: scale(zoom)`, but keeping the DOM container logic simple.

```mermaid
graph TD
    A[Mouse Click (Screen Coords)] -->|Minus Offset| B[Container Coords]
    B -->|Divide by Zoom| C[Image Coords]
    C -->|Hex Formula| D[Axial Q,R]
    
    style C fill:#ff9,stroke:#333
```

```typescript
const handleMapClick = (e) => {
    // Get click position relative to the container
    const rect = e.currentTarget.getBoundingClientRect();
    const clickX = e.clientX - rect.left;
    
    // Reverse the zoom math to find the "Real" pixel
    const realX = (clickX + scrollOffset) / zoom;
    
    // Use realX for the Hex math
    const hex = pixelToHex(realX, ...);
}
```

This allows the user to inspect pixels closely while the underlying math engine always operates on the original image dimensions.

## 3. State Management: The "Brain" Dump

Since this is a client-side app, if you refresh the page, your trained model vanishes. That’s bad UX.

We needed a way to Save/Load the "Knowledge Base."

### Limitations of TensorFlow.js
You can save a TensorFlow model, but KNN classifiers are unique - they are defined by the *data* they hold, not just weights.

### The JSON Solution
We opted for a "Hydration" strategy. We don't save the vectors. We save the **Source Data**.

When you export the Knowledge Base, we create a JSON file containing:
1.  **Base64 Images:** The tiny cropped images of every tile you labeled.
2.  **Labels:** The text tags for those tiles.

```json
{
  "examples": [
    {
      "id": "10,5",
      "imageBase64": "data:image/png;base64,iVBORw...",
      "labels": { "terrain": "Forest", "unit": "Tank" }
    }
  ]
}
```

**On Import:**
1.  We load the JSON.
2.  We iterate through every example.
3.  We re-run `mobilenet.infer()` on the Base64 image to regenerate the vector.
4.  We re-populate the KNN classifier.

This makes the save files slightly larger (because of the images), but it makes the system **robust**. You can upgrade the underlying MobileNet model version in the future, and your save file will still work because it rebuilds the math from the raw pixels every time.

---

## Conclusion

**WeeMap Scanner** was a fun experiment in combining old-school game geometry with modern browser-based AI.

By keeping the ML simple (Transfer Learning + KNN) and focusing heavily on the "Teacher" UX, we built a tool that feels powerful, responsive, and surprisingly accurate - all without a single API call to the cloud.

The code is available in the repo. Go forth and reverse-engineer your childhood favorites!
