const totalSlides = document.querySelectorAll('.slide').length;
let notesWindow = null;

// Get initial slide from URL hash (e.g., #5 or #slide-5)
function getSlideFromHash() {
    const hash = window.location.hash;
    if (hash) {
        const match = hash.match(/^#(?:slide-)?(\d+)$/);
        if (match) {
            const slideNum = parseInt(match[1], 10);
            if (slideNum >= 1 && slideNum <= totalSlides) {
                return slideNum;
            }
        }
    }
    return 1;
}

let currentSlide = getSlideFromHash();

document.getElementById('totalSlides').textContent = totalSlides;

// Get speaker notes from the DOM
function getNotesForSlide(n) {
    const slide = document.querySelectorAll('.slide')[n - 1];
    const notesEl = slide.querySelector('.speaker-notes');
    return notesEl ? notesEl.innerHTML : '<p>No notes for this slide.</p>';
}

// Get slide title from the DOM
function getSlideTitleForSlide(n) {
    const slide = document.querySelectorAll('.slide')[n - 1];
    const h1 = slide.querySelector('h1');
    return h1 ? h1.textContent : `Slide ${n}`;
}

function showSlide(n) {
    const slides = document.querySelectorAll('.slide');

    if (n > totalSlides) currentSlide = 1;
    if (n < 1) currentSlide = totalSlides;

    slides.forEach(slide => slide.classList.remove('active'));
    slides[currentSlide - 1].classList.add('active');

    document.getElementById('slideNum').textContent = currentSlide;

    // Update URL hash for bookmarking/refresh
    history.replaceState(null, null, `#${currentSlide}`);

    // Update navigation buttons
    document.getElementById('prevBtn').disabled = currentSlide === 1;
    document.getElementById('nextBtn').disabled = currentSlide === totalSlides;

    // Update speaker notes window if it exists
    updateNotesWindow();
}

function changeSlide(n) {
    currentSlide += n;
    showSlide(currentSlide);
}

function openNotesWindow() {
    // Close existing notes window if open
    if (notesWindow && !notesWindow.closed) {
        notesWindow.close();
    }

    // Open new notes window
    notesWindow = window.open('', 'speakerNotes', 'width=900,height=700,scrollbars=yes,resizable=yes');

    if (notesWindow) {
        const title = getSlideTitleForSlide(currentSlide);
        const notes = getNotesForSlide(currentSlide);

        notesWindow.document.write(`
            <html>
            <head>
                <title>Speaker Notes - ${title}</title>
                <style>
                    body {
                        font-family: 'Segoe UI', Arial, sans-serif;
                        padding: 30px;
                        background: #2c3e50;
                        color: white;
                        line-height: 1.6;
                        margin: 0;
                    }
                    .header {
                        background: #34495e;
                        padding: 20px;
                        margin: -30px -30px 30px -30px;
                        border-bottom: 3px solid #3498db;
                    }
                    h1 {
                        color: #3498db;
                        margin: 0;
                        font-size: 1.8em;
                    }
                    .slide-info {
                        color: #bdc3c7;
                        margin-top: 5px;
                        font-size: 1.1em;
                    }
                    h2 {
                        color: #3498db;
                        margin-top: 30px;
                        margin-bottom: 15px;
                        font-size: 1.4em;
                    }
                    h3 {
                        color: #e74c3c;
                        margin-top: 25px;
                        margin-bottom: 10px;
                        font-size: 1.2em;
                    }
                    p, li {
                        font-size: 1em;
                        line-height: 1.6;
                        margin-bottom: 12px;
                    }
                    ul, ol {
                        margin-left: 20px;
                        margin-bottom: 15px;
                    }
                    li {
                        margin-bottom: 8px;
                    }
                    strong {
                        color: #f39c12;
                    }
                    .content {
                        max-width: 800px;
                    }
                    code {
                        background: #1a252f;
                        padding: 2px 6px;
                        border-radius: 3px;
                        font-family: monospace;
                    }
                </style>
            </head>
            <body>
                <div class="header">
                    <h1 id="notesTitle">${title}</h1>
                    <div class="slide-info">Slide <span id="slideNumber">${currentSlide}</span> of ${totalSlides}</div>
                </div>
                <div class="content" id="notesContent">
                    ${notes}
                </div>
            </body>
            </html>
        `);

        notesWindow.document.close();
    }
}

function updateNotesWindow() {
    if (notesWindow && !notesWindow.closed) {
        const title = getSlideTitleForSlide(currentSlide);
        const notes = getNotesForSlide(currentSlide);

        const titleElement = notesWindow.document.getElementById('notesTitle');
        const slideNumberElement = notesWindow.document.getElementById('slideNumber');
        const contentElement = notesWindow.document.getElementById('notesContent');

        if (titleElement && slideNumberElement && contentElement) {
            titleElement.textContent = title;
            slideNumberElement.textContent = currentSlide;
            contentElement.innerHTML = notes;
            notesWindow.document.title = `Speaker Notes - ${title}`;
        }
    }
}

function closeNotesWindow() {
    if (notesWindow && !notesWindow.closed) {
        notesWindow.close();
    }
}

// Keyboard navigation
document.addEventListener('keydown', function(e) {
    if (e.key === 'ArrowLeft') changeSlide(-1);
    if (e.key === 'ArrowRight') changeSlide(1);
    if (e.key === 'Escape') closeNotesWindow();
    if (e.key === 'n' || e.key === 'N') openNotesWindow();
});

// Handle browser back/forward buttons
window.addEventListener('hashchange', function() {
    const slideNum = getSlideFromHash();
    if (slideNum !== currentSlide) {
        currentSlide = slideNum;
        showSlide(currentSlide);
    }
});

// Initialize
showSlide(currentSlide);
