
# Building a "Human-in-the-Loop" Interface for Local AI (Part 2.3)

**Topic:** Performance, Polish, and The "Scroll Sync"
**Stack:** React Refs, DOM APIs, Keyboard Events

In the final section of this deep dive, we need to talk about **Performance** and **Polish**. 

It’s easy to build a React app that renders a list of items. It’s hard to build a React app that renders 2,000 interactive SVG nodes on top of a 4K image while maintaining 60 FPS scrolling and syncing state across two different panels.

Here is how we optimized **WeeMap Scanner**.

## 1. The DOM Node Explosion (Optimization)

Our `GridOverlay` component renders the hex grid. 

```tsx
// Naive implementation
return (
  <svg>
    {tiles.map(t => (
      <g>
        <path d="..." /> 
        <text>{t.q},{t.r}</text>
      </g>
    ))}
  </svg>
)
```

In a standard *Advance Wars* map (30x20), that is 600 tiles. That means 1,200 DOM elements (Path + Text). React handles this fine.

But if a user uploads a large map (e.g., *Civilization V* screenshot at 100x100), we suddenly have **20,000 DOM nodes**. 

If you try to render 20,000 `<text>` elements with shadows and strokes, the browser's layout engine chokes. Scrolling becomes a slideshow.

### The Fix: Conditional Complexity
We realized the user doesn't *always* need to see the `Q,R` coordinates. They mostly need to see the grid lines to align the map.

We added a toggle `showLabels` to the config state.

```tsx
// GridOverlay.tsx
{showLabels && (
  <text ...>{q},{r}</text>
)}
```

By default, we allow the user to turn this off. Removing those text nodes cuts the DOM size in half and removes the most expensive rendering operation (text layout), restoring 60 FPS immediately.

*Note: A more advanced solution would be using an HTML5 Canvas layer for the grid instead of SVG, but SVG offers superior distinct click-handling logic for irregular shapes like hexagons.*

## 2. The "Scroll Sync" (UX Polish)

A common pattern in data tools is the "List vs. Map" split.
*   **The List:** Shows structured data (confidence scores, labels).
*   **The Map:** Shows spatial context.

We needed them to talk to each other.

### Sidebar -> Map
When the user clicks a result in the sidebar, we don't just highlight the tile. We **scroll** the map to it.

Since our map is inside a `overflow: auto` container and potentially zoomed in, we can't just use `element.scrollIntoView()`. We have to do the math.

```typescript
const handleResultClick = (q, r) => {
    // 1. Calculate Hex Center in pixels
    const { x, y } = getHexCenter(q, r, config);
    
    // 2. Adjust for Zoom
    const scrollLeft = (x * zoom) - (containerWidth / 2);
    const scrollTop = (y * zoom) - (containerHeight / 2);
    
    // 3. Smooth Scroll
    container.scrollTo({
        left: scrollLeft,
        top: scrollTop,
        behavior: 'smooth'
    });
};
```

### Map -> Sidebar
When the user clicks a tile on the map, we need to find that entry in the long list of results in the sidebar.

We assigned a unique DOM ID to every result item: `id={`result-item-${q}-${r}`}`.

```typescript
// After clicking the map...
setTimeout(() => {
    const el = document.getElementById(`result-item-${q}-${r}`);
    el?.scrollIntoView({ block: 'nearest', behavior: 'smooth' });
}, 50);
```

This bi-directional syncing creates a "tight" feeling tool. The user feels like they are manipulating one object, just looking at it through two different lenses.

## 3. Power User Features: Keyboard Nav

Finally, data entry is boring. Clicking miles of UI with a mouse is slow.

To make the "Review" phase fast, we added keyboard listeners (`ArrowUp`, `ArrowDown`).

This allows a specific workflow:
1.  **Look at Map.**
2.  **Press Down:** The selection moves to the next tile.
3.  **Map Scrolls:** The view centers on the new tile.
4.  **Verify:** User checks if the label is correct.

We implemented this by sorting the results spatially (sorting by `Q` then `R`) so that pressing "Down" usually moves you to the visual neighbor, rather than jumping randomly around the map based on confidence score.

---

# Series Conclusion

We started this journey with a screenshot and a problem: **How do we parse game state without a backend?**

By combining **Math** (Axial Coordinates), **Machine Learning** (MobileNet + KNN), and **Frontend Engineering** (React + Canvas), we built a solution that is:

1.  **Private:** No data leaves the device.
2.  **Free:** No API costs.
3.  **Fast:** Zero network latency.
4.  **Flexible:** Works on almost any hex-based 2D game.

**WeeMap Scanner** demonstrates that the browser is becoming a powerful platform for AI application development. You don't always need Python, PyTorch, and a massive cloud bill to solve computer vision problems. Sometimes, all you need is a little bit of Geometry and a few good examples.
